package mock

import "github.com/pinub/pinub"

// Ensure UserService implements pinub.UserService
var _ pinub.UserService = &UserService{}

type UserService struct {
	Client *Client

	UserFn      func(string) (*pinub.User, error)
	UserInvoked bool

	UserByTokenFn      func(string) (*pinub.User, error)
	UserByTokenInvoked bool

	CreateUserFn      func(*pinub.User) error
	CreateUserInvoked bool

	UpdateUserFn      func(*pinub.User) error
	UpdateUserInvoked bool

	DeleteUserFn      func(*pinub.User) error
	DeleteUserInvoked bool

	AddTokenFn      func(*pinub.User) error
	AddTokenInvoked bool

	RefreshTokenFn      func(*pinub.User) error
	RefreshTokenInvoked bool
}

func (s *UserService) User(email string) (*pinub.User, error) {
	s.UserInvoked = true
	return s.UserFn(email)
}

func (s *UserService) UserByToken(token string) (*pinub.User, error) {
	s.UserByTokenInvoked = true
	return s.UserByTokenFn(token)
}

func (s *UserService) CreateUser(user *pinub.User) error {
	s.CreateUserInvoked = true
	return s.CreateUserFn(user)
}

func (s *UserService) UpdateUser(user *pinub.User) error {
	s.UpdateUserInvoked = true
	return s.UpdateUserFn(user)
}

func (s *UserService) DeleteUser(user *pinub.User) error {
	s.DeleteUserInvoked = true
	return s.DeleteUserFn(user)
}

func (s *UserService) AddToken(user *pinub.User) error {
	s.AddTokenInvoked = true
	return s.AddTokenFn(user)
}

func (s *UserService) RefreshToken(user *pinub.User) error {
	s.RefreshTokenInvoked = true
	return s.RefreshTokenFn(user)
}
