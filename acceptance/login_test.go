package acceptance_test

import (
	"fmt"
	"net/http"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Logging In", func() {
	var (
		session *gexec.Session
	)

	BeforeEach(func() {
		cmd := exec.Command(serverPath)
		session = execBin(cmd)
	})

	AfterEach(func() {
		session.Kill().Wait()
	})

	FWhen("the client has previously been authorised", func() {
		// user -> login-server
		// login-server redirecting the user to oauth provider
		// oauth provider redirects the user to the login-server
		// login-server requests details from the oauth provider (e.g. userid, email)
		// login-server returns a token to the user

		// test = user
		// 1. Send a request to login-server
		// 2. Redirect to oauth provider
		// 3. Redirect to the login-server
		// 4. Receive a token

		It("gets a token", func() {

			loginServerURL := "https://stories.storyscriptapp.com/github"
			state := "random-id-for-now"

			loginURL := fmt.Sprintf("%s?state=%s", loginServerURL, state)

			resp, err := http.Get(loginURL)
			Expect(err).NotTo(HaveOccurred())

			fmt.Println(resp.Header)

			Expect(resp.StatusCode).To(Equal(302))

		})
		//        request redirect url: "https://github.com/login/oauth/authorize" query: {"scope": "user:email", "state": state, "client_id": app.secrets.github_client_id, "redirect_uri": redirect_url}

	})
})
