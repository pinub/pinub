package pinub_test

import (
	"testing"

	"github.com/pinub/pinub"
	"golang.org/x/crypto/bcrypt"
)

const validPassword = "1234test"
const validEmail = "test@test.test"

func TestUser_IsValid(t *testing.T) {
	t.Parallel()
	t.Run("valid emails", func(t *testing.T) {
		emails := []string{
			validEmail,
			"d@t.c",
			"test@example.com",
			"test.test@example.com",
			"t+t@t.d",
		}
		for _, email := range emails {
			u := &pinub.User{
				Email:    email,
				Password: validPassword,
			}
			assert(t, u.IsValid(), "email should be valid")
		}
	})

	t.Run("invalid emails", func(t *testing.T) {
		emails := []string{
			"@t.c",
			"test@localhost",
			"t+t@t.",
		}
		for _, email := range emails {
			u := &pinub.User{
				Email:    email,
				Password: validPassword,
			}
			assert(t, !u.IsValid(), "email should be invalid")
		}
	})

	t.Run("valid passwords", func(t *testing.T) {
		passwords := []string{
			validPassword,
			"1234",
			"abcd",
			"254qjgeZ4VGn24cbyG4axU",
			"hQKVKsnB9LnzQgQWEbadAp8m",
			"ZL6aJ56MzqwAhX8NjgTUgYGX",
			"y8PkDQ6TnKhRpy2RJjn2yL56",
			"JWBacJKPmYjqa632k4uRCUbV",
		}
		for _, password := range passwords {
			u := &pinub.User{
				Email:    validEmail,
				Password: password,
			}
			assert(t, u.IsValid(), "password should be valid")
		}
	})

	t.Run("invalid passwords", func(t *testing.T) {
		passwords := []string{
			"123",
			"abc",
		}
		for _, password := range passwords {
			u := &pinub.User{
				Email:    validEmail,
				Password: password,
			}
			assert(t, !u.IsValid(), "password should be invalid")
		}
	})
}

func TestUser_IsValidPassword(t *testing.T) {
	t.Parallel()
	hash, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
	u := &pinub.User{
		Email:    validEmail,
		Password: string(hash),
	}

	tests := map[string]bool{
		validPassword: true,
		"1234":        false,
		//		"abcd":        false,
		//		"sthinvalid":  false,
		//		"1234tes":     false,
		//		"":            false,
	}
	for password, want := range tests {
		if got := u.IsValidPassword(password); got != want {
			t.Errorf("want %v - got %v for password %v", want, got, password)
		}
	}
}
