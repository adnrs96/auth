package gh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/storyscript/auth"
)

type UserClient struct{}

func (c UserClient) GetUser(accessToken string) (auth.User, error) {
	user, err := getUser(accessToken)
	if err != nil {
		return auth.User{}, err
	}

	email, err := getPrimaryEmail(accessToken)
	if err != nil {
		return auth.User{}, err
	}

	return auth.User{
		Service:   "github",
		ServiceID: user.ID,
		Username:  user.Login,
		Email:     email,
	}, nil
}

func getUser(accessToken string) (user, error) {
	body, err := getWithAuth("https://api.github.com/user", accessToken)
	if err != nil {
		return user{}, errors.Wrap(err, "failed to get user details")
	}

	var u user
	if err := json.Unmarshal(body, &u); err != nil {
		return user{}, errors.Wrap(err, "failed to unmarshal response: %s")
	}

	return u, nil
}

func getPrimaryEmail(accessToken string) (string, error) {
	body, err := getWithAuth("https://api.github.com/user/emails", accessToken)
	if err != nil {
		return "", errors.Wrap(err, "failed to get emails")
	}

	var emails userEmails
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal response: %s")
	}

	return emails[0].Email, nil
}

func getWithAuth(url string, accessToken string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", accessToken))

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type user struct {
	ID int `json:"id"`

	Login string `json:"login"`
	Name  string `json:"name"`
}

type userEmails []userEmail

type userEmail struct {
	Email string `json:"email"`
}
