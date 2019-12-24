package idp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tquach/prestige-api/platform/logger"
)

// Ensure interface compliance.
var _ IdentityProvider = &FacebookIdentityProvider{}

// Constants for reuse
const (
	FbValidationURL = "https://graph.facebook.com/me?access_token=%s"
)

// FacebookIdentityProvider is a Facebook implementation of the IdentityProvider.
type FacebookIdentityProvider struct {
	client http.Client
	logger logger.Logger
}

// ValidateToken will check Facebook and ensure the token is valid.
func (f *FacebookIdentityProvider) ValidateToken(oauthResponse OAuthResponse) (IdentityProviderResponse, error) {
	client := http.Client{}
	idpResp := IdentityProviderResponse{}

	resp, err := client.Get(fmt.Sprintf(FbValidationURL, oauthResponse.AccessToken))
	if err != nil {
		return idpResp, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return idpResp, err
	}

	idpResp.RawResponse = contents

	f.logger.Debugf("received response from Facebook: %s", string(contents))
	if err := json.Unmarshal(contents, &idpResp); err != nil {
		return idpResp, err
	}

	if idpResp.Error != nil {
		return idpResp, idpResp.Error
	}

	if idpResp.ID != oauthResponse.UserID {
		return idpResp, fmt.Errorf("invalid token for requested user %q != %q", oauthResponse.UserID, idpResp.ID)
	}

	return idpResp, nil
}
