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
	HasStems      bool       `json:"has_stems"`
	IsMessageCard bool       `json:"is_message_card"`
	IsFlowers     bool       `json:"is_flowers"`
	IsAddOn       bool       `json:"is_add_on"`
	ImageUrl      []string   `json:"image_url"`
	StockQuantity int64      `json:"stock_quantity"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	CategoryData  *Category  `json:"category_data,omitempty"`

	// extended fields
	Stems []ProductStem `json:"stems,omitempty"`
}

type UpdateProduct struct {
	ID            uint32    `json:"id"`
	Name          *string   `json:"name"`
	Description   *string   `json:"description"`
	HasStems      *bool     `json:"has_stems"`
	IsMessageCard *bool     `json:"is_message_card"`
	IsFlowers     *bool     `json:"is_flowers"`
	IsAddOn       *bool     `json:"is_add_on"`
	Price         *float64  `json:"price"`
	CategoryID    *uint32   `json:"category_id"`
	ImageURL      *[]string `json:"image_url"`
	StockQuantity *int64    `json:"stock_quantity"`

	// extended fields
	Stems []UpdateProductStem `json:"stems,omitempty"`
}

type ProductStem struct {
	ID        uint32  `json:"id"`
	ProductID uint32  `json:"product_id"`
	StemCount uint32  `json:"stem_count"`
	Price     float64 `json:"price"`
}

type UpdateProductStem struct {
	StemCount *uint32  `json:"stem_count"`
	Price     *float64 `json:"price"`
}

type ProductFilter struct {
	Pagination    *pkg.Pagination
	Search        *string
	PriceFrom     *float64
	PriceTo       *float64
	CategoryIDs   *[]int64
	IsMessageCard *bool
	IsFlowers     *bool
	IsAddOn       *bool
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *Product) (*Product, error)
	GetProductByID(ctx context.Context, id int64) (*Product, error)
	UpdateProduct(ctx context.Context, product *UpdateProduct) (*Product, error)
	ListProducts(ctx context.Context, filter *ProductFilter) ([]*Product, *pkg.Pagination, error)
	DeleteProduct(ctx context.Context, id int64) error

	GetDashboardData(ctx context.Context) (interface{}, error)
	ListAddOns(ctx context.Context) ([]*Product, error)
	ListMessageCards(ctx context.Context) ([]*Product, error)
}
