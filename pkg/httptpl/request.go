package httptpl

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/netutil/httpctype"
	"github.com/gookit/goutil/netutil/httpreq"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/strutil/textutil"
)

// Template definition for http request
//
// TIP:
//
//	allow use vars in URL, Query, Header, Body, JSON, Form and BodyFile.
type Template struct {
	typ  string // see TypeDefinition
	src  string // the contents
	path string // definition file path
	// Index value
	Index int

	// Version info
	Version string `json:"version"`
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`

	// URL string for request remote, full URL.
	URL    string `json:"url"`
	Method string `json:"method"`
	// Query for request url
	Query map[string]any `json:"query"`
	// Header for request
	Header map[string]string `json:"header"`
	// BasicAuth config
	BasicAuth *httpreq.BasicAuthConf `json:"basic_auth"`

	// Body for request
	Body any `json:"body"`
	// JSON body data, will auto add content type
	JSON any `json:"json"`
	// Form body data, will auto add content type
	Form map[string]any `json:"form"`
	// BodyFile will read file contents as body
	BodyFile string `json:"body_file"`
	// body buffer on build request
	bodyBuf *bytes.Buffer

	// BeforeSend hook
	BeforeSend func(r *http.Request, body *bytes.Buffer)

	// Resp http response data check
	Resp *httpreq.Resp `json:"response"`
}

// NewTemplate instance
func NewTemplate() *Template {
	return &Template{
		typ:     TypeDefinition,
		Version: "1.0-beta",
		Kind:    "kite.http-template",
		Method:  "GET",
	}
}

// FromJSONBytes init definition template
func (t *Template) FromJSONBytes(bs []byte) error {
	return json.Unmarshal(bs, t)
}

// FromJSONString init definition template
func (t *Template) FromJSONString(s string) error {
	return t.FromJSONBytes([]byte(s))
}

// FromJSONFile init definition template
func (t *Template) FromJSONFile(p string) error {
	bs, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	return t.FromJSONBytes(bs)
}

// FromHCString parse from hc-file request part contents
func (t *Template) FromHCString(s string) error {
	t.src = s
	// TODO parse
	return nil
}

var rpl = textutil.NewVarReplacer("{{,}}", func(vp *textutil.VarReplacer) {
	vp.WithParseEnv().WithParseDefault().DisableFlatten().KeepMissingVars()
})

// Send request with variables
func (t *Template) Send(vars maputil.Data, hs map[string]string, opt *httpreq.ReqOption) error {
	req, err := t.BuildRequest(vars, hs)
	if err != nil {
		return err
	}

	if t.BeforeSend != nil {
		t.BeforeSend(req, t.bodyBuf)
	}

	// send request
	resp, err := httpreq.Std().SendRequest(req, opt)
	if err != nil {
		return err
	}

	t.Resp = httpreq.NewResp(resp)
	return nil
}

// BuildRequest instance
func (t *Template) BuildRequest(vars maputil.Data, hs map[string]string) (*http.Request, error) {
	rpl.ResetMissVars()

	// build URL
	apiUrl := rpl.Replace(t.URL, vars)
	if len(t.Query) > 0 {
		q := make(url.Values)
		for k, v := range t.Query {
			strV := strutil.SafeString(v)
			q.Set(k, rpl.Replace(strV, vars))
		}

		if strings.ContainsRune(apiUrl, '?') {
			apiUrl += "&" + q.Encode()
		} else {
			apiUrl += "?" + q.Encode()
		}
	}

	if missVars := rpl.MissVars(); len(missVars) > 0 {
		return nil, errorx.Rawf("input missing variables %v", missVars)
	}

	// build body
	body, err := t.BuildRequestBody(vars)
	if err != nil {
		return nil, err
	}

	// create request
	r, err := http.NewRequest(t.Method, apiUrl, body)
	if err != nil {
		return nil, err
	}

	// set headers for request
	hs = maputil.MergeSMap(t.Header, hs, false)
	for name, val := range hs {
		r.Header.Set(name, rpl.Replace(val, vars))
	}

	// add auth headers
	if t.BasicAuth != nil && t.BasicAuth.IsValid() {
		value := rpl.Replace(t.BasicAuth.String(), vars)
		value = "Basic " + base64.StdEncoding.EncodeToString([]byte(value))
		r.Header.Set("Authorization", value)
	}

	return r, nil
}

// BuildRequestBody for request
func (t *Template) BuildRequestBody(vars maputil.Data) (io.Reader, error) {
	if httpreq.IsNoBodyMethod(t.Method) {
		return nil, nil
	}

	var data string

	if t.JSON != nil {
		bs, err := json.Marshal(t.JSON)
		if err != nil {
			return nil, err
		}

		data = byteutil.String(bs)
		t.Header[httpctype.Key] = httpctype.JSON
	} else if len(t.Form) > 0 {
		data = httpreq.ToQueryValues(t.Form).Encode()
		t.Header[httpctype.Key] = httpctype.Form
	} else if t.BodyFile != "" {
		bs, err := os.ReadFile(t.BodyFile)
		if err != nil {
			return nil, err
		}
		data = byteutil.String(bs)
	} else if t.Body != nil {
		switch typeVal := t.Body.(type) {
		case string:
			data = typeVal
		case []byte:
			data = byteutil.String(typeVal)
		default: // encode by content type
			cType := t.ContentType()
			switch httpctype.ToKind(cType, "") {
			case httpctype.KindJSON:
				bs, err := json.Marshal(t.JSON)
				if err != nil {
					return nil, err
				}

				data = byteutil.String(bs)
			case httpctype.KindForm:
				data = httpreq.ToQueryValues(t.Body).Encode()
			default:
				return nil, errorx.Rawf("invalid body type for request %s", t.URL)
			}
		}
	}

	if len(data) > 0 {
		data = rpl.Replace(data, vars)
		if len(rpl.MissVars()) > 0 {
			return nil, errorx.Rawf("input missing variables %v", rpl.MissVars())
		}

		t.bodyBuf = new(bytes.Buffer)
		t.bodyBuf.WriteString(data)
		return t.bodyBuf, nil
	}
	return nil, nil
}

func (t *Template) RequestString(vars maputil.Data) string {
	req, err := t.BuildRequest(vars, nil)
	if err != nil {
		return ""
	}
	return httpreq.RequestToString(req)
}

// ContentType get
func (t *Template) ContentType() string {
	if ct, ok := t.Header[httpctype.Key]; ok {
		return ct
	}
	return httpctype.Form // default
}

// Type name
func (t *Template) Type() string {
	return t.typ
}

// Path of template file
func (t *Template) Path() string {
	return t.path
}

// String of the request template
func (t *Template) String() string {
	var sb strutil.Builder
	sb.WriteStrings("Name: ", t.Name)
	sb.WriteStrings("Type: ", t.typ)

	if len(t.path) > 0 {
		sb.WriteStrings("Path: ", t.path)
	}

	return sb.String()
}
