package http

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	Save(user login.User) error
}

//go:generate counterfeiter . TokenGenerator

type TokenGenerator interface {
	Generate(user login.User) (string, error)
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
	TokenGenerator TokenGenerator
}

func (h CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authCode := r.FormValue("code")

	accessToken, err := h.TokenProvider.GetAccessToken(authCode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.UserInfoFetcher.GetUser(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.OAuthToken = accessToken

	if err := h.UserRepository.Save(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := h.TokenGenerator.Generate(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	// claims := StoryscriptClaims{
	// 	StandardClaims: jwt.StandardClaims{
	// 		Issuer:    "storyscript",
	// 		IssuedAt:  time.Now().UTC().Unix(),
	// 		ExpiresAt: time.Now().Add(60 * 60 * 24 * 365).UTC().Unix(),
	// 	},
	// 	OwnerUUID: ownerUUID,
	// 	TokenUUID: tokenUUID,
	// }
	//
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//
	// fmt.Println(os.Getenv("SECRET_KEY"))
	//
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
