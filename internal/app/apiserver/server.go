package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/honyshyota/constanta-rest-api/internal/app/model"
	"github.com/honyshyota/constanta-rest-api/internal/app/store"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "usersession"
	ctxKeyUser  ctxKey = iota
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	srv := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	srv.configureRouter()

	return srv
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func (srv *server) configureRouter() {
	srv.router.HandleFunc("/users", srv.handleUsersCreate()).Methods("POST")
	srv.router.HandleFunc("/sessions", srv.handleSessionsCreate()).Methods("POST")
}

func (srv *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := srv.sessionStore.Get(r, sessionName)
		if err != nil {
			srv.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["id"]
		if !ok {
			srv.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := srv.store.User().Find(id.(int))
		if err != nil {
			srv.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (srv *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
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
			Email:    req.Email,
			Password: req.Password,
		}
		if err := srv.store.User().Create(u); err != nil {
			srv.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()

		srv.respond(w, r, http.StatusCreated, u)
	}
}

func (srv *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := srv.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			srv.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := srv.sessionStore.Get(r, sessionName)
		if err != nil {
			srv.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["id"] = u.ID
		if err := srv.sessionStore.Save(r, w, session); err != nil {
			srv.error(w, r, http.StatusInternalServerError, err)
			return
		}

		srv.respond(w, r, http.StatusOK, nil)
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
