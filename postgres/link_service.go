package postgres

import (
	"database/sql"

	"github.com/pinub/pinub"
)

// Ensure LinkService implements pinub.LinkService
var _ pinub.LinkService = &LinkService{}

// LinkService represents a service for managing links.
type LinkService struct {
	Client *Client
}

// Links returns all links.
func (s *LinkService) Links(u *pinub.User) ([]pinub.Link, error) {
	query := `
		SELECT id, url, ul.created_at FROM links AS l
			JOIN user_links AS ul ON l.id = ul.link_id AND ul.user_id = $1
		ORDER BY ul.created_at DESC
	`
	rows, err := s.Client.Query(query, u.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // nolint: errcheck

	var links []pinub.Link
	for rows.Next() {
		var link pinub.Link
		if err = rows.Scan(&link.ID, &link.URL, &link.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

// Link returns a link by url.
//func (s *LinkService) Link(url string) (*pinub.Link, error) {
//	var l pinub.Link
//
//	query := "SELECT id, url FROM links WHERE url = $1 LIMIT 1"
//	err := s.Client.QueryRow(query, url).Scan(&l.ID, &l.URL)
//	if err != nil {
//		return nil, err
//	}
//
//	return &l, nil
//}

// CreateLink creates a new link for the given user. Four steps are necessary:
//   1. check if link exists in table
//   2. if no - create it
//   3. check if user has a relation to link
//   4. if no - create it
// We can make a shortcut when link does not exists, we create it and also
// create the relation to the user.
func (s *LinkService) CreateLink(l *pinub.Link, u *pinub.User) error {
	query := "SELECT id FROM links WHERE url = $1 LIMIT 1"
	if err := s.Client.QueryRow(query, l.URL).Scan(&l.ID); err == sql.ErrNoRows {
		query = "INSERT INTO links (url) VALUES ($1) RETURNING id"
		if err = s.Client.QueryRow(query, l.URL).Scan(&l.ID); err != nil {
			return err
		}
	}

	query = "SELECT created_at FROM user_links WHERE link_id = $1 AND user_id = $2 LIMIT 1"
	if err := s.Client.QueryRow(query, l.ID, u.ID).Scan(&l.CreatedAt); err == sql.ErrNoRows {
		query = "INSERT INTO user_links (link_id, user_id) VALUES ($1, $2) RETURNING created_at"
		if err = s.Client.QueryRow(query, l.ID, u.ID).Scan(&l.CreatedAt); err != nil {
			return err
		}
	}

	return nil
}

// DeleteLink remove the link for the given user. Also removes the link
// when no user stored that link anymore.
func (s *LinkService) DeleteLink(l *pinub.Link, u *pinub.User) error {
	query := "DELETE FROM user_links WHERE user_id = $1 AND link_id = $2"
	if _, err := s.Client.Exec(query, u.ID, l.ID); err != nil {
		return err
	}

	query = `
		DELETE FROM links WHERE id = $1 AND
		 (SELECT count(link_id) FROM user_links WHERE link_id = $1) = 0
	`
	_, err := s.Client.Exec(query, l.ID)
	return err
}
