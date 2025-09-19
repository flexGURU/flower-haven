package repository

import (
	"context"
	"time"

	"github.com/flexGURU/flower-haven/backend/pkg"
)

type PaystackPayment struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Amount    string    `json:"amount"`
	Reference string    `json:"reference"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaystackEvent struct {
	ID        int64     `json:"id"`
	Event     string    `json:"event"`
	Data      any       `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

type PaystackRepository interface {
	CreatePayment(ctx context.Context, email string, amount int64, reference string) error
	GetPaymentByReference(ctx context.Context, reference string) (PaystackPayment, error)
	UpdatePaymentStatus(ctx context.Context, reference string, status string) error
	ListPaystackPayments(ctx context.Context, status string, pagination *pkg.Pagination) ([]PaystackPayment, *pkg.Pagination, error)

	LogPaystackEvent(ctx context.Context, event string, payload []byte) error
	ListPaystackEvents(ctx context.Context, event string, pagination *pkg.Pagination) ([]PaystackEvent, *pkg.Pagination, error)
}
