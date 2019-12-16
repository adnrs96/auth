package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/storyscript/login"
	. "github.com/storyscript/login/http"
	"github.com/storyscript/login/http/httpfakes"
)

var _ = Describe("The auth handlers", func() {

	var (
		tokenProvider   *httpfakes.FakeTokenProvider
		userInfoFetcher *httpfakes.FakeUserInfoFetcher
		userRepository  *httpfakes.FakeUserRepository
		tokenGenerator  *httpfakes.FakeTokenGenerator
		server          Server

		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		tokenProvider = &httpfakes.FakeTokenProvider{}
		userInfoFetcher = &httpfakes.FakeUserInfoFetcher{}
		userRepository = &httpfakes.FakeUserRepository{}
		tokenGenerator = &httpfakes.FakeTokenGenerator{}

		tokenProvider.GetAccessTokenReturns("fake-access-token", nil)
		userInfoFetcher.GetUserReturns(login.User{
			Name: "test-user-name",
		}, nil)
		userRepository.SaveReturns("fake-owner-uuid", nil)
		tokenGenerator.GenerateReturns("fake-token", nil)
		tokenProvider.GetConsentURLReturns("https://fake-consent-url.com")

		server = Server{
			TokenProvider:   tokenProvider,
			UserInfoFetcher: userInfoFetcher,
			UserRepository:  userRepository,
			TokenGenerator:  tokenGenerator,
		}

		recorder = httptest.NewRecorder()
	})

	Describe("The login request handler", func() {
		It("redirects to the consent URL", func() {
			handler := http.HandlerFunc(server.HandleLogin)

			request, err := http.NewRequest("GET", "/login", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusFound))
			Expect(recorder.Header()).To(HaveKeyWithValue("Location", []string{"https://fake-consent-url.com"}))
		})
	})

	Describe("the callback request handler", func() {

		JustBeforeEach(func() {
			handler := http.HandlerFunc(server.HandleCallback)

			formValues := url.Values{"code": {"fake-auth-code"}}
			request, err := http.NewRequest("POST", "/callback", strings.NewReader(formValues.Encode()))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			handler.ServeHTTP(recorder, request)
		})

		It("gets an access token based on the auth code", func() {
			Expect(tokenProvider.GetAccessTokenCallCount()).NotTo(BeZero())
			Expect(tokenProvider.GetAccessTokenArgsForCall(0)).To(Equal("fake-auth-code"))
		})

		It("uses the access token to fetch user info", func() {
			Expect(userInfoFetcher.GetUserCallCount()).NotTo(BeZero())
			Expect(userInfoFetcher.GetUserArgsForCall(0)).To(Equal("fake-access-token"))
		})

		It("saves the fetched user along with their access token", func() {
			Expect(userRepository.SaveCallCount()).NotTo(BeZero())
			Expect(userRepository.SaveArgsForCall(0)).To(Equal(login.User{
				Name:       "test-user-name",
				OAuthToken: "fake-access-token",
			}))
		})

		It("generates a token for the user", func() {
			Expect(tokenGenerator.GenerateCallCount()).NotTo(BeZero())
			Expect(tokenGenerator.GenerateArgsForCall(0)).To(Equal("fake-owner-uuid"))
		})

		It("sets a cookie containing the token", func() {
			cookie := recorder.Result().Cookies()[0]

			Expect(cookie.Name).To(Equal("storyscript-access-token"))
			Expect(cookie.Path).To(Equal("/"))
			Expect(cookie.Expires).To(BeTemporally("~", time.Now().Add(time.Hour*24*365), time.Minute))
			Expect(cookie.HttpOnly).To(BeTrue())
			Expect(cookie.SameSite).To(Equal(http.SameSiteStrictMode))
			Expect(cookie.Value).To(Equal("fake-token"))
		})

		When("fetching an access token fails", func() {
			BeforeEach(func() {
				tokenProvider.GetAccessTokenReturns("", errors.New("explode"))
			})

			It("returns a 400 Bad Request", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("fetching a user fails", func() {
			BeforeEach(func() {
				userInfoFetcher.GetUserReturns(login.User{}, errors.New("explode"))
			})

			It("returns a 400 Bad Request", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("saving a user fails", func() {
			BeforeEach(func() {
				userRepository.SaveReturns("", errors.New("explode"))
			})

			It("returns a 500 Internal Server Error", func() {
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("generating a token fails", func() {
			BeforeEach(func() {
				tokenGenerator.GenerateReturns("", errors.New("explode"))
			})

			It("returns a 500 Internal Server Error", func() {
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
