package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	pb "github.com/abdullohsattorov/catalog-service/genproto/catalog_service"
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


func (r *catalogRepo) CreateCategory(category pb.Category) (pb.Category, error) {
	var id string
	err := r.db.QueryRow(`
		INSERT INTO categories(category_id, name, parent_uuid, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) returning category_id`, category.CategoryId, category.Name, category.ParentUuid, time.Now().UTC(), time.Now().UTC()).Scan(&id)
	if err != nil {
		return pb.Category{}, err
	}

	newCategory, err := r.GetCategory(id)
	if err != nil {
		return pb.Category{}, err
	}

	return newCategory, nil
}

func (r *catalogRepo) GetCategory(id string) (pb.Category, error) {
	var category pb.Category

	err := r.db.QueryRow(`
		SELECT category_id, name, parent_uuid, created_at, updated_at FROM categories 
		WHERE category_id=$1 and deleted_at is null`, id).Scan(&category.CategoryId, &category.Name, &category.ParentUuid, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return pb.Category{}, err
	}

	return category, nil
}



func (r *catalogRepo) CreateAuthor(author pb.Author) (pb.Author, error) {
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

func (r *catalogRepo) GetAuthor(id string) (pb.Author, error) {
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

func (r *catalogRepo) ListAuthor(page, limit int64) ([]*pb.Author, int64, error) {
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
  

func (r *catalogRepo) ListCategory(page, limit int64) ([]*pb.Category, int64, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Queryx(`SELECT category_id, name, parent_uuid, created_at, updated_at FROM categories WHERE deleted_at is null order by category_id LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, 0, err
	}
  
  var(
		categories []*pb.Category
		count      int64
	)

	for rows.Next() {
		var category pb.Category
		err = rows.Scan(&category.CategoryId, &category.Name, &category.ParentUuid, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, &category)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM categories where deleted_at is null`).Scan(&count)

	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func (r *catalogRepo) UpdateCategory(category pb.Category) (pb.Category, error) {
	result, err := r.db.Exec(`UPDATE categories SET name=$1, parent_uuid=$2, updated_at = $3 WHERE category_id=$4 and deleted_at is null`,
		category.Name, category.ParentUuid, time.Now().UTC(), category.CategoryId,
	)
	if err != nil {
		return pb.Category{}, err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Category{}, sql.ErrNoRows
	}

	var newCategory pb.Category

	newCategory, err = r.GetCategory(category.CategoryId)
	if err != nil {
		return pb.Category{}, err
	}

	return newCategory, nil
}

func (r *catalogRepo) DeleteCategory(id string) error {
	result, err := r.db.Exec(`UPDATE categories SET deleted_at = $1 WHERE category_id=$2`, time.Now().UTC(), id)
  if err != nil {
		  return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *catalogRepo) UpdateAuthor(update pb.Author) (pb.Author, error) {
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

func (r *catalogRepo) DeleteAuthor(id string) error {
	result, err := r.db.Exec("UPDATE authors SET deleted_at = $2 WHERE author_id = $1", id, time.Now().UTC())
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
