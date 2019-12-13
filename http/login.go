package http

import (
	"fmt"
	"net/http"

	"github.com/storyscript/login"
)

type TokenProvider interface {
	GetConsentURL(state string) string
	GetAccessToken(authCode string) (string, error)
}

type UserInfoFetcher interface {
	GetUser(accessToken string) (login.User, error)
}

type UserRepository interface {
	Save(user login.User) error
}

type LoginHandler struct {
	TokenProvider TokenProvider
}

func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, h.TokenProvider.GetConsentURL("random-state"), http.StatusFound)
}

type CallbackHandler struct {
	TokenProvider   TokenProvider
	UserInfoFetcher UserInfoFetcher

	UserRepository UserRepository
}

func (h CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// state := r.FormValue("state")
	authCode := r.FormValue("code")

	// check state is valid

	accessToken, _ := h.TokenProvider.GetAccessToken(authCode)
	user, _ := h.UserInfoFetcher.GetUser(accessToken)

	user.Name = "will"
	user.OAuthToken = accessToken

	_ = h.UserRepository.Save(user)

	w.Write([]byte(fmt.Sprintf("%v", user)))
}