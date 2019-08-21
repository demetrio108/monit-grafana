package router

import (
	"github.com/demetrio108/monit-grafana/exporter"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Router struct {
	addr   string
	router *mux.Router
}

func New(addr string, exporters []*exporter.Exporter) *Router {
	r := Router{}

	r.addr = addr
	r.router = mux.NewRouter()

	for _, exp := range exporters {
		log.Print("Registering exporter endpoint ", addr, "/"+exp.Name()+"/")
		r.router.HandleFunc("/"+exp.Name()+"/", exp.RootHandler)
		log.Print("Registering exporter endpoint ", addr, "/"+exp.Name()+"/search")
		r.router.HandleFunc("/"+exp.Name()+"/search", exp.SearchHandler)
		log.Print("Registering exporter endpoint ", addr, "/"+exp.Name()+"/query")
		r.router.HandleFunc("/"+exp.Name()+"/query", exp.QueryHandler)
	}

	return &r
}

func (r *Router) ListenAndServe() error {
	log.Printf("Ready to serve at %s", r.addr)
	return http.ListenAndServe(r.addr, r.router)
}
