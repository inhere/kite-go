package lcproxy

import (
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/strutil"
)

const (
	HttpKey  = "HTTP_PROXY"
	HttpsKey = "HTTPS_PROXY"
)

// LocalProxy local proxy ENV setting.
type LocalProxy struct {
	// HttpProxy host url. eg: http://127.0.0.1:1080
	HttpProxy string `json:"http_proxy"`
	// HttpsProxy host url. eg: http://127.0.0.1:1080
	HttpsProxy string `json:"https_proxy"`
}

// NewLocalProxy instance
func NewLocalProxy() *LocalProxy {
	return &LocalProxy{}
}

// Apply proxy ENV setting.
func (lp *LocalProxy) Apply(beforeFn func(lp *LocalProxy)) {
	if !lp.IsEmpty() {
		beforeFn(lp)
		envutil.SetEnvs(HttpKey, lp.HttpProxy, HttpsKey, lp.HttpsProxy)
	}
}

// RunFunc on open proxy ENV, unset ENV after run.
func (lp *LocalProxy) RunFunc(fn func()) {
	if lp.IsEmpty() {
		fn()
		return
	}

	lp.Apply(nil)
	defer lp.Unset()
	fn()
}

// EnvKeys for proxy ENV
func (lp *LocalProxy) EnvKeys() []string {
	return []string{HttpKey, HttpsKey}
}

// Unset proxy ENV
func (lp *LocalProxy) Unset() {
	if !lp.IsEmpty() {
		envutil.UnsetEnvs(HttpKey, HttpsKey)
	}
}

// IsEmpty proxy ENV setting.
func (lp *LocalProxy) IsEmpty() bool {
	return strutil.Valid(lp.HttpProxy, lp.HttpsProxy) == ""
}
