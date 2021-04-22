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

func (g *HandlerGroup) viewHandle(w http.ResponseWriter, r *http.Request) {
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

func startServer(addr, pacFile, maxAge string) error {
	hg := &HandlerGroup{}
	err := hg.loadPData(pacFile, maxAge)
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	mux := http.NewServeMux()
	mux.HandleFunc("/",hg.missingHandle)
	mux.HandleFunc("/pac",hg.viewHandle)
	mux.HandleFunc("/wpad.dat", hg.wpadHandle)
	mux.HandleFunc("/edit/", hg.editHandle)
	mux.HandleFunc("/save/", hg.saveHandle)
	mux.HandleFunc("/favicon.ico", hg.missingHandle)

	slog.Info("server start on", addr)

	// TODO loggingHandler := logging.NewApacheLoggingHandler(mux, os.Stdout)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
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
