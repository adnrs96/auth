package helpers

import (
	"database/sql"

	_ "github.com/lib/pq"
	. "github.com/onsi/gomega"
)

type Database struct {
	client *sql.DB
}

func NewDB(connectionString string) Database {
	pg, err := sql.Open("postgres", connectionString)
	Expect(err).NotTo(HaveOccurred())

	return Database{client: pg}
}

func (d Database) GetEmails() []string {
	rows, err := d.client.Query("SELECT email FROM owner_emails;")
	Expect(err).NotTo(HaveOccurred())
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		Expect(rows.Scan(&email)).To(Succeed())
		emails = append(emails, email)
	}

	Expect(rows.Err()).NotTo(HaveOccurred())

	return emails
}

func (d Database) GetTokenByEmail(email string) string {
	ownerUUID := d.GetOwnerUUIDByEmail(email)
	return d.getRow("SELECT secret FROM token_secrets INNER JOIN tokens ON token_secrets.token_uuid = tokens.uuid WHERE tokens.owner_uuid = $1;", ownerUUID)
}

func (d Database) GetOwnerUUIDByEmail(email string) string {
	return d.getRow("SELECT owner_uuid FROM owner_emails WHERE email = $1;", email)
}

func (d Database) GetTokenUUIDByOwnerUUID(ownerUUID string) string {
	return d.getRow("SELECT uuid FROM owner_vcs WHERE owner_uuid = $1;", ownerUUID)
}

func (d Database) PurgeOwnerByEmail(email string) {
	ownerUUID := d.GetOwnerUUIDByEmail(email)
	tokenUUID := d.GetTokenUUIDByOwnerUUID(ownerUUID)
	ownerVCSUUID := d.getRow("SELECT uuid FROM owner_vcs WHERE owner_uuid = $1;", ownerUUID)

	var err error
	_, err = d.client.Exec("DELETE FROM tokens WHERE uuid = $1;", tokenUUID)
	Expect(err).NotTo(HaveOccurred())

	_, err = d.client.Exec("DELETE FROM owner_vcs WHERE uuid = $1;", ownerVCSUUID)
	Expect(err).NotTo(HaveOccurred())

	_, err = d.client.Exec("DELETE FROM owners WHERE uuid = $1;", ownerUUID)
	Expect(err).NotTo(HaveOccurred())
}

func (d Database) getRow(query string, param string) string {
	row := d.client.QueryRow(query, param)

	var result string
	err := row.Scan(&result)
	if err == sql.ErrNoRows {
		return ""
	}
	Expect(err).NotTo(HaveOccurred())

	return result
}
