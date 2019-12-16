package postgres_test

import (
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/storyscript/login"
	"github.com/storyscript/login/acceptance/helpers"
	. "github.com/storyscript/login/postgres"
)

var _ = Describe("Postgres Client", func() {

	var dbHelper helpers.Database

	BeforeEach(func() {
		dbHelper = helpers.NewDB(dbConnStr)
	})

	Describe("Saving a User", func() {

		var savedOwnerUUID string

		BeforeEach(func() {
			client := Client{
				DB: openDB(),
			}

			userToSave := login.User{
				Service:   "github",
				ServiceID: 123,

				Username:   "test-username",
				Name:       "test-name",
				Email:      "test-email@example.com",
				OAuthToken: "test-token",
			}

			var err error
			savedOwnerUUID, err = client.Save(userToSave)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			dbHelper.PurgeOwnerByEmail("test-email@example.com")
		})

		It("persists the user to postgres and returns a uuid for the owner", func() {
			fetchedOwnerUUID := dbHelper.GetOwnerUUIDByEmail("test-email@example.com")
			Expect(savedOwnerUUID).To(Equal(fetchedOwnerUUID))
		})

		When("saving the same user twice", func() {
			It("returns the same owner uuid", func() {
				client := Client{
					DB: openDB(),
				}

				userToSave := login.User{
					Service:   "github",
					ServiceID: 123,

					Username:   "test-username",
					Name:       "test-name",
					Email:      "test-email@example.com",
					OAuthToken: "test-token",
				}

				secondSavedOwnerUUID, err := client.Save(userToSave)
				Expect(err).NotTo(HaveOccurred())
				Expect(secondSavedOwnerUUID).To(Equal(savedOwnerUUID))
			})
		})
	})
})

func openDB() *sql.DB {
	dbDriver := "postgres"
	db, err := sql.Open(dbDriver, dbConnStr)
	Expect(err).NotTo(HaveOccurred())

	return db
}
