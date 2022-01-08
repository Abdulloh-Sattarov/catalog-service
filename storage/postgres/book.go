package postgres

import (
	"database/sql"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"

	pb "github.com/abdullohsattorov/catalog-service/genproto/catalog_service"
	"github.com/abdullohsattorov/catalog-service/pkg/utils"
)

type bookRepo struct {
	db *sqlx.DB
}

// NewBookRepo ...
func NewBookRepo(db *sqlx.DB) *bookRepo {
	return &bookRepo{db: db}
}

func (r *bookRepo) CreateBook(book pb.Book) (pb.Book, error) {
	var id string
	err := r.db.QueryRow(`
        INSERT INTO books(book_id, name, author_id, price, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6) returning book_id`, book.BookId, book.Name, book.AuthorId, book.Price, time.Now().UTC(), time.Now().UTC()).Scan(&id)
	if err != nil {
		return pb.Book{}, err
	}

	_, err = r.db.Exec(`
        INSERT INTO book_categories(book_id, category_id)
        VALUES ($1, $2)`, book.BookId, book.CategoryId)
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

func (r *bookRepo) GetBook(id string) (pb.Book, error) {
	var book pb.Book
	err := r.db.QueryRow(`
        SELECT book_id, name, author_id, price, created_at, updated_at FROM books
        WHERE book_id=$1 and deleted_at is null`, id).Scan(&book.BookId, &book.Name, &book.AuthorId, &book.Price, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		return pb.Book{}, err
	}

	err = r.db.QueryRow(`
        select c.category_id, c.name
			from book_categories
		join books b on book_categories.book_id = b.book_id
		join categories c on book_categories.category_id = c.category_id
		where b.book_id = $1`, id).Scan(&book.CategoryId, &book.CategoryName)
	if err != nil {
		return pb.Book{}, err
	}

	return book, nil
}

func (r *bookRepo) ListBook(page, limit int64, filters map[string]string) ([]*pb.Book, int64, error) {
	offset := (page - 1) * limit

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("b.book_id", "b.name", "b.author_id", "b.price", "c.category_id", "c.name", "b.created_at", "b.updated_at")
	sb.From("book_categories")
	sb.Join("books b", "book_categories.book_id=b.book_id")
	sb.Join("categories c", "book_categories.category_id=c.category_id")

	if value, ok := filters["authors"]; ok {
		args := utils.StringSliceToInterfaceSlice(utils.ParseFilter(value))
		sb.Join("authors a", "b.author_id=a.author_id")
		sb.Where(sb.In("a.author_id", args...))
	}

	if value, ok := filters["category"]; ok {
		sb.Where(sb.Equal("c.category_id", value))
	}

	sb.Limit(int(limit))
	sb.Offset(int(offset))

	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var (
		books []*pb.Book
		count int64
	)

	for rows.Next() {
		var book pb.Book
		err = rows.Scan(&book.BookId, &book.Name, &book.AuthorId, &book.Price, &book.CategoryId, &book.CategoryName, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		books = append(books, &book)
	}

	sbc := sqlbuilder.NewSelectBuilder()
	sbc.Select("count(*)")
	sbc.From("book_categories")
	sbc.Join("books b", "book_categories.book_id=b.book_id")
	sbc.Join("categories c", "book_categories.category_id=c.category_id")

	if value, ok := filters["authors"]; ok {
		elements := utils.StringSliceToInterfaceSlice(utils.ParseFilter(value))
		sbc.Join("authors a", "b.author_id=a.author_id")
		sbc.Where(sbc.In("a.author_id", elements...))
	}

	if value, ok := filters["category"]; ok {
		sbc.Where(sbc.Equal("c.category_id", value))
	}

	query, args = sbc.BuildWithFlavor(sqlbuilder.PostgreSQL)

	err = r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return books, count, nil
}

func (r *bookRepo) UpdateBook(book pb.Book) (pb.Book, error) {
	result, err := r.db.Exec(`UPDATE books SET name=$1, author_id=$2, price=$3, updated_at = $4 WHERE book_id=$5 and deleted_at is null`,
		book.Name, book.AuthorId, book.Price, time.Now().UTC(), book.BookId,
	)
	if err != nil {
		return pb.Book{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Book{}, sql.ErrNoRows
	}

	resultNext, err := r.db.Exec(`UPDATE book_categories SET category_id = $1 WHERE book_id=$2`,
		book.CategoryId, book.BookId,
	)
	if err != nil {
		return pb.Book{}, err
	}

	if i, _ := resultNext.RowsAffected(); i == 0 {
		return pb.Book{}, sql.ErrNoRows
	}

	var NewBook pb.Book

	NewBook, err = r.GetBook(book.BookId)

	if err != nil {
		return pb.Book{}, err
	}

	return NewBook, nil
}

func (r *bookRepo) DeleteBook(id string) error {
	result, err := r.db.Exec(`UPDATE books SET deleted_at = $1 WHERE book_id=$2`, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
