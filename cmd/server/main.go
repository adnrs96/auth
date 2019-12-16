package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/storyscript/login/gh"
	"github.com/storyscript/login/http"
	"github.com/storyscript/login/jwt"
	"github.com/storyscript/login/postgres"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	ghClient := gh.UserClient{}
	ghOAuthClient := gh.OAuthClient{
		Config: &oauth2.Config{
			ClientID:     os.Getenv("GH_CLIENT_ID"),
			ClientSecret: os.Getenv("GH_CLIENT_SECRET"),
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
	}

	postgresClient := postgres.Client{
		DB: openDB(),
	}

	jwtGenerator := jwt.Generator{
		SigningKey: os.Getenv("SECRET_KEY"),
	}

	server := http.Server{
		TokenProvider:   ghOAuthClient,
		UserInfoFetcher: ghClient,

		UserRepository: postgresClient,
		TokenGenerator: jwtGenerator,
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

	return fmt.Sprintf("%s&sslmode=%s", getEnvOrPanic("DB_CONNECTION_STRING"), sslMode)
}

func getEnvOrPanic(env string) string {
	value := os.Getenv(env)
	if value == "" {
		panic(fmt.Sprintf("Environment variable '%s' must be set", env))
	}

	return value
}
