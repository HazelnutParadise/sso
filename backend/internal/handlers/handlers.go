package handlers

import "log"

var (
	AuthHandlers          *authHandler
	OAuthProviderHandlers = &oauthProviderHandler{}
)

func init() {
	var err error
	AuthHandlers, err = NewAuthHandler()
	if err != nil {
		log.Fatal("Failed to initialize auth handler: " + err.Error())
	}
}
