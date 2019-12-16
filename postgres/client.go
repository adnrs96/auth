package postgres

import (
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/storyscript/login"
)

type Client struct {
	DB *sql.DB
}

func (c Client) Save(user login.User) (string, error) {
	_, err := c.DB.Exec("SELECT create_owner_by_login($1, $2, $3, $4, $5, $6)",
		user.Service,
		strconv.Itoa(user.ServiceID),
		user.Username,
		user.Name,
		user.Email,
		user.OAuthToken)
	if err != nil {
		return "", errors.Wrap(err, "failed to save user")
	}

	ownerUUID, err := c.getOwnerUUIDByEmail(user.Email)
	if err != nil {
		return "", errors.Wrap(err, "failed to get ownerUUID")
	}

	return ownerUUID, nil
}

func (c Client) getOwnerUUIDByEmail(email string) (string, error) {
	return c.getRow("SELECT owner_uuid FROM owner_emails WHERE email = $1;", email)
}

func (c Client) getRow(query string, param string) (string, error) {
	row := c.DB.QueryRow(query, param)

	var result string
	if err := row.Scan(&result); err != nil {
		return "", err
	}

	return result, nil
}
