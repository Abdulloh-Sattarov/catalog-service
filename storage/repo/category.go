package repo

import pb "github.com/abdullohsattorov/catalog-service/genproto/catalog_service"

type CategoryStorageI interface {
	CreateCategory(category pb.Category) (pb.Category, error)
	GetCategory(id string) (pb.Category, error)
	ListCategory(page, limit int64) ([]*pb.Category, int64, error)
	UpdateCategory(category pb.Category) (pb.Category, error)
	DeleteCategory(id string) error
}
