package pacutil

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/gookit/rux"
	"github.com/gookit/slog"
)

const editTpl = `
<!DOCTYPE html>
<html>
  <head>
    <title>Editing PAC</title>
    <meta charset="utf-8">
  </head>
  <body>
    <h1>Editing PAC</h1>
    <div>{{printf "%s" .Message}}</div>

    <form action="/save/" method="POST" accept-charset="utf-8">
      <div><textarea name="body" rows="20" cols="80">{{printf "%s" .Body}}</textarea></div>
      <div><input type="submit" value="Save"></div>
    </form>
  </body>
</html>
`

type PData struct {
	file    string
	Body    []byte
	Message string
	Etag    string
	MaxAge  string
}

func (p *PData) save() error {
	return ioutil.WriteFile(p.file, p.Body, 0600)
}

func (p *PData) edit(w http.ResponseWriter) {
	// t, _ := template.ParseFiles(*templatePath + "/edit.html")
	t := template.New("edit")
	template.Must(t.Parse(editTpl))
	// t.Execute(w, p)
	slog.ErrorT(t.Execute(w, p))
}

// refer https://github.com/ceaser/pac-server/blob/master/main.go
type HandlerGroup struct {
	pd *PData
}

func (g *HandlerGroup) loadPData(pacFile, maxAge string) error {
	body, err := ioutil.ReadFile(pacFile)
	if err != nil {
		return err
	}

	g.pd = &PData{Body: body, MaxAge: maxAge}

	// Generate etag
	hasher := md5.New()
	hasher.Write(g.pd.Body)
	g.pd.Etag = hex.EncodeToString(hasher.Sum(nil))

	return err
}

func (g *HandlerGroup) gfwHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	dst, err := DecodeGfwList(pacOpts.gwfile)
	if err != nil {
		_, err = w.Write([]byte(`error`))
		slog.ErrorT(err)
		return
	}

	_, err = w.Write(dst)
	slog.ErrorT(err)
}

func (g *HandlerGroup) viewHandle(w http.ResponseWriter, r *http.Request) {
	e := g.pd.Etag

	// w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig; charset=utf-8")
	// Caching
	w.Header().Set("Etag", e)
	w.Header().Set("Cache-Control", "max-age="+g.pd.MaxAge)
	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, e) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	// fmt.Fprintf(w, "%s", g.pd.Body)
	_, err := w.Write(g.pd.Body)
	slog.ErrorT(err)
}

func (g *HandlerGroup) editHandle(w http.ResponseWriter, r *http.Request) {
	g.pd.edit(w)
}

func (g *HandlerGroup)  saveHandle(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")

	g.pd.Body = []byte(body)
	// p := &PData{Body: []byte(body)}
	err := g.pd.save()
	if err != nil {
		g.pd.Message = fmt.Sprintf("Error: %s", err.Error())
		g.pd.edit(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (g *HandlerGroup) wpadHandle(w http.ResponseWriter, r *http.Request) {
	// TODO: Refactor duplication between wpadHandler and viewHandler
	// p, _ := loadPage()
	e := g.pd.Etag

	// Caching
	w.Header().Set("Etag", e)
	w.Header().Set("Cache-Control", "max-age="+g.pd.MaxAge)
	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, e) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig; charset=utf-8")
	// fmt.Fprintf(w, "%s", g.pd.Body)
	_, err := w.Write(g.pd.Body)
	slog.ErrorT(err)
}

func (g *HandlerGroup)  missingHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	// fmt.Fprint(w, "404")
	_, err := w.Write([]byte("404"))
	slog.ErrorT(err)
}

func startServer(opts PacOpts) error {
	hg := &HandlerGroup{}
	err := hg.loadPData(opts.file, opts.maxAge)
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	h := createHttpRux(hg)
	// h := createHttpMux(hg)

	slog.Info("server start on", opts.addr)

	// TODO loggingHandler := logging.NewApacheLoggingHandler(mux, os.Stdout)
	server := &http.Server{
		Addr:    opts.addr,
		Handler: h,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
				slog.Error(err)
			} else {
				color.Success.Println("Server closed")
			}
		}
	}()
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	return server.Shutdown(ctx)
}

func createHttpRux(hg *HandlerGroup) *rux.Router {
	r := rux.New()
	rux.Debug(true)

	r.GET("/", rux.WrapHTTPHandlerFunc(hg.missingHandle))
	r.GET("/pac",rux.WrapHTTPHandlerFunc(hg.viewHandle))
	r.GET("/gfw", rux.WrapHTTPHandlerFunc(hg.gfwHandle))
	r.GET("/wpad.dat", rux.WrapHTTPHandlerFunc(hg.wpadHandle))
	r.GET("/edit/", rux.WrapHTTPHandlerFunc(hg.editHandle))
	r.GET("/save/", rux.WrapHTTPHandlerFunc(hg.saveHandle))
	r.GET("/favicon.ico", rux.WrapHTTPHandlerFunc(hg.missingHandle))

	return r
}

func createHttpMux(hg *HandlerGroup) *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("/",hg.missingHandle)
	r.HandleFunc("/pac",hg.viewHandle)
	r.HandleFunc("/gfw", hg.gfwHandle)
	r.HandleFunc("/wpad.dat", hg.wpadHandle)
	r.HandleFunc("/edit/", hg.editHandle)
	r.HandleFunc("/save/", hg.saveHandle)
	r.HandleFunc("/favicon.ico", hg.missingHandle)

	return r
}
