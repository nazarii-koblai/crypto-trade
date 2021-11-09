package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/crypto-trade/config"
	"github.com/crypto-trade/handler"
	"github.com/crypto-trade/storage"
	"github.com/crypto-trade/token"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	addr          = ":8080"
	versionPrefix = "/v1"

	signupSubroute = "/signup"
	signinSubroute = "/signin"

	transferSubroute     = "/transfer"
	transactionsSubroute = "/transactions"
)

// Server represents server.
type Server struct {
	logger *logrus.Logger
	srv    *http.Server
	config config.Config
	db     *storage.Storage
	token  token.Token
}

// New creates new Server.
func New(
	logger *logrus.Logger,
	config config.Config,
	db *storage.Storage,
	token token.Token,
) *Server {
	s := new(Server)
	s.logger = logger
	s.config = config
	s.db = db
	s.token = token
	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.setupMux(),
	}
	return s
}

// Run starts server.
func (s *Server) Run() {
	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return
	}
	if err != nil {
		s.logger.Fatal(err)
	}
}

func (s *Server) setupMux() *mux.Router {
	r := mux.NewRouter()

	v1 := r.PathPrefix(versionPrefix).Subrouter()
	v1.Use(csrf.Protect(s.config.CSRF.Key, csrf.Secure(false)))

	auth := handler.NewAuth(s.logger, s.db, s.token)
	v1.Path(signupSubroute).Handler(auth.Signup()).Methods(http.MethodPost)
	v1.Path(signinSubroute).Handler(auth.Signin()).Methods(http.MethodPost, http.MethodGet)

	transferSubroute := v1.PathPrefix(transferSubroute).Subrouter()
	transactionsSubroute := v1.PathPrefix(transactionsSubroute).Subrouter()

	_, _ = transactionsSubroute, transferSubroute

	return r
}

// Shutdown shutdowns server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
