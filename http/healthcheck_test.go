package http_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/storyscript/login/http"
)

var _ = Describe("The healthcheck handler", func() {

	It("returns 200 OK", func() {
		recorder := httptest.NewRecorder()
		handler := HealthcheckHandler{}

		request, err := http.NewRequest("GET", "http://localhost:3000/healthcheck", nil)
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(recorder, request)

		Expect(recorder.Code).To(Equal(http.StatusOK))
	})
})
