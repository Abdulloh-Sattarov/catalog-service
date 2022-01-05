package repo

import (
	pb "github.com/abdullohsattorov/catalog-service/genproto"
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
}
