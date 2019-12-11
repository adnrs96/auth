package main

import (
	"os"

	"github.com/storyscript/login/gh"
	"github.com/storyscript/login/http"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	ghClient := gh.Client{}
	ghOAuthClient := gh.OAuthClient{
		Config: &oauth2.Config{
			ClientID:     os.Getenv("GH_CLIENT_ID"),
			ClientSecret: os.Getenv("GH_CLIENT_SECRET"),
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
	}
	server := http.Server{
		TokenProvider:   ghOAuthClient,
		UserInfoFetcher: ghClient,
	}

	if err := server.Start(); err != nil {
		panic(err)
	}
}
