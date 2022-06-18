package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/honyshyota/constanta-rest-api/internal/app/model"
	"github.com/honyshyota/constanta-rest-api/internal/app/store"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "usersession"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
	errNotFound                 = errors.New("page not found")
	errNotModified              = errors.New("unable to modify record")
	errInvalidStatus            = errors.New("invalid transaction status")
	errTransNotExists           = errors.New("transaction does not exist")
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
	// Logger RequestID and CORS
	srv.router.Use(srv.setRequestID)
	srv.router.Use(srv.logRequest)
	srv.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	// external funcs
	srv.router.HandleFunc("/users", srv.handleUsersCreate()).Methods("POST")
	srv.router.HandleFunc("/sessions", srv.handleSessionsCreate()).Methods("POST")

	// private funcs
	private := srv.router.PathPrefix("/private").Subrouter()
	private.Use(srv.authenticateUser)
	private.HandleFunc("/whoami", srv.handleWhoami()).Methods("GET")
	private.HandleFunc("/pay", srv.handleTransactionCreate()).Methods("POST")
	private.HandleFunc("/update", srv.handleTransactionAdmin()).Methods("POST")
	private.HandleFunc("/checkstatus", srv.handleCheckTransStatus()).Methods("GET")
	private.HandleFunc("/findtrans", srv.handleFindTransactions()).Methods("GET")
	private.HandleFunc("/delete", srv.handleDeleteTransaction()).Methods("POST")
}

func (srv *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (srv *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := srv.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Since(start),
		)
	})
}

// Authentification users
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

// Yourself monitor)))
func (srv *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

// Transactions create
func (srv *server) handleTransactionCreate() http.HandlerFunc {
	type request struct {
		Pay      string `json:"pay"`
		Currency string `json:"currency"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*model.User)

		pay, err := strconv.ParseFloat(req.Pay, 32)
		if err != nil {
			srv.error(w, r, http.StatusInternalServerError, err)
			return
		}

		transaction := &model.Transaction{
			UserID:     u.ID,
			Email:      u.Email,
			Pay:        float32(pay),
			Currency:   req.Currency,
			TimeCreate: time.Now(),
			TimeUpdate: time.Now(),
			Status:     randomizer(), // Simple randomizer
		}

		if err := srv.store.Transaction().Create(transaction); err != nil {
			srv.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		srv.respond(w, r, http.StatusCreated, transaction)
	}
}

// Transaction status update (Admin use only)
func (srv *server) handleTransactionAdmin() http.HandlerFunc {
	type request struct {
		TransID int    `json:"trans_id,string"`
		Status  string `json:"trans_status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)
		if u.Email == "admin@example.com" {
			req := &request{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				srv.error(w, r, http.StatusBadRequest, err)
				return
			}

			trans, err := srv.store.Transaction().FindTrans(req.TransID)
			if err != nil {
				srv.error(w, r, http.StatusNotFound, err)
				return
			}

			if trans.Status != "error" && trans.Status != "success" && trans.Status != "failure" {
				if err := srv.store.Transaction().StatusUpdate(req.Status, req.TransID); err != nil {
					srv.error(w, r, http.StatusNotFound, err)
					return
				}
			} else {
				srv.error(w, r, http.StatusNotModified, errNotModified)
			}

			srv.respond(w, r, http.StatusCreated, req)
		} else {
			srv.error(w, r, http.StatusBadRequest, errNotFound)
			return
		}
	}
}

// Transaction check status (admin use only)
func (srv *server) handleCheckTransStatus() http.HandlerFunc {
	type request struct {
		TransID int `json:"trans_id,string"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)
		if u.Email == "admin@example.com" {
			req := &request{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				srv.error(w, r, http.StatusBadRequest, err)
				return
			}

			trans, err := srv.store.Transaction().FindTrans(req.TransID)
			if err != nil {
				srv.error(w, r, http.StatusNotFound, err)
				return
			}

			srv.respond(w, r, http.StatusOK, trans.Status)

		} else {
			srv.error(w, r, http.StatusBadRequest, errNotFound)
			return
		}
	}
}

// Find all user transactions
func (srv *server) handleFindTransactions() http.HandlerFunc {
	type request struct {
		Data string `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, r, http.StatusBadRequest, err)
			return
		}

		trans, err := srv.store.Transaction().Find(req.Data)
		if err != nil {
			srv.error(w, r, http.StatusBadRequest, err)
		}

		srv.respond(w, r, http.StatusOK, trans)
	}
}

// Delete transaction
func (srv *server) handleDeleteTransaction() http.HandlerFunc {
	type request struct {
		TransID int `json:"trans_id,string"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, r, http.StatusBadRequest, err)
			return
		}

		trans, err := srv.store.Transaction().FindTrans(req.TransID)
		if err != nil {
			srv.error(w, r, http.StatusNotFound, err)
			return
		}

		if trans.UserID != u.ID {
			srv.error(w, r, http.StatusBadRequest, errTransNotExists)
		}

		if trans.Status == "success" || trans.Status == "error" {
			if err := srv.store.Transaction().Delete(req.TransID); err != nil {
				fmt.Println(trans.UserID, u.ID)
				srv.error(w, r, http.StatusInternalServerError, err)
			}
			srv.respond(w, r, http.StatusCreated, "transaction deleted")
		} else {
			srv.error(w, r, http.StatusBadRequest, errInvalidStatus)
		}
	}
}

// Users create
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

// Set sessions
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
