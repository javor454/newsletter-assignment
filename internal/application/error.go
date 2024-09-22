package application

import "errors"

var (
	InvalidPasswordError              = errors.New("invalid password")
	UserNotFoundError                 = errors.New("user not found")
	NewsletterNotFoundError           = errors.New("newsletter not found")
	EmailTakenError                   = errors.New("email already taken")
	AlreadySubscibedToNewsletterError = errors.New("already subscibed to newsletter")
	UnknownUserError                  = errors.New("unknown user")
	InvalidUUIDError                  = errors.New("invalid uuid")
	InvalidTokenError                 = errors.New("invalid token")
)
