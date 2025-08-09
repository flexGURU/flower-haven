package repository

import (
	"context"
	"time"

	"github.com/flexGURU/flower-haven/backend/pkg"
)

type Subscription struct {
	ID           uint32     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	ProductIds   []uint32   `json:"product_ids"`
	ProductsData []Product  `json:"products_data,omitempty"`
	AddOns       []uint32   `json:"add_ons"`
	AddOnsData   []Product  `json:"add_ons_data,omitempty"`
	Price        float64    `json:"price"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type UpdateSubscription struct {
	ID          uint32    `json:"id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	ProductIds  *[]uint32 `json:"product_ids"`
	AddOns      *[]uint32 `json:"add_ons"`
	Price       *float64  `json:"price"`
}

type SubscriptionFilter struct {
	Pagination *pkg.Pagination
	Search     *string
	PriceFrom  *float64
	PriceTo    *float64
}

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, subscription *Subscription) (*Subscription, error)
	GetSubscriptionByID(ctx context.Context, id int64) (*Subscription, error)
	UpdateSubscription(ctx context.Context, subscription *UpdateSubscription) (*Subscription, error)
	ListSubscriptions(ctx context.Context, filter *SubscriptionFilter) ([]*Subscription, *pkg.Pagination, error)
	DeleteSubscription(ctx context.Context, id int64) error
}
