package route

import (
	"github.com/gorilla/mux"
	"net/http"
)

type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	HEAD    HTTPMethod = "HEAD"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	PATCH   HTTPMethod = "PATCH"
	DELETE  HTTPMethod = "DELETE"
	CONNECT HTTPMethod = "CONNECT"
	OPTIONS HTTPMethod = "OPTIONS"
	TRACE   HTTPMethod = "TRACE"
)

type Muxer struct {
	muxer *mux.Router
}

func NewMuxer() *Muxer {
	return &Muxer{
		muxer: mux.NewRouter(),
	}
}

func (m *Muxer) AddRoute(method HTTPMethod, path string, handler func(http.ResponseWriter, *http.Request)) {
	m.muxer.HandleFunc(path, handler).Methods(string(method))
}

func (m *Muxer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.muxer.ServeHTTP(w, r)
}
