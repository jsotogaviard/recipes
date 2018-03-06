package configuration

import (
	"os"
	"errors"
	"jsotogaviard-api-test/application/constants"
)

// Configuration configures the database and port of the server
type Config struct {
	Database   Database
	Port string
}

// The database configuration
type Database struct {
	Host     string
	User     string
	Password string
	DBName   string
	Charset  string
}

// Retrieve configuration from environment
func GetConfig() *Config {
	hostVar, ok := os.LookupEnv(constants.GetHost());
	if !ok {
		errors.New(constants.GetHost() + " is not present")
	}

	userVar, ok := os.LookupEnv(constants.GetUser());
	if !ok {
		errors.New(constants.GetUser() + " is not present")
	}

	passwordVar, ok := os.LookupEnv(constants.GetPassword());
	if !ok {
		errors.New(constants.GetPassword() + " is not present")
	}

	dbnameVar, ok := os.LookupEnv(constants.GetDbname());
	if !ok {
		errors.New(constants.GetDbname() + " is not present")
	}

	charsetVar, ok := os.LookupEnv(constants.GetCharset());
	if! ok {
		errors.New(constants.GetCharset() + " is not present")
	}

	portVar, ok := os.LookupEnv(constants.GetPort());
	if !ok {
		errors.New(constants.GetPort() + " is not present")
	}

	dbVar := Database{
		Host:     hostVar,
		User:     userVar,
		Password: passwordVar,
		DBName:   dbnameVar,
		Charset:  charsetVar,
	}

	return &Config{
		Database: dbVar,
		Port: portVar,
	}
}