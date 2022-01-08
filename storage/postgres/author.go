package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	pb "github.com/abdullohsattorov/catalog-service/genproto/catalog_service"
)

type authorRepo struct {
	db *sqlx.DB
}

// NewAuthorRepo ...
func NewAuthorRepo(db *sqlx.DB) *authorRepo {
	return &authorRepo{db: db}
}

func (r *authorRepo) CreateAuthor(author pb.Author) (pb.Author, error) {
	var id string

	err := r.db.QueryRow(`
				   INSERT INTO authors (author_id, name, created_at, updated_at) 
				   VALUES ($1, $2, $3, $4) RETURNING author_id `,
		author.AuthorId, author.Name, time.Now().UTC(), time.Now().UTC()).
		Scan(&id)
	if err != nil {
		return pb.Author{}, err
	}

	NewAuthor, err := r.GetAuthor(id)
	if err != nil {
		return pb.Author{}, err
	}

	return NewAuthor, nil
}

func (r *authorRepo) GetAuthor(id string) (pb.Author, error) {
	var NewAuthor pb.Author

	err := r.db.QueryRow(`
						SELECT author_id, name, created_at, updated_at FROM authors 
						WHERE author_id = $1 AND deleted_at IS NULL`, id).
		Scan(&NewAuthor.AuthorId, &NewAuthor.Name, &NewAuthor.CreatedAt, &NewAuthor.UpdatedAt)
	if err != nil {
		return pb.Author{}, err
	}

	return NewAuthor, nil
}

func (r *authorRepo) ListAuthor(page, limit int64) ([]*pb.Author, int64, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Queryx(`
				SELECT author_id, name, created_at, updated_at FROM authors 
				WHERE deleted_at is NULL ORDER BY author_id LIMIT $1 OFFSET $2
				`, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var (
		authors []*pb.Author
		count   int64
	)

	for rows.Next() {
		var author pb.Author
		err = rows.Scan(&author.AuthorId, &author.Name, &author.CreatedAt, &author.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		authors = append(authors, &author)
	}

	err = r.db.QueryRow("SELECT count(*) FROM authors WHERE deleted_at IS NULL").Scan(&count)
	if err != nil {
		return nil, 0, err
	}
	return authors, count, nil
}

func (r *authorRepo) UpdateAuthor(update pb.Author) (pb.Author, error) {
	result, err := r.db.Exec("UPDATE authors SET name = $2, updated_at = $3 WHERE author_id = $1",
		update.AuthorId, update.Name, time.Now().UTC())
	if err != nil {
		return pb.Author{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Author{}, sql.ErrNoRows
	}

	var NewAuthor pb.Author

	NewAuthor, err = r.GetAuthor(update.AuthorId)

	if err != nil {
		return pb.Author{}, err
	}
	return NewAuthor, nil
}

func (r *authorRepo) DeleteAuthor(id string) error {
	result, err := r.db.Exec("UPDATE authors SET deleted_at = $2 WHERE author_id = $1", id, time.Now().UTC())
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
