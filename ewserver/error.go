package ewserver

// Error represents any constant error string we can return
type Error string

func (e Error) Error() string {
	return string(e)
}

// common errors
const (
	ErrUserNotFound      = Error("user not found")
	ErrInvalidUser       = Error("invalid username or fields")
	ErrUserAlreadyExists = Error("user already exists")
	ErrInvalidPassword   = Error("invalid password for user")
)
