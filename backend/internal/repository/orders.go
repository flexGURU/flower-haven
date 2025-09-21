package repository

import (
	"context"
	"time"

	"github.com/flexGURU/flower-haven/backend/pkg"
)

type Order struct {
	ID              uint32      `json:"id"`
	UserName        string      `json:"user_name"`
	UserPhoneNumber string      `json:"user_phone_number"`
	TotalAmount     float64     `json:"total_amount"`
	PaymentStatus   bool        `json:"payment_status"`
	Status          string      `json:"status"`
	DeliveryDate    time.Time   `json:"delivery_date"`
	TimeSlot        string      `json:"time_slot"`
	ByAdmin         bool        `json:"by_admin"`
	ShippingAddress *string     `json:"shipping_address,omitempty"`
	DeletedAt       *time.Time  `json:"deleted_at,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	OrderItemsData  []OrderItem `json:"order_item_data,omitempty"`
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
	ID                    uint32   `json:"id"`
	OrderID               uint32   `json:"order_id"`
	ProductID             uint32   `json:"product_id"`
	StemID                uint32   `json:"stem_id,omitempty"`
	PaymentMethod         string   `json:"payment_method"`
	Frequency             string   `json:"frequency"`
	Quantity              int32    `json:"quantity"`
	Amount                float64  `json:"amount"`
	OrderData             *Order   `json:"order_data,omitempty"`
	CurrentProductDetails *Product `json:"current_product_details,omitempty"`
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order, orderItems []OrderItem) (*Order, error)
	GetOrderByID(ctx context.Context, id int64) (*Order, error)
	UpdateOrder(ctx context.Context, order *UpdateOrder) (*Order, error)
	ListOrders(ctx context.Context, filter *OrderFilter) ([]*Order, *pkg.Pagination, error)
	DeleteOrder(ctx context.Context, id int64) error

	// Order Items
	GetOrderItemsByProductID(ctx context.Context, productID int64, filter *OrderFilter) ([]*OrderItem, *pkg.Pagination, error)
}
