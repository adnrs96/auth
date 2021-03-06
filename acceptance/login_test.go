package acceptance_test

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/sclevine/agouti"
	"github.com/storyscript/auth/acceptance/helpers"
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
		session.Kill().Wait()
		Expect(page.Destroy()).To(Succeed())
		db.PurgeOwnerByEmail(os.Getenv("ACCEPTANCE_EMAIL"))
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

			It("sets a cookie containing a JWT token", func() {
				cookies, err := page.GetCookies()
				Expect(err).NotTo(HaveOccurred())
				cookie := cookies[0]

				Expect(cookie.Name).To(Equal("storyscript-access-token"))
				Expect(cookie.Path).To(Equal("/"))
				Expect(cookie.Expires).To(BeTemporally("~", time.Now().Add(time.Hour*24*365), time.Minute))
				Expect(cookie.HttpOnly).To(BeTrue())
				Expect(cookie.Domain).To(Equal(os.Getenv("DOMAIN")))

				var claims StoryscriptClaims
				token, err := jwt.ParseWithClaims(cookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(token.Valid).To(BeTrue())

				Expect(claims.Issuer).To(Equal("storyscript"))
				Expect(claims.OwnerUUID).To(Equal(db.GetOwnerUUIDByEmail(os.Getenv("ACCEPTANCE_EMAIL"))))
			})

			It("redirects to the provided redirect URI after login", func() {
				Eventually(pageURL).Should(Equal(os.Getenv("POST_LOGIN_REDIRECT_URI")))
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

type StoryscriptClaims struct {
	jwt.StandardClaims
	OwnerUUID string `json:"owner_uuid"`
}
