package http

import "net/http"

type Server struct {
	TokenProvider   TokenProvider
	UserInfoFetcher UserInfoFetcher

	UserRepository UserRepository
	TokenGenerator TokenGenerator
}

func (s Server) Start() error {
	loginHandler := LoginHandler{
		TokenProvider: s.TokenProvider,
	}

	callbackHandler := CallbackHandler{
		TokenProvider:   s.TokenProvider,
		UserInfoFetcher: s.UserInfoFetcher,

		UserRepository: s.UserRepository,
		TokenGenerator: s.TokenGenerator,
	}

	var routes = []Route{
		Route{
			"Login",
			"GET",
			"/login",
			loginHandler,
		},
		Route{
			"Callback",
			"GET",
			"/callback",
			callbackHandler,
		},
		Route{
			"Healthcheck",
			"Get",
			"/healthcheck",
			HealthcheckHandler{},
		},
	}

	return http.ListenAndServe(":3000", NewRouter(routes))
}

type OAuthHandler struct{}
