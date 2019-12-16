package jwt_test

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/storyscript/auth/jwt"
)

var _ = Describe("JWT Generator", func() {

	It("generates a verifiable JWT containing the owner uuid", func() {
		generator := Generator{
			SigningKey: "secret",
		}
		tokenString, _ := generator.Generate("fake-owner-uuid")

		var claims StoryscriptClaims
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(token.Valid).To(BeTrue())

		Expect(claims.Issuer).To(Equal("storyscript"))
		Expect(claims.IssuedAt).To(BeNumerically("~", time.Now().UTC().Unix(), 100))
		Expect(claims.ExpiresAt).To(BeNumerically("~", time.Now().Add(60*60*24*365).UTC().Unix(), 100))
		Expect(claims.OwnerUUID).To(Equal("fake-owner-uuid"))
	})
})
