package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	pb "github.com/abdullohsattorov/catalog-service/genproto"
)

type catalogRepo struct {
	db *sqlx.DB
}

// NewCatalogRepo ...
func NewCatalogRepo(db *sqlx.DB) *catalogRepo {
	return &catalogRepo{db: db}
}

func (r *catalogRepo) CreateBook(book pb.Book) (pb.Book, error) {
	var id string
	err := r.db.QueryRow(`
        INSERT INTO books(book_id, name, author_id, price, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6) returning book_id`, book.BookId, book.Name, book.AuthorId, book.Price, time.Now().UTC(), time.Now().UTC()).Scan(&id)
	if err != nil {
		return pb.Book{}, err
	}

	var NewBook pb.Book

	NewBook, err = r.GetBook(id)

	if err != nil {
		return pb.Book{}, err
	}

	return NewBook, nil
}

func (r *catalogRepo) GetBook(id string) (pb.Book, error) {
	var book pb.Book
	err := r.db.QueryRow(`
        SELECT book_id, name, author_id, price, created_at, updated_at FROM books
        WHERE book_id=$1 and deleted_at is null`, id).Scan(&book.BookId, &book.Name, &book.AuthorId, &book.Price, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		return pb.Book{}, err
	}

	return book, nil
}

func (r *catalogRepo) ListBook(page, limit int64) ([]*pb.Book, int64, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Queryx(
		`SELECT book_id, name, author_id, price, created_at, updated_at FROM books WHERE deleted_at is null order by book_id LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close() // nolint:err check

	var (
		books []*pb.Book
		count int64
	)
	for rows.Next() {
		var book pb.Book
		err = rows.Scan(&book.BookId, &book.Name, &book.AuthorId, &book.Price, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		books = append(books, &book)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM books where deleted_at is null`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return books, count, nil
}

func (r *catalogRepo) UpdateBook(book pb.Book) (pb.Book, error) {
	result, err := r.db.Exec(`UPDATE books SET name=$1, author_id=$2, price=$3, updated_at = $4 WHERE book_id=$5 and deleted_at is null`,
		book.Name, book.AuthorId, book.Price, time.Now().UTC(), book.BookId,
	)
	if err != nil {
		return pb.Book{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Book{}, sql.ErrNoRows
	}

	var NewBook pb.Book

	NewBook, err = r.GetBook(book.BookId)

	if err != nil {
		return pb.Book{}, err
	}

	return NewBook, nil
}

func (r *catalogRepo) DeleteBook(id string) error {
	result, err := r.db.Exec(`UPDATE books SET deleted_at = $1 WHERE book_id=$2`, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
