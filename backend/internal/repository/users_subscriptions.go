package repository

import (
	"context"
	"time"

	"github.com/flexGURU/flower-haven/backend/pkg"
)

type UserSubscription struct {
	ID               uint32       `json:"id"`
	UserID           uint32       `json:"user_id"`
	UserData         User         `json:"user_data,omitempty"`
	SubscriptionID   uint32       `json:"subscription_id"`
	SubscriptionData Subscription `json:"subscription_data,omitempty"`
	DayOfWeek        int16        `json:"day_of_week"`
	Status           bool         `json:"status"`
	StartDate        time.Time    `json:"start_date"`
	EndDate          time.Time    `json:"end_date"`
	DeletedAt        *time.Time   `json:"deleted_at,omitempty"`
	CreatedAt        time.Time    `json:"created_at"`
	PaymentData      Payment      `json:"payment_data,omitempty"`
}

type UpdateUserSubscription struct {
	ID        uint32     `json:"id"`
	DayOfWeek *int16     `json:"day_of_week"`
	Status    *bool      `json:"status"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

type UserSubscriptionFilter struct {
	Pagination *pkg.Pagination
	Status     *bool
}

type UserSubscriptionRepository interface {
	CreateUserSubscription(ctx context.Context, subscription *UserSubscription) (*UserSubscription, error)
	GetUserSubscriptionByID(ctx context.Context, id int64) (*UserSubscription, error)
	GetUsersSubscriptionsByUserID(ctx context.Context, userId int64) ([]*UserSubscription, error)
	UpdateUserSubscription(ctx context.Context, subscription *UpdateUserSubscription) (*UserSubscription, error)
	ListUserSubscriptions(ctx context.Context, filter *UserSubscriptionFilter) ([]*UserSubscription, *pkg.Pagination, error)
	DeleteUserSubscription(ctx context.Context, id int64) error
}
