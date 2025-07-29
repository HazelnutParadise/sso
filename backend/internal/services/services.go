package services

var (
	AuthorizeService     = &authorizationService{}
	TokenService         = &tokenService{}
	UserService          = &userService{}
	OAuthClientService   = &oauthClientService{}
	LogService           = &logService{}
	SuspendedUserService = &suspendedUserService{}
)
