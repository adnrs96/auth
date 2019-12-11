package gh

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/storyscript/login"
)

type Client struct{}

func (c Client) GetUser(accessToken string) (login.User, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return login.User{}, fmt.Errorf("failed creating user request: %s", err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", accessToken))
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return login.User{}, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)

	var user login.User
	if err := decoder.Decode(&user); err != nil {
		return login.User{}, fmt.Errorf("could not unmarshal response: %s", err.Error())
	}

	return user, nil
}
