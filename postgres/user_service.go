package postgres

import (
	"github.com/pinub/pinub"
)

// Ensure UserService implements pinub.UserService
var _ pinub.UserService = &UserService{}

// UserService represents a service for managing users.
type UserService struct {
	Client *Client
}

// User searches the users table for a user with the given email address.
func (s *UserService) User(email string) (*pinub.User, error) {
	var u pinub.User

	query := "SELECT id, email, password FROM users WHERE email = $1 LIMIT 1"
	err := s.Client.QueryRow(query, email).Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// CreateUser stores the given user object on the users table. Does not
// create a new user when the given users email address is already in the
// database.
func (s *UserService) CreateUser(u *pinub.User) error {
	query := "SELECT id FROM users WHERE email = $1 LIMIT 1"
	if err := s.Client.QueryRow(query, u.Email).Scan(&u.ID); err == nil {
		return nil
	}

	query = "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	err := s.Client.QueryRow(query, u.Email, u.Password).Scan(&u.ID)

	return err
}

// UpdateUser updates the fields of the user object.
func (s *UserService) UpdateUser(u *pinub.User) error {
	query := "UPDATE users SET (email, password) = ($1, $2) WHERE id = $3"
	_, err := s.Client.Exec(query, u.Email, u.Password, u.ID)

	return err
}

// DeleteUser removes the given user object from the users table.
func (s *UserService) DeleteUser(u *pinub.User) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := s.Client.Exec(query, u.ID)

	return err
}

// AddToken creates a new token for the given user object and sets the last
// active date to now.
func (s *UserService) AddToken(u *pinub.User) error {
	query := "INSERT INTO logins (user_id) VALUES ($1) RETURNING token, active_at"
	err := s.Client.QueryRow(query, u.ID).Scan(&u.Token, &u.ActiveAt)

	return err
}

// UserByToken queries the database for a user by the given token.
func (s *UserService) UserByToken(token string) (*pinub.User, error) {
	var u pinub.User

	query := "SELECT id, email, password, token, active_at FROM users u " +
		" JOIN logins l ON u.id = l.user_id AND l.token = $1"
	err := s.Client.QueryRow(query, token).
		Scan(&u.ID, &u.Email, &u.Password, &u.Token, &u.ActiveAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// RefreshToken updates the last seen token of the user.
func (s *UserService) RefreshToken(u *pinub.User) error {
	query := "UPDATE logins SET active_at = now() WHERE token = $1 RETURNING active_at"
	err := s.Client.QueryRow(query, u.Token).Scan(&u.ActiveAt)

	return err
}
