package apiserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/honyshyota/constanta-rest-api/internal/app/model"
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
	type request struct {
		ID       int    `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			ID:         req.ID,
			Email:      req.Email,
			Password:   req.Password,
			TimeCreate: time.Now(),
			TimeUpdate: time.Now(),
		}
		if err := srv.store.User().Create(u); err != nil {
			srv.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()

		srv.respond(w, r, http.StatusCreated, u)
	}
}

func (srv *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	srv.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (srv *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
