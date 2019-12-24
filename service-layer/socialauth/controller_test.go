//+build integration

package socialauth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/tquach/prestige-api/platform/config"

	"github.com/tquach/prestige-api/platform/auth"
	"github.com/tquach/prestige-api/platform/logger"
	"github.com/tquach/prestige-api/service-layer/users"
	pg "gopkg.in/pg.v4"
)

var (
	db *pg.DB
)

func setupTest() {
	db = pg.Connect(pgOptions())
	config.Env = "local"
}

func pgOptions() *pg.Options {
	return &pg.Options{
		User:               "prestige",
		Password:           "changeme",
		Database:           "test",
		Addr:               "192.168.99.100:5432",
		DialTimeout:        30 * time.Second,
		ReadTimeout:        10 * time.Second,
		WriteTimeout:       10 * time.Second,
		PoolSize:           10,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        10 * time.Second,
		IdleCheckFrequency: 100 * time.Millisecond,
	}
}

func TestAuthenticateWithBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c := Controller{
		TokenService: auth.NewTokenGenerator("secret"),
		logger:       logger.New("auth"),
	}

	req, err := http.NewRequest("GET", "/", strings.NewReader(
		`{
			"socialNetwork": "b",
			"authResponse": {
				"userID": 132908123
			}
		}`))
	if err != nil {
		t.Fatal(err)
	}

	c.Authenticate(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatal(w.Body.String())
	}
}

func TestAuthReturnsUser(t *testing.T) {
	setupTest()

	// Create test user
	user := users.User{
		SocialAccountID: "132908123",
	}

	if err := db.Create(&user); err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = db.Delete(&user)
	}()

	w := httptest.NewRecorder()
	c := Controller{
		TokenService: auth.NewTokenGenerator("secret"),
		DB:           db,
		logger:       logger.New("test"),
	}

	req, err := http.NewRequest("POST", "/auth", strings.NewReader(
		`{
			"socialNetwork": "local",
			"authResponse": {
				"userID": "132908123",
				"accessToken": "abcdefgh"
			}
		}`))
	if err != nil {
		t.Fatal(err)
	}

	c.Authenticate(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected an OK status but got %d instead", w.Code)
	}
}
