package login

type User struct {
	Service   string
	ServiceID int

	Username   string
	Name       string
	Email      string
	OAuthToken string
}
