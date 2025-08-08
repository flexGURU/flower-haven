package repository

import (
	"context"
	"time"

	"github.com/flexGURU/flower-haven/backend/pkg"
)

type Order struct {
	ID              int64      `json:"id"`
	UserName        string     `json:"user_name"`
	UserPhoneNumber string     `json:"user_phone_number"`
	TotalAmount     float64    `json:"total_amount"`
	PaymentStatus   bool       `json:"payment_status"`
	Status          string     `json:"status"`
	ShippingAddress *string    `json:"shipping_address,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type UpdateOrder struct {
	ID              uint32  `json:"id"`
	UserName        *string `json:"user_name"`
	UserPhoneNumber *string `json:"user_phone_number"`
	PaymentStatus   *bool   `json:"payment_status"`
	Status          *string `json:"status"`
	ShippingAddress *string `json:"shipping_address,omitempty"`
}

type OrderFilter struct {
	Pagination    *pkg.Pagination
	Search        *string
	PaymentStatus *bool
	Status        *string
}

type OrderItem struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Amount    float64 `json:"amount"`
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrderByID(ctx context.Context, id int64) (*Order, error)
	UpdateOrder(ctx context.Context, order *UpdateOrder) (*Order, error)
	ListOrders(ctx context.Context, filter *OrderFilter) ([]*Order, *pkg.Pagination, error)
	DeleteOrder(ctx context.Context, id int64) error

	// Order Items
	CreateOrderItem(ctx context.Context, item *OrderItem) (*OrderItem, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]*OrderItem, error)
	GetOrderItemsByProductID(ctx context.Context, productID int64) ([]*OrderItem, error)
}
