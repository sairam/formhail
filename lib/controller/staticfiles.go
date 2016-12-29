package controller

import (
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"../common"

	"github.com/gorilla/mux"
)

func initStatic(r *mux.Router) {
	// TODO add 404 handler

	// Handle Static Files
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, strings.Join([]string{"public", "favicon.ico"}, string(os.PathSeparator)))
	})

	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "robots.txt", time.Now(), strings.NewReader("User-agent: *\nSitemap: /sitemap.xml"))
	})

	r.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		t, _ := template.New("foo").Parse(sitemapTemplate)
		t.Execute(w, &struct {
			Host    string
			Pages   []string
			Changed string
		}{common.Config.WebsiteURL, common.Config.StaticFilesList, time.Now().Format("2006-01-02T15:00:00-07:00")})
	})

}

var sitemapTemplate = `<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  {{$host:=.Host}}
  {{$changed:=.Changed}}
  {{range $page:=.Pages}}
  <url>
    <loc>{{$host}}{{$page}}</loc>
    <lastmod>{{$changed}}</lastmod>
  </url>
  {{end}}
</urlset>`
