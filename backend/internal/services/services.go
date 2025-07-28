package services

var (
	AuthService           = &AuthenticationService{}
	AuthorizeService      = &AuthorizationService{}
	TokenServiceInstance  = &TokenService{}
	UserServiceInstance   = &UserService{}
	OAuthClientServiceIns = &OAuthClientService{}
	LogServiceInstance    = &LogService{}
	SuspendedUserServiceI = &SuspendedUserService{}
)
