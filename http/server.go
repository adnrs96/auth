package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	TokenProvider   TokenProvider
	UserInfoFetcher UserInfoFetcher

	UserRepository UserRepository
	TokenGenerator TokenGenerator
}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func (s Server) Start() error {
	var routes = []route{
		{
			"Login",
			"GET",
			"/login",
			s.HandleLogin,
		},
		{
			"Callback",
			"GET",
			"/callback",
			s.HandleCallback,
		},
		{
			"Healthcheck",
			"Get",
			"/healthcheck",
			s.HandleHealthcheck,
		},
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(route.HandlerFunc)
	}

	return http.ListenAndServe(":3000", router)
}
