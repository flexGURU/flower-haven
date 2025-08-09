package repository

import (
	"context"
	"time"

	"github.com/flexGURU/flower-haven/backend/pkg"
)

type Payment struct {
	ID                 uint32    `json:"id"`
	Description        *string   `json:"description,omitempty"`
	OrderID            *uint32   `json:"order_id,omitempty"`
	UserSubscriptionID *uint32   `json:"user_subscription_id,omitempty"`
	PaymentMethod      string    `json:"payment_method"`
	Amount             float64   `json:"amount"`
	PaidAt             time.Time `json:"paid_at"`
	CreatedAt          time.Time `json:"created_at"`
}

type UpdatePayment struct {
	ID            uint32     `json:"id"`
	Description   *string    `json:"description,omitempty"`
	PaymentMethod *string    `json:"payment_method"`
	Amount        *float64   `json:"amount"`
	PaidAt        *time.Time `json:"paid_at"`
}

type PaymentFilter struct {
	Pagination    *pkg.Pagination
	PaymentMethod *string
	StartDate     *time.Time
	EndDate       *time.Time
}

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *Payment) (*Payment, error)
	GetPaymentByID(ctx context.Context, id int64) (*Payment, error)
	UpdatePayment(ctx context.Context, payment *UpdatePayment) (*Payment, error)
	ListPayments(ctx context.Context, filter *PaymentFilter) ([]*Payment, *pkg.Pagination, error)
}
