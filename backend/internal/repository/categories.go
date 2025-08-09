package repository

import (
	"context"
	"time"

	"github.com/flexGURU/flower-haven/backend/pkg"
)

type Category struct {
	ID          uint32     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ImageUrl    []string   `json:"image_url"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type UpdateCategory struct {
	ID          uint32    `json:"id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	ImageUrl    *[]string `json:"image_url"`
}

type CategoryFilter struct {
	Pagination *pkg.Pagination
	Search     *string
}

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	GetCategoryByID(ctx context.Context, id int64) (*Category, error)
	UpdateCategory(ctx context.Context, category *UpdateCategory) (*Category, error)
	ListCategories(ctx context.Context, filter *CategoryFilter) ([]*Category, *pkg.Pagination, error)
	DeleteCategory(ctx context.Context, id int64) error
}
