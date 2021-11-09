package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/crypto-trade/model"
	"github.com/crypto-trade/repository"
	"github.com/crypto-trade/storage"
	"github.com/crypto-trade/token"
	"github.com/crypto-trade/web"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
)

var (
	errInvalidCredentials = errors.New("invalid credentials")
)

// Auth represents auth handler.
type Auth struct {
	logger *logrus.Logger
	repo   repository.User
	token  token.Token
}

// NewAuth returns new auth hanlder.
func NewAuth(logger *logrus.Logger, repo repository.User, token token.Token) *Auth {
	return &Auth{
		logger: logger,
		repo:   repo,
		token:  token,
	}
}

// Signup returns signup handler func.
func (a *Auth) Signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err)
			return
		}

		defer r.Body.Close()

		var user model.User
		if err = json.Unmarshal(body, &user); err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err)
			return
		}

		if err = user.Validate(); err != nil {
			web.RespondWithError(w, http.StatusBadRequest, err)
			return
		}

		if err := user.Password.Hash(); err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err)
			return
		}

		user, err = a.repo.AddNewUser(r.Context(), user)
		if errors.Is(err, storage.ErrDuplicateKeyValueViolatesUniqueConstraint) {
			web.RespondWithError(w, http.StatusInternalServerError, errors.New("user with such email already registred"))
			return
		}

		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			Username string         `json:"username"`
			Surname  string         `json:"surname"`
			Wallets  []model.Wallet `json:"wallets"`
		}

		b, err := json.Marshal(response{
			Username: string(user.Name),
			Surname:  string(user.Surname),
			Wallets:  user.Wallets,
		})
		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err)
			return
		}

		web.RespondWithSuccess(w, http.StatusOK, b)
	}
}

// Signin returns signin handler func.
func (a *Auth) Signin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err)
			return
		}

		defer r.Body.Close()

		type request struct {
			Email    model.Email    `json:"email"`
			Password model.Password `json:"password"`
		}

		var req request
		if err = json.Unmarshal(body, &req); err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err)
			return
		}

		hasedPassword, err := a.repo.HashedPassword(r.Context(), req.Email)
		if errors.Is(err, storage.ErrNotFound) {
			web.RespondWithError(w, http.StatusBadRequest, errInvalidCredentials)
			return
		}

		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err)
			return
		}

		if err = req.Password.CompareWithHashed([]byte(hasedPassword)); err != nil {
			web.RespondWithError(w, http.StatusBadRequest, errInvalidCredentials)
			return
		}

		token, err := a.token.GenerateWithClaims(jwt.MapClaims{
			"email": req.Email,
		})
		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			Token string `json:"token"`
		}

		b, err := json.Marshal(response{
			Token: token,
		})
		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("X-CSRF-Token", csrf.Token(r))

		web.RespondWithSuccess(w, http.StatusOK, b)
	}
}
