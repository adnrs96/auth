package http_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/storyscript/login/http"
	"github.com/storyscript/login/http/httpfakes"
)

var _ = Describe("The auth handlers", func() {

	var tokenProvider *httpfakes.FakeTokenProvider

	Describe("The login request handler", func() {

		BeforeEach(func() {
			tokenProvider = &httpfakes.FakeTokenProvider{}
			tokenProvider.GetConsentURLReturns("https://fake-consent-url.com")
		})

		It("redirects to the consent URL", func() {
			recorder := httptest.NewRecorder()
			handler := LoginHandler{
				TokenProvider: tokenProvider,
			}

			request, err := http.NewRequest("GET", "/login", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusFound))
			Expect(recorder.Header()).To(HaveKeyWithValue("Location", []string{"https://fake-consent-url.com"}))
		})
	})
})
