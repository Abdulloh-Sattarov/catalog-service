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

	for _, j := range book.CategoryId {
		_, err = r.db.Exec(`
        INSERT INTO book_categories(book_id, category_id)
        VALUES ($1, $2)`, book.BookId, j)
		if err != nil {
			return pb.Book{}, err
		}
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

	rows, err := r.db.Queryx(`
		select c.category_id, c.name, c.parent_uuid, cat.name, c.created_at, c.updated_at
			from book_categories
		join books b on book_categories.book_id = b.book_id
		join categories c on book_categories.category_id = c.category_id
		left join categories as cat ON c.parent_uuid = cat.category_id
		where b.book_id = $1`, id)
	if err != nil {
		return pb.Book{}, err
	}

	var (
		categories     []*pb.Category
		parentUUID     sql.NullString
		parentCategory sql.NullString
	)

	for rows.Next() {
		var category pb.Category
		err = rows.Scan(&category.CategoryId, &category.Name, &parentUUID, &parentCategory, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return pb.Book{}, err
		}

		if !parentUUID.Valid {
			parentUUID.String = ""
		}
		category.ParentUuid = parentUUID.String

		if !parentCategory.Valid {
			parentCategory.String = ""
		}
		category.ParentCategory = parentCategory.String

		book.CategoryId = append(book.CategoryId, category.CategoryId)

		categories = append(categories, &category)
	}

	book.Categories = categories

	return book, nil
}

func (r *bookRepo) ListBook(page, limit int64, filters map[string]string) ([]*pb.Book, int64, error) {
	offset := (page - 1) * limit

	sb := sqlbuilder.NewSelectBuilder()

	sb.Select("book_categories.book_id")
	sb.From("book_categories")
	sb.JoinWithOption("LEFT", "books b", "book_categories.book_id=b.book_id")
	sb.GroupBy("book_categories.book_id")
	if value, ok := filters["authors"]; ok {
		sb.JoinWithOption("LEFT", "authors a", "b.author_id=a.author_id")
		args := utils.StringSliceToInterfaceSlice(utils.ParseFilter(value))
		sb.Where(sb.In("a.author_id", args...))
	}

	if value, ok := filters["category"]; ok {
		sb.JoinWithOption("LEFT", "categories c", "c.category_id=book_categories.category_id")
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
		var byId pb.ByIdReq
		err = rows.Scan(&byId.Id)
		if err != nil {
			return nil, 0, err
		}
		book, _ := r.GetBook(byId.Id)
		books = append(books, &book)
	}

	sbc := sqlbuilder.NewSelectBuilder()

	sbc.Select("count(*)")
	sbc.From("book_categories")
	sbc.JoinWithOption("LEFT", "books b", "book_categories.book_id=b.book_id")
	sbc.GroupBy("book_categories.book_id")
	if value, ok := filters["authors"]; ok {
		args1 := utils.StringSliceToInterfaceSlice(utils.ParseFilter(value))
		sbc.JoinWithOption("LEFT", "authors a", "b.author_id=a.author_id")
		sbc.Where(sbc.In("a.author_id", args1...))
	}

	if value, ok := filters["category"]; ok {
		sbc.JoinWithOption("LEFT", "categories c", "c.category_id=book_categories.category_id")
		sbc.Where(sbc.Equal("c.category_id", value))
	}

	query, args = sbc.BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err = r.db.Queryx(query, args...)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		count = count + 1
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
