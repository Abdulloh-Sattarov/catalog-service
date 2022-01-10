package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	pb "github.com/abdullohsattorov/catalog-service/genproto/catalog_service"
)

type categoryRepo struct {
	db *sqlx.DB
}

// NewCategoryRepo ...
func NewCategoryRepo(db *sqlx.DB) *categoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) CreateCategory(category pb.Category) (pb.Category, error) {
	var (
		parentID sql.NullString
		id       string
	)
	parentID = stringToNullString(category.ParentUuid)

	err := r.db.QueryRow(`
	INSERT INTO categories(category_id, name, parent_uuid, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5) returning category_id`, category.CategoryId, category.Name, parentID, time.Now().UTC(), time.Now().UTC()).Scan(&id)
	if err != nil {
		return pb.Category{}, err
	}

	newCategory, err := r.GetCategory(id)
	if err != nil {
		return pb.Category{}, err
	}

	return newCategory, nil
}

func (r *categoryRepo) GetCategory(id string) (pb.Category, error) {
	var category pb.Category
	var (
		parentUUID     sql.NullString
		parentCategory sql.NullString
	)

	err := r.db.QueryRow(`
		SELECT cat.category_id, cat.name AS category_name, cat.parent_uuid, cat2.name AS parent_category, cat.created_at, cat.updated_at
		FROM categories AS cat 
		LEFT JOIN categories AS cat2 ON cat.parent_uuid = cat2.category_id
		WHERE cat.category_id=$1 AND cat.deleted_at is null;
		`, id).Scan(&category.CategoryId, &category.Name, &parentUUID, &parentCategory, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return pb.Category{}, err
	}

	if !parentUUID.Valid {
		parentUUID.String = ""
	}
	category.ParentUuid = parentUUID.String

	if !parentCategory.Valid {
		parentCategory.String = ""
	}
	category.ParentCategory = parentCategory.String

	return category, nil
}

func (r *categoryRepo) ListCategory(page, limit int64) ([]*pb.Category, int64, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Queryx(`
		SELECT cat.category_id, cat.name AS category_name, cat.parent_uuid, cat2.name AS parent_category, cat.created_at, cat.updated_at
		FROM categories AS cat 
		LEFT JOIN categories AS cat2 ON cat.parent_uuid = cat2.category_id
		WHERE cat.deleted_at is null 
		ORDER BY category_id LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var (
		categories []*pb.Category
		count      int64
	)
	var (
		parentUUID     sql.NullString
		parentCategory sql.NullString
	)
	for rows.Next() {
		var category pb.Category

		err = rows.Scan(&category.CategoryId, &category.Name, &parentUUID, &parentCategory, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}

		if !parentUUID.Valid {
			parentUUID.String = ""
		}
		category.ParentUuid = parentUUID.String

		if !parentCategory.Valid {
			parentCategory.String = ""
		}
		category.ParentCategory = parentCategory.String

		categories = append(categories, &category)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM categories where deleted_at is null`).Scan(&count)

	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func (r *categoryRepo) UpdateCategory(category pb.Category) (pb.Category, error) {
	var (
		parentID    sql.NullString
		newCategory pb.Category
	)

	parentID = stringToNullString(category.ParentUuid)

	result, err := r.db.Exec(`UPDATE categories SET name=$1, parent_uuid=$2, updated_at = $3 WHERE category_id=$4 and deleted_at is null`,
		category.Name, parentID, time.Now().UTC(), category.CategoryId,
	)
	if err != nil {
		return pb.Category{}, err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Category{}, sql.ErrNoRows
	}

	newCategory, err = r.GetCategory(category.CategoryId)
	if err != nil {
		return pb.Category{}, err
	}

	return newCategory, nil
}

func (r *categoryRepo) DeleteCategory(id string) error {
	var count int
	err := r.db.QueryRow(`select count(*) from categories where parent_uuid=$1 `, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		err1 := errors.New("this category has subcategories")
		return err1
	}

	result, err := r.db.Exec(`UPDATE categories SET deleted_at = $1 WHERE category_id=$2`, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func stringToNullString(s string) (ns sql.NullString) {
	if s != "" {
		ns.Valid = true
		ns.String = s
		return ns
	}

	return ns
}
