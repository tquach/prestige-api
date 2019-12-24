package socialauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	pg "gopkg.in/pg.v4"

	"github.com/tquach/prestige-api/platform/auth"
	"github.com/tquach/prestige-api/platform/idp"
	"github.com/tquach/prestige-api/platform/logger"
	"github.com/tquach/prestige-api/platform/render"
	"github.com/tquach/prestige-api/service-layer/users"
)

// Errors declared for reuse
var (
	ErrUnauthorized = errors.New("unauthorized")
)

// AuthenticationResponse contains the response data from the authentication request.
type AuthenticationResponse struct {
	AccessToken string     `json:"accessToken"`
	User        users.User `json:"user"`
}

// AuthenticationRequest wraps an authentication request.
type AuthenticationRequest struct {
	SocialNetwork string            `json:"socialNetwork"`
	OAuthResponse idp.OAuthResponse `json:"authResponse"`
}

// Controller contains all the attributes needed for a full controller.
type Controller struct {
	DB           *pg.DB
	logger       logger.Logger
	TokenService auth.TokenService
}

// Authenticate will handle the authentication request.
func (c *Controller) Authenticate(w http.ResponseWriter, req *http.Request) {
	authReq := AuthenticationRequest{}
	if err := json.NewDecoder(req.Body).Decode(&authReq); err != nil {
		c.logger.Errorf("invalid request: %s", err)
		render.JSONError(err, http.StatusBadRequest, w)
		return
	}

	socialAccountID := authReq.OAuthResponse.UserID
	user := users.User{}
	if err := c.DB.Model(&user).Where("social_account_id = ?", socialAccountID).Limit(1).Select(); err != nil {
		c.logger.Errorf("unrecognized social account id: %s", authReq.OAuthResponse.UserID)
		render.JSONError(fmt.Errorf("user not found %s", authReq.OAuthResponse.UserID), http.StatusNotFound, w)
		return
	}

	provider, err := idp.New(authReq.SocialNetwork)
	if err != nil {
		c.logger.Errorf("invalid social network: %s\n", err.Error())
		render.JSONError(err, http.StatusBadRequest, w)
		return
	}

	idpResponse, err := provider.ValidateToken(authReq.OAuthResponse)
	if err != nil {
		c.logger.Error("invalid oauth response: ", err)
		render.JSONError(err, http.StatusInternalServerError, w)
		return
	}

	if idpResponse.Error != nil {
		c.logger.Error("error from idp", idpResponse.Error.Message)
		render.JSON(idpResponse.Error, http.StatusUnauthorized, w)
		return
	}

	claimsPayload := map[string]interface{}{}
	claimsPayload["sid"] = authReq.OAuthResponse.UserID
	claimsPayload["uid"] = user.ID
	claimsPayload["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Create JWT token
	w.Header().Add("Content-type", "application/json")

	accessToken, err := c.TokenService.GenerateAccessToken(claimsPayload)
	if err != nil {
		c.logger.Error(err)
		render.JSONError(err, http.StatusInternalServerError, w)
		return
	}

	authResponse := AuthenticationResponse{
		AccessToken: accessToken,
		User:        user,
	}

	if err := json.NewEncoder(w).Encode(authResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// NewController creates a new instance of a Controller.
func NewController(tokenSvc auth.TokenService, db *pg.DB) *Controller {
	return &Controller{
		DB:           db,
		logger:       logger.New("auth"),
		TokenService: tokenSvc,
	}
}
