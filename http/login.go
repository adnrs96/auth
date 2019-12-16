package http

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/storyscript/login"
)

//go:generate counterfeiter . TokenProvider

type TokenProvider interface {
	GetConsentURL(state string) string
	GetAccessToken(authCode string) (string, error)
}

type UserInfoFetcher interface {
	GetUser(accessToken string) (login.User, error)
}

type UserRepository interface {
	Save(user login.User) (string, string, error)
}

type LoginHandler struct {
	TokenProvider TokenProvider
}

func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, h.TokenProvider.GetConsentURL("random-state"), http.StatusFound)
}

type StoryscriptClaims struct {
	jwt.StandardClaims
	OwnerUUID string `json:"owner_uuid"`
	TokenUUID string `json:"token_uuid"`
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

	ownerUUID, tokenUUID, _ := h.UserRepository.Save(user)

	claims := StoryscriptClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "storyscript",
			IssuedAt:  time.Now().UTC().Unix(),
			ExpiresAt: time.Now().Add(60 * 60 * 24 * 365).UTC().Unix(),
		},
		OwnerUUID: ownerUUID,
		TokenUUID: tokenUUID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(os.Getenv("SECRET_KEY"))

	http.SetCookie(w, &http.Cookie{
		Name:     "storyscript-jwt",
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 365),
		MaxAge:   60 * 60 * 24 * 365,
		HttpOnly: true,

		Value: tokenString,
	})
}
