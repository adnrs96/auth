package acceptance_test

import (
	"fmt"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/sclevine/agouti"
	"github.com/storyscript/login/acceptance/helpers"
)

var _ = Describe("The Login Process", func() {
	var (
		session *gexec.Session

		page *agouti.Page
		db   helpers.Database
	)

	BeforeEach(func() {
		db = helpers.NewDB(dbConnStr)

		cmd := exec.Command(serverPath)
		cmd.Env = append(os.Environ(), fmt.Sprintf("DB_CONNECTION_STRING=%s", dbConnStr))
		session = execBin(cmd)

		var err error
		page, err = agoutiDriver.NewPage()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		db.PurgeOwnerByEmail(os.Getenv("ACCEPTANCE_EMAIL"))

		Expect(page.Destroy()).To(Succeed())
		session.Kill().Wait()
	})

	var pageURL = func() string {
		url, err := page.URL()
		Expect(err).NotTo(HaveOccurred())
		return url
	}

	When("the user has previously authorized a client", func() {
		// Ensure that the user provided in the ACCEPTANCE_EMAIL env var has already authorized the client
		// otherwise there will be an additional prompt the test doesn't handle

		When("going through the login flow", func() {

			BeforeEach(func() {
				Expect(page.Navigate("http://localhost:3000/login")).To(Succeed())

				Eventually(pageURL).Should(HavePrefix("https://github.com/login"))
				loginToGitHub(page)
			})

			It("creates a token in the database associated with the user", func() {
				Eventually(db.GetEmails).Should(HaveLen(1))
				Expect(db.GetEmails()).To(ContainElement(os.Getenv("ACCEPTANCE_EMAIL")))

				Eventually(func() string {
					return db.GetTokenByEmail(os.Getenv("ACCEPTANCE_EMAIL"))
				}).ShouldNot(BeEmpty())
			})
		})
	})
})

func loginToGitHub(page *agouti.Page) {
	userField := page.FindByName("login")
	passwordField := page.FindByName("password")
	loginButton := page.FindByName("commit")

	Expect(userField.Fill(os.Getenv("ACCEPTANCE_EMAIL"))).To(Succeed())
	Expect(passwordField.Fill(os.Getenv("ACCEPTANCE_PASSWORD"))).To(Succeed())
	Expect(loginButton.Submit()).To(Succeed())
}