package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/storyscript/auth/gh"
	"github.com/storyscript/auth/http"
	"github.com/storyscript/auth/jwt"
	"github.com/storyscript/auth/postgres"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	ghClient := gh.UserClient{}
	ghOAuthClient := gh.OAuthClient{
		Config: &oauth2.Config{
			ClientID:     getEnvOrPanic("GH_CLIENT_ID"),
			ClientSecret: getEnvOrPanic("GH_CLIENT_SECRET"),
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
	}

	postgresClient := postgres.Client{
		DB: openDB(),
	}

	jwtGenerator := jwt.Generator{
		SigningKey: getEnvOrPanic("JWT_SIGNING_KEY"),
	}

	server := http.Server{
		TokenProvider:   ghOAuthClient,
		UserInfoFetcher: ghClient,

		UserRepository: postgresClient,
		TokenGenerator: jwtGenerator,

		Domain:      getEnvOrPanic("DOMAIN"),
		RedirectURI: getEnvOrPanic("POST_LOGIN_REDIRECT_URI"),
	}

	if err := server.Start(); err != nil {
		panic(err)
	}
}

func openDB() *sql.DB {
	dbConnectionString := dbConnectionString()
	dbDriver := "postgres"
	db, err := sql.Open(dbDriver, dbConnectionString)
	if err != nil {
		panic(err)
	}

	return db
}

func dbConnectionString() string {
	dbConnectionString := getEnvOrPanic("DB_CONNECTION_STRING")
	sslMode := os.Getenv("SSL_MODE")

	if sslMode == "" {
		return dbConnectionString
	}

	return fmt.Sprintf("%s&sslmode=%s", dbConnectionString, sslMode)
}

func getEnvOrPanic(env string) string {
	value := os.Getenv(env)
	if value == "" {
		panic(fmt.Sprintf("Environment variable '%s' must be set", env))
	}

	return value
}
