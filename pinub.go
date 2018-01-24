package pinub

// LinkService interface takes care of the database connections related to
// the Link.
type LinkService interface {
	//Link(string) (*Link, error)
	Links(*User) ([]Link, error)
	CreateLink(*Link, *User) error
	DeleteLink(*Link, *User) error
}

// UserService interface takes care of all database related connections with
// the user object.
type UserService interface {
	User(string) (*User, error)
	UserByToken(string) (*User, error)
	CreateUser(*User) error
	UpdateUser(*User) error
	DeleteUser(*User) error
	AddToken(*User) error
	RefreshToken(*User) error
}

// Client interface has all the connections to the services.
type Client interface {
	LinkService() LinkService
	UserService() UserService
}
