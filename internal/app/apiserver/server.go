package apiserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/honyshyota/constanta-rest-api/internal/app/store"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	srv := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	srv.configureRouter()

	return srv
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func (srv *server) configureRouter() {
	srv.router.HandleFunc("/users", srv.handleUsersCreate()).Methods("POST")
}

func (srv *server) handleUsersCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}
