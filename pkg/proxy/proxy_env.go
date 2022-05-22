package proxy

const (
	HttpKey  = "HTTP_PROXY"
	HttpsKey = "HTTPS_PROXY"
)

// Config struct
type Config struct {
	// ProxyHost url. eg:  http://127.0.0.1:1080
	ProxyHost string
}

func New(conf *Config) {

}
