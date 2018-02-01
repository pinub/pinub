package postgres_test

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pinub/pinub"
	"github.com/pinub/pinub/claim"
	"github.com/pinub/pinub/postgres"
)

const (
	email    = "test@example.de"
	password = "test"
)

func TestUserService(t *testing.T) {
	t.Parallel()

	client, err := postgres.New(os.Getenv("DATABASE_URL"))
	claim.Ok(t, err)

	user := pinub.User{
		Email:    email,
		Password: password,
	}
	s := client.UserService()

	t.Run("create user", func(t *testing.T) {
		claim.Ok(t, s.CreateUser(&user))

		claim.Equals(t, email, user.Email)
		claim.Equals(t, password, user.Password)

		claim.Ok(t, s.DeleteUser(&user))
	})

	t.Run("update user", func(t *testing.T) {
		claim.Ok(t, s.CreateUser(&user))

		newEmail := "another@email.com"
		user.Email = newEmail
		claim.Ok(t, s.UpdateUser(&user))
		claim.Equals(t, newEmail, user.Email)

		newPassword := "another password"
		user.Password = newPassword
		claim.Ok(t, s.UpdateUser(&user))
		claim.Equals(t, newPassword, user.Password)

		user.Email = email
		user.Password = password

		claim.Ok(t, s.DeleteUser(&user))
	})

	t.Run("get user", func(t *testing.T) {
		claim.Ok(t, s.CreateUser(&user))

		u, err := s.User(email)
		claim.Ok(t, err)

		claim.Equals(t, email, u.Email)
		claim.Equals(t, password, u.Password)

		claim.Ok(t, s.DeleteUser(&user))
	})

	t.Run("duplicate email", func(t *testing.T) {
		claim.Ok(t, s.CreateUser(&user))

		user2 := pinub.User{
			Email:    email,
			Password: "bogous",
		}
		s.CreateUser(&user2)

		claim.Equals(t, user.ID, user2.ID)

		claim.Ok(t, s.DeleteUser(&user))
	})

	t.Run("add token", func(t *testing.T) {
		claim.Ok(t, s.CreateUser(&user))
		claim.Ok(t, s.AddToken(&user))

		claim.Equals(t, len(user.Token), len("747dacda-c33d-4647-a7f8-0cb87a5b4a24"))

		user2, err := s.UserByToken(user.Token)
		claim.Ok(t, err)
		claim.Equals(t, user.ID, user2.ID)
		claim.Ok(t, s.RefreshToken(&user))

		claim.Ok(t, s.DeleteUser(&user))
	})
}
