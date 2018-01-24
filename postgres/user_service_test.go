package postgres_test

import (
	"testing"

	"github.com/pinub/pinub"
)

const userEmail = "test@example.de"
const userPassword = "test"

func createUser(s pinub.UserService) *pinub.User {
	u := pinub.User{
		Email:    userEmail,
		Password: userPassword,
	}
	if err := s.CreateUser(&u); err != nil {
		panic(err)
	}

	return &u
}

func TestUserService_CreateUser_DeleteUser(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	u := createUser(s)

	equals(t, u.Password, userPassword)
	equals(t, u.Email, userEmail)

	ok(t, s.DeleteUser(u))
}

func TestUserService_Update(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	u := createUser(s)

	u.Email = "another@email.com"
	ok(t, s.UpdateUser(u))
	equals(t, u.Email, "another@email.com")

	u.Password = "new password"
	ok(t, s.UpdateUser(u))
	equals(t, u.Password, "new password")

	ok(t, s.DeleteUser(u))
}

func TestUserService_User(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	createUser(s)

	user, err := s.User(userEmail)
	ok(t, err)

	equals(t, user.Password, userPassword)
	equals(t, user.Email, userEmail)

	ok(t, s.DeleteUser(user))
}

func TestUserService_DuplicateEmail(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	u1 := createUser(s)
	u2 := createUser(s)
	equals(t, u1.ID, u2.ID)

	ok(t, s.DeleteUser(u2))
}

func TestUserService_AddToken_ByToken_RefreshToken(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	u := createUser(s)
	ok(t, s.AddToken(u))

	equals(t, len(u.Token), len("747dacda-c33d-4647-a7f8-0cb87a5b4a24"))

	u1, err := s.UserByToken(u.Token)
	ok(t, err)
	equals(t, u, u1)
	ok(t, s.RefreshToken(u))

	ok(t, s.DeleteUser(u))
}
