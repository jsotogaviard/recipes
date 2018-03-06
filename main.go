package main

import (
	"jsotogaviard-api-test/application"
	"jsotogaviard-api-test/configuration"
	"jsotogaviard-api-test/application/security"
	"jsotogaviard-api-test/application/constants"
)

func main() {

	// Get configuration
	c := configuration.GetConfig()

	// Get security
	s := security.GetSecurity()

	// Get application
	a := application.GetApplication(c, s)

	// Run it
	a.Run(constants.GetDoubleDot() + c.Port)
}
