package httptpl

import (
	"io/fs"
	"os"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/pkg/pkgutil"
)

// Templates type constants
// const TypeTemplate   = "template"    // request definition template files
const (
	TypeDefinition = "definition"  // request definition spec/template files
	TypeHttpClient = "http-client" // a IDE http-client file
)

// Templates request set definitions for a group
type Templates struct {
	init bool
	name string
	// types: definition, http-client. see TypeDefinition
	typ string
	// src the hc-file contents
	src string
	// typ=TypeHttpClient: path is the hc-file path.
	//
	// typ=TypeTemplate: path is the template file dir path
	path string
	// ext names for auto load template. eg: [json, yaml]
	exts []string
	// loaded template set
	set map[string]*Template
	// name to index
	// err error
}

// NewTemplates instance
func NewTemplates(name string) *Templates {
	return &Templates{
		name: name,
		typ:  TypeDefinition,
		set:  make(map[string]*Template),
	}
}

func (ts *Templates) FromHCFile(hcFile string) error {
	if ts.init {
		return nil
	}

	ts.init = true
	ts.path = hcFile
	bs, err := os.ReadFile(hcFile)
	if err != nil {
		return err
	}

	return ts.FromHCString(string(bs))
}

func (ts *Templates) FromHCString(s string) error {
	ts.src = s
	ss := strutil.Split("\n"+s, RequestSplit)

	for i, part := range ss {
		t := NewTemplate()
		t.typ = ts.typ
		t.Index = i

		if err := t.FromHCString(part); err != nil {
			return err
		}

		ts.set[t.Name] = t
	}
	return nil
}

func (ts *Templates) Send(keywords string) error {
	return nil
}

func (ts *Templates) Get(name string) *Template {
	t, _ := ts.Lookup(name)
	return t
}

func (ts *Templates) Lookup(name string) (*Template, error) {
	// find and load template file
	if ts.IsDefType() {
		t, err := ts.LoadTemplate(name)
		if err != nil {
			return nil, err
		}
		return t, nil
	}

	// is from hc-file
	t, ok := ts.set[name]
	if ok {
		return t, nil
	}
	return nil, errorx.Rawf("http-client template %q not found", name)
}

// LoadTemplate file and register
func (ts *Templates) LoadTemplate(name string) (*Template, error) {
	c := pkgutil.NewConfig()
	t := NewTemplate()

	tplFile := ts.path + "/" + name
	if fsutil.IsFile(tplFile) {
		t.Name = name
		t.path = tplFile
		t.Index = len(ts.set)
		if err := c.LoadFiles(tplFile); err != nil {
			return nil, err
		}

		if err := c.Decode(t); err != nil {
			return nil, err
		}
		return t, nil
	}

	// try with ext for load
	for _, ext := range ts.exts {
		tplFile := ts.path + "/" + name + "." + ext
		if fsutil.IsFile(tplFile) {
			t.Name = name
			t.path = tplFile
			t.Index = len(ts.set)
			if err := c.LoadFiles(tplFile); err != nil {
				return nil, err
			}

			if err := c.Decode(t); err != nil {
				return nil, err
			}
			return t, nil
		}
	}
	return nil, errorx.Rawf("not found template %q on %q", name, ts.name)
}

// GetByPath get Template by request uri path.
func (ts *Templates) GetByPath(path string) *Template {
	return nil
}

func (ts *Templates) FindOne(keywords string) *Template {
	return nil
}

func (ts *Templates) Search(keywords string, limit int) []*Template {
	return nil
}

// Each template
func (ts *Templates) Each(fn func(t *Template)) {
	if ts.IsHcType() {
		for _, t := range ts.set {
			fn(t)
		}
		return
	}

	err := ts.LoadAll()
	if err != nil {
		return
	}
}

// LoadAll request definition template files
func (ts *Templates) LoadAll() error {
	if ts.init || ts.IsHcType() {
		return nil
	}

	ts.init = true
	return fsutil.FindInDir(ts.path, func(fPath string, ent fs.DirEntry) error {
		_, err := ts.LoadTemplate(ent.Name())
		return err
	}, fsutil.OnlyFindFile, fsutil.IncludeSuffix(ts.exts...))
}

// All templates
func (ts *Templates) All() map[string]*Template {
	return ts.set
}

// IsDefType templates
func (ts *Templates) IsDefType() bool {
	return ts.typ == TypeDefinition
}

// IsTplType templates
func (ts *Templates) IsTplType() bool {
	return ts.typ == TypeDefinition
}

// IsHcType templates
func (ts *Templates) IsHcType() bool {
	return ts.typ == TypeHttpClient
}

func (ts *Templates) String() string {
	var sb strutil.Builder
	sb.WriteStrings("Name: ", ts.name, "\n")
	sb.WriteStrings("Type: ", ts.typ, "\n")
	sb.WriteStrings("Path: ", ts.path, "\n")

	if ts.IsTplType() {
		sb.WriteString("Exts: ")
		sb.WriteAnys(ts.exts)
	} else {
		sb.WriteAnys("Parts number: ", len(ts.set))
	}

	sb.WriteByteNE('\n')
	return sb.String()
}
