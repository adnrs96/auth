package acceptance_test

import (
	"fmt"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Healthchecking", func() {
	var (
		session *gexec.Session
	)

	BeforeEach(func() {
		cmd := exec.Command(serverPath)
		cmd.Env = append(os.Environ(), fmt.Sprintf("DB_CONNECTION_STRING=%s", dbConnStr))
		session = execBin(cmd)
	})

	AfterEach(func() {
		session.Kill().Wait()
	})

	It("eventually responds 200 OK", func() {
		Eventually(healthcheck).Should(Succeed())
	})
})
