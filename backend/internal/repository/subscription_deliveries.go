package repository

import (
	"context"
	"time"

	"github.com/flexGURU/flower-haven/backend/pkg"
)

type SubscriptionDelivery struct {
	ID                 uint32     `json:"id"`
	Description        *string    `json:"description,omitempty"`
	UserSubscriptionID uint32     `json:"user_subscription_id"`
	DeliveredOn        time.Time  `json:"delivered_on"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

type UpdateSubscriptionDelivery struct {
	ID          uint32     `json:"id"`
	Description *string    `json:"description,omitempty"`
	DeliveredOn *time.Time `json:"delivered_on,omitempty"`
}

type SubscriptionDeliveryFilter struct {
	Pagination *pkg.Pagination
}

type SubscriptionDeliveryRepository interface {
	CreateSubscriptionDelivery(ctx context.Context, delivery *SubscriptionDelivery) (*SubscriptionDelivery, error)
	GetSubscriptonDeliveryByUserSubscriptionID(ctx context.Context, userSubscriptionID int64) ([]*SubscriptionDelivery, error)
	UpdateSubscriptionDelivery(ctx context.Context, delivery *UpdateSubscriptionDelivery) (*SubscriptionDelivery, error)
	ListSubscriptionDeliveries(ctx context.Context, filter *SubscriptionDeliveryFilter) ([]*SubscriptionDelivery, *pkg.Pagination, error)
	DeleteSubscriptionDelivery(ctx context.Context, id int64) error
}
