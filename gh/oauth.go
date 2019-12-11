package gh

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

type OAuthClient struct {
	Config *oauth2.Config
}

func (c OAuthClient) GetConsentURL(state string) string {
	return c.Config.AuthCodeURL(state)
}

func (c OAuthClient) GetAccessToken(authCode string) (string, error) {
	token, err := c.Config.Exchange(context.Background(), authCode)
	if err != nil {
		return "", fmt.Errorf("code exchange failed: %s", err.Error())
	}

	return token.AccessToken, nil
}
