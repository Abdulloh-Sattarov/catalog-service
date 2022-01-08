package storage

import (
	"github.com/jmoiron/sqlx"

	"github.com/abdullohsattorov/catalog-service/storage/postgres"
	"github.com/abdullohsattorov/catalog-service/storage/repo"
)

// IStorage ...
type IStorage interface {
	Book() repo.BookStorageI
	Author() repo.AuthorStorageI
	Category() repo.CategoryStorageI
}

type storagePg struct {
	db           *sqlx.DB
	bookRepo     repo.BookStorageI
	authorRepo   repo.AuthorStorageI
	categoryRepo repo.CategoryStorageI
}

// NewStoragePg ...
func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:           db,
		bookRepo:     postgres.NewBookRepo(db),
		authorRepo:   postgres.NewAuthorRepo(db),
		categoryRepo: postgres.NewCategoryRepo(db),
	}
}

func (s storagePg) Book() repo.BookStorageI {
	return s.bookRepo
}

func (s storagePg) Author() repo.AuthorStorageI {
	return s.authorRepo
}

func (s storagePg) Category() repo.CategoryStorageI {
	return s.categoryRepo
}
