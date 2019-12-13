package postgres

import (
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/storyscript/login"
)

type Client struct {
	DB *sql.DB
}

func (c Client) Save(user login.User) (string, string, error) {
	_, err := c.DB.Exec("SELECT create_owner_by_login($1, $2, $3, $4, $5, $6)",
		user.Service,
		strconv.Itoa(user.ServiceID),
		user.Username,
		user.Name,
		user.Email,
		user.OAuthToken)
	if err != nil {
		panic(err)
	}

	ownerUUID, err := c.GetOwnerUUIDByEmail(user.Email)
	if err != nil {
		panic(err)
	}

	tokenUUID, err := c.GetTokenUUIDByOwnerUUID(ownerUUID)
	if err != nil {
		panic(err)
	}

	return ownerUUID, tokenUUID, nil
}

func (c Client) GetOwnerUUIDByEmail(email string) (string, error) {
	return c.getRow("SELECT owner_uuid FROM owner_emails WHERE email = $1;", email)
}

func (c Client) GetTokenUUIDByOwnerUUID(ownerUUID string) (string, error) {
	return c.getRow("SELECT uuid FROM owner_vcs WHERE owner_uuid = $1;", ownerUUID)
}

func (c Client) getRow(query string, param string) (string, error) {
	row := c.DB.QueryRow(query, param)

	var result string
	if err := row.Scan(&result); err != nil {
		return "", err
	}

	return result, nil
}
