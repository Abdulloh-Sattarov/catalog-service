package repo

import pb "github.com/abdullohsattorov/catalog-service/genproto/catalog_service"

type AuthorStorageI interface {
	CreateAuthor(author pb.Author) (pb.Author, error)
	GetAuthor(id string) (pb.Author, error)
	ListAuthor(page, limit int64) ([]*pb.Author, int64, error)
	UpdateAuthor(update pb.Author) (pb.Author, error)
	DeleteAuthor(id string) error
}
