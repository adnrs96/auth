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

			FIt("sets a cookie containing a JWT token", func() {
				cookies, err := page.GetCookies()
				Expect(err).NotTo(HaveOccurred())
				cookie := cookies[0]

				Expect(cookie.Name).To(Equal("storyscript-access-token"))
				Expect(cookie.Path).To(Equal("/"))
				Expect(cookie.Expires).To(BeTemporally("~", time.Now().Add(time.Hour*24*365), time.Minute))
				Expect(cookie.HttpOnly).To(BeTrue())

				var claims StoryscriptClaims
				token, err := jwt.ParseWithClaims(cookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("SECRET_KEY")), nil
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(token.Valid).To(BeTrue())

				Expect(claims.Issuer).To(Equal("storyscript"))
				Expect(claims.OwnerUUID).To(Equal(db.GetOwnerUUIDByEmail(os.Getenv("ACCEPTANCE_EMAIL"))))
				Expect(claims.TokenUUID).To(Equal(db.GetTokenUUIDByOwnerUUID(claims.OwnerUUID)))

				// Expect(claims.Secret).To(Equal(db.GetTokenByEmail(os.Getenv("ACCEPTANCE_EMAIL"))))

				//	token :=jqt.NewWithClaims(jwt.SigningMethodHMC256)
				//	ss, err := token.SignedString
				//
				//
				//Calendar c = Calendar.getInstance();
				// c.add(Calendar.YEAR, 1);
				//
				// Algorithm algorithm = Algorithm.HMAC256(Constants.JWT_COOKIE_SECRET_KEY);
				// String token = JWT.create()
				//         .withIssuer(Constants.JWT_ISSUER)
				//         .withClaim(Constants.JWT_CLAIM_KEY_SECRET, SecretsUtil.hashTokenSecret(loginTokenSecret))
				//         .withClaim(Constants.JWT_CLAIM_KEY_OWNER_UUID, owner.getId())
				//         .withClaim(Constants.JWT_CLAIM_KEY_TOKEN_UUID, owner.getTokenUuid())
				//         .withIssuedAt(new Date())
				//         .withExpiresAt(c.getTime())
				//         .sign(algorithm);
				// final Cookie authCookie = new Cookie(Constants.JWT_COOKIE_NAME, token);
				// authCookie.setHttpOnly(true);
				// authCookie.setPath("/");
				// authCookie.setDomain(Constants.HOST);
				// authCookie.setMaxAge(60 * 60 * 24 * 365); // 1 year cookie - safety first
				// res.addCookie(authCookie);
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
	TokenUUID string `json:"token_uuid"`
}
