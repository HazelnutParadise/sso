package services

var (
	AuthService          = &authenticationService{}
	AuthorizeService     = &authorizationService{}
	TokenService         = &tokenService{}
	UserService          = &userService{}
	OAuthClientService   = &oauthClientService{}
	LogService           = &logService{}
	SuspendedUserService = &suspendedUserService{}
)
