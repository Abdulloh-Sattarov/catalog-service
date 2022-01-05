package repo

import (
	pb "github.com/abdullohsattorov/catalog-service/genproto/catalog_service"
)

// CatalogStorageI ...
type CatalogStorageI interface {
	CreateBook(book pb.Book) (pb.Book, error)
	GetBook(id string) (pb.Book, error)
	ListBook(page, limit int64) ([]*pb.Book, int64, error)
	UpdateBook(update pb.Book) (pb.Book, error)
	DeleteBook(id string) error

	CreateCategory(category pb.Category) (pb.Category, error)
	GetCategory(id string) (pb.Category, error)
	ListCategory(page, limit int64) ([]*pb.Category, int64, error)
	UpdateCategory(category pb.Category) (pb.Category, error)
	DeleteCategory(id string) error


	CreateAuthor(author pb.Author) (pb.Author, error)
	GetAuthor(id string) (pb.Author, error)
	ListAuthor(page, limit int64) ([]*pb.Author, int64, error)
	UpdateAuthor(update pb.Author) (pb.Author, error)
	DeleteAuthor(id string) error

}
