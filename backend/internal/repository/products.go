package repository

import (
	"context"
	"time"

	"github.com/flexGURU/flower-haven/backend/pkg"
)

type Product struct {
	ID            uint32     `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Price         float64    `json:"price"`
	CategoryID    uint32     `json:"category_id"`
	CategoryData  *Category  `json:"category_data,omitempty"`
	ImageUrl      []string   `json:"image_url"`
	StockQuantity int64      `json:"stock_quantity"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type UpdateProduct struct {
	ID            uint32    `json:"id"`
	Name          *string   `json:"name"`
	Description   *string   `json:"description"`
	Price         *float64  `json:"price"`
	CategoryID    *uint32   `json:"category_id"`
	ImageURL      *[]string `json:"image_url"`
	StockQuantity *int64    `json:"stock_quantity"`
}

type ProductFilter struct {
	Pagination  *pkg.Pagination
	Search      *string
	PriceFrom   *float64
	PriceTo     *float64
	CategoryIDs *[]int64
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *Product) (*Product, error)
	GetProductByID(ctx context.Context, id int64) (*Product, error)
	UpdateProduct(ctx context.Context, product *UpdateProduct) (*Product, error)
	ListProducts(ctx context.Context, filter *ProductFilter) ([]*Product, *pkg.Pagination, error)
	DeleteProduct(ctx context.Context, id int64) error
}
