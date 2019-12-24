package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tquach/prestige-api/platform/auth"
	"github.com/tquach/prestige-api/platform/config"
	"github.com/tquach/prestige-api/platform/httpserver"
	"github.com/tquach/prestige-api/platform/logger"
	"github.com/tquach/prestige-api/service-layer/accounts"
	"github.com/tquach/prestige-api/service-layer/socialauth"
	"github.com/tquach/prestige-api/service-layer/trips"
	"github.com/tquach/prestige-api/service-layer/users"
	pg "gopkg.in/pg.v4"
	"gopkg.in/yaml.v2"
)

// Packgae global constants
const (
	SecretKeyEnv = "SECRET_KEY"
)

// Flags to pass in from the command line.
var (
	hostname    = flag.String("hostname", "localhost:9000", "Hostname and port to bind to.")
	databaseURL = flag.String("databaseURL", "postgres://localhost:5432", "url of database, eg. postgres://localhost:5432")
	configFile  = flag.String("configFile", os.Getenv("CONFIG_FILE"), "Configuration file. Optional.")
)

func main() {
	settings := parseConfig()
	opts := &httpserver.Opts{
		Hostname: settings.AppSettings.Hostname,
		Port:     settings.AppSettings.Port,
	}

	s, err := httpserver.NewAppServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Opening connection to", settings.Database.URL)
	db := pg.Connect(&pg.Options{
		Addr:     settings.Database.URL,
		User:     settings.Database.Username,
		Password: settings.Database.Password,
		SSL:      settings.Database.SSLMode,
		Database: settings.Database.DBName,
	})

	tokenSvc := auth.NewTokenGenerator(settings.SecretKey)
	defer db.Close()

	// Register services
	if err := registerServices(db, tokenSvc, s); err != nil {
		log.Fatalf("failed to register services: %s\n", err)
	}

	if settings.Debug {
		pg.SetQueryLogger(logger.NewStdlogAdapter("pg"))
	}

	log.Fatal(s.Start())
}

func registerServices(db *pg.DB, tokenSvc auth.TokenService, appServer *httpserver.AppServer) error {
	// Social Authentication
	socialAuthSvc := socialauth.NewService(db, tokenSvc)
	if err := appServer.RegisterService("/auth/", socialAuthSvc); err != nil {
		return fmt.Errorf("register auth service: %s", err)
	}

	// Account Management
	accountsSvc := accounts.NewService(db, tokenSvc)
	if err := appServer.RegisterService("/accounts/", accountsSvc); err != nil {
		return fmt.Errorf("register accounts service: %s", err)
	}

	// User profiles
	userSvc := users.NewService(db, tokenSvc)
	if err := appServer.RegisterService("/users/", userSvc); err != nil {
		return fmt.Errorf("register user service: %s", err)
	}

	// Trips & Destinations
	tripsSvc := trips.NewService(db, tokenSvc)
	if err := appServer.RegisterService("/trips/", tripsSvc); err != nil {
		return fmt.Errorf("register trips service: %s", err)
	}
	return nil
}

func parseConfig() config.EnvSettings {
	flag.Parse()
	urlParts := strings.Split(*hostname, ":")
	port, err := strconv.Atoi(urlParts[1])
	if err != nil {
		log.Fatal(err)
	}

	settings := config.EnvSettings{
		SecretKey: os.Getenv(SecretKeyEnv),
		Database: config.DatabaseSettings{
			Driver: "postgres",
			URL:    *databaseURL,
		},
		AppSettings: config.AppSettings{
			Hostname: urlParts[0],
			Port:     port,
		},
	}

	// config file will override other command line args
	if *configFile != "" {
		log.Printf("Using settings from %q\n", *configFile)
		contents, err := ioutil.ReadFile(*configFile)
		if err != nil {
			log.Fatal(err)
		}

		if err = yaml.Unmarshal(contents, &settings); err != nil {
			log.Fatal(err)
		}
	}
	return settings
}
