package postgres_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var dbConnStr string

var _ = BeforeSuite(func() {
	dbDriver := "postgres"
	dbName := getEnvOrError("TEST_DB_NAME")
	dbUser := getEnvOrError("TEST_DB_USER")
	dbPassword := getEnvOrError("TEST_DB_PASSWORD")
	dbHost := getEnvOrError("TEST_DB_HOST")
	sslMode := "disable"
	searchPath := "public,app_public,app_private,app_secret,app_hidden"
	dbConnStr = fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=%s&search_path=%s", dbDriver, dbUser, dbPassword, dbHost, dbName, sslMode, searchPath)
})

func TestPostgres(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postgres Suite")
}

func getEnvOrError(env string) string {
	value := os.Getenv(env)
	if value == "" {
		Fail(fmt.Sprintf("Environment variable '%s' must be set", env))
	}

	return value
}
