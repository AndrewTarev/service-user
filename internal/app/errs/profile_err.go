package errs

import "errors"

var (
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInvalidUserId        = errors.New("invalid user id")
	ErrProfileAlreadyExists = errors.New("user profile already exists")
	ErrCreateUserProfile    = errors.New("error create user-profile")
	ErrProfileNotFound      = errors.New("profile not found")
	ErrDeleteUserProfile    = errors.New("error delete user-profile")
	ErrUpdateUserProfile    = errors.New("error update user-profile")
)

var (
	ErrTokenExpired = errors.New("token has expired")
	ErrTokenInvalid = errors.New("token invalid")
)
