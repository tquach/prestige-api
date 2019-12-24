package accounts

import (
	"github.com/tquach/prestige-api/platform/idp"

	"github.com/tquach/prestige-api/service-layer/users"
)

// Account is the domain model for account information.
type Account struct {
	SocialNetwork string            `json:"socialNetwork"`
	OAuthResponse idp.OAuthResponse `json:"authResponse"`
	User          users.User        `json:"user"`
}

// RegistrationResponse holds the response data from the registration result.
type RegistrationResponse struct {
	User        users.User `json:"user"`
	AccessToken string     `json:"accessToken"`
}
