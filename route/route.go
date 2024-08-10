package route

import (
	"context"
	"crud/controller"
	"crud/repository"
	"crud/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
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

func (m *Muxer) AddMiddleware(middleware mux.MiddlewareFunc) {
	m.muxer.Use(middleware)
}

var Router *Muxer

func init() {
	Router = NewMuxer()

	conn, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		panic(err)
	}

	db := conn.Database(os.Getenv("MONGODB_DB"))

	userColl := repository.NewUserRepository(db, "users")

	Router.AddMiddleware(util.InsertToRequest("userDb", userColl))

	Router.AddRoute(POST, "/register", controller.Register)
	Router.AddRoute(POST, "/login", controller.Login)
}
