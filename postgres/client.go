package postgres

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/storyscript/login"
)

type Client struct {
	DB *sql.DB
}

func (c Client) Save(user login.User) error {
	fmt.Println(user)
	_, err := c.DB.Exec("SELECT create_owner_by_login($1, $2, $3, $4, $5, $6)",
		user.Service,
		strconv.Itoa(user.ServiceID),
		user.Username,
		user.Name,
		user.Email,
		user.OAuthToken)
	if err != nil {
		panic(err)
		// return err
	}

	return nil
}
