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
	var (
		ownerUUID string
		token     string
	)

	// Get uuid of the owner relating to the test email used
	row := d.client.QueryRow("SELECT owner_uuid FROM owner_emails WHERE email = $1;", email)

	err := row.Scan(&ownerUUID)
	if err == sql.ErrNoRows {
		return ""
	}
	Expect(err).NotTo(HaveOccurred())

	// Get the secret token for that owner
	row = d.client.QueryRow("SELECT secret FROM token_secrets INNER JOIN tokens ON token_secrets.token_uuid = tokens.uuid WHERE tokens.owner_uuid = $1;", ownerUUID)

	err = row.Scan(&token)
	if err == sql.ErrNoRows {
		return ""
	}
	Expect(err).NotTo(HaveOccurred())

	return token
}

func (d Database) PurgeOwnerByEmail(email string) {
	var (
		ownerUUID    string
		ownerVCSUUID string
		tokenUUID    string
	)

	row := d.client.QueryRow("SELECT owner_uuid FROM owner_emails WHERE email = $1;", email)
	Expect(row.Scan(&ownerUUID)).To(Succeed())

	row = d.client.QueryRow("SELECT uuid FROM owner_vcs WHERE owner_uuid = $1;", ownerUUID)
	Expect(row.Scan(&ownerVCSUUID)).To(Succeed())

	row = d.client.QueryRow("SELECT uuid FROM tokens WHERE owner_uuid = $1;", ownerUUID)
	Expect(row.Scan(&tokenUUID)).To(Succeed())

	var err error
	_, err = d.client.Exec("DELETE FROM tokens WHERE uuid = $1;", tokenUUID)
	Expect(err).NotTo(HaveOccurred())

	_, err = d.client.Exec("DELETE FROM owner_vcs WHERE uuid = $1;", ownerVCSUUID)
	Expect(err).NotTo(HaveOccurred())

	_, err = d.client.Exec("DELETE FROM owners WHERE uuid = $1;", ownerUUID)
	Expect(err).NotTo(HaveOccurred())
}
