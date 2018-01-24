package pinub

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const emailRegexp = `.+@.+\..+`
const minPasswordLength = 4

var (
	// ErrUserInvalidEmail when the entered email address is not correct.
	ErrUserInvalidEmail = errors.New("Invalid Email")
	// ErrUserShortPassword when password is shorter than minPasswordLength.
	ErrUserShortPassword = errors.New("Password is too short")
	// ErrUserPasswordNotCorrect when the given password is not correct.
	ErrUserPasswordNotCorrect = errors.New("Email or Password is incorrect")
)

// User represents a stored user in database.
type User struct {
	ID        string
	Email     string
	Password  string
	Token     string
	CreatedAt *time.Time
	ActiveAt  *time.Time
	Errors    map[string]error
}

// IsValid checks the current set values for validity.
func (u *User) IsValid() bool {
	u.Errors = make(map[string]error)

	if matched, _ := regexp.MatchString(emailRegexp, u.Email); !matched {
		u.Errors["Email"] = ErrUserInvalidEmail
	}

	if len(u.Password) < minPasswordLength {
		u.Errors["Password"] = ErrUserShortPassword
	}

	return len(u.Errors) == 0
}

// IsValidPassword checks the given password to be valid with the set user
// password.
func (u *User) IsValidPassword(password string) bool {
	u.Errors = make(map[string]error)

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		u.Errors["Password"] = ErrUserPasswordNotCorrect
	}

	return len(u.Errors) == 0
}

// HashPassword does exactly what it is names: hashes the given password and
// either returns err when an error occured or the newly created hash.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash), err
}
