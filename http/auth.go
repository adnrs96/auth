package http

import (
	"net/http"
	"time"

	"github.com/storyscript/login"
)

//go:generate counterfeiter . TokenProvider

type TokenProvider interface {
	GetConsentURL(state string) string
	GetAccessToken(authCode string) (string, error)
}

//go:generate counterfeiter . UserInfoFetcher

type UserInfoFetcher interface {
	GetUser(accessToken string) (login.User, error)
}

//go:generate counterfeiter . UserRepository

type UserRepository interface {
	Save(user login.User) (string, error)
}

//go:generate counterfeiter . TokenGenerator

type TokenGenerator interface {
	Generate(ownerUUID string) (string, error)
}

func (s Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.TokenProvider.GetConsentURL("random-state"), http.StatusFound)
}

func (s Server) HandleCallback(w http.ResponseWriter, r *http.Request) {
	authCode := r.FormValue("code")

	accessToken, err := s.TokenProvider.GetAccessToken(authCode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := s.UserInfoFetcher.GetUser(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.OAuthToken = accessToken

	ownerUUID, err := s.UserRepository.Save(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := s.TokenGenerator.Generate(ownerUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "storyscript-access-token",
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 365),
		MaxAge:   60 * 60 * 24 * 365,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,

		Value: token,
	})
}
