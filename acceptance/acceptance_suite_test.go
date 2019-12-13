package acceptance_test

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/sclevine/agouti"
)

var (
	dbConnStr    string
	serverPath   string
	agoutiDriver *agouti.WebDriver
)

var _ = BeforeSuite(func() {
	dbDriver := "postgres"
	dbName := getEnvOrError("TEST_DB_NAME")
	dbUser := getEnvOrError("TEST_DB_USER")
	dbPassword := getEnvOrError("TEST_DB_PASSWORD")
	dbHost := getEnvOrError("TEST_DB_HOST")
	sslMode := "disable"
	searchPath := "public,app_public,app_private,app_secret,app_hidden"
	dbConnStr = fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=%s&search_path=%s", dbDriver, dbUser, dbPassword, dbHost, dbName, sslMode, searchPath)

	var err error
	serverPath, err = gexec.Build("../cmd/server/main.go")
	Expect(err).NotTo(HaveOccurred())

	agoutiDriver = agouti.ChromeDriver(agouti.ChromeOptions("args", []string{
		// "--headless",
		"--no-sandbox",
		"--disable-gpu"}))
	Expect(agoutiDriver.Start()).To(Succeed())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
	Expect(agoutiDriver.Stop()).To(Succeed())
})

func execBin(cmd *exec.Cmd) *gexec.Session {
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}

func getEnvOrError(env string) string {
	value := os.Getenv(env)
	if value == "" {
		Fail(fmt.Sprintf("Environment variable '%s' must be set", env))
	}

	return value
}

func healthcheck() error {
	resp, err := http.Get("http://localhost:3000/healthcheck")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	return fmt.Errorf("expected status code 200 but got %d", resp.StatusCode)
}

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(time.Second * 5)
	RunSpecs(t, "Acceptance Suite")
}
