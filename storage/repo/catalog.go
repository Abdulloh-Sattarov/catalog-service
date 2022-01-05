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
}
