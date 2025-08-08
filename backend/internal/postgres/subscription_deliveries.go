package postgres

import (
	"context"

	"github.com/flexGURU/flower-haven/backend/internal/postgres/generated"
	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ repository.SubscriptionDeliveryRepository = (*SubscriptionDeliveryRepository)(nil)

type SubscriptionDeliveryRepository struct {
	queries *generated.Queries
}

func NewSubscriptionDeliveryRepository(queries *generated.Queries) *SubscriptionDeliveryRepository {
	return &SubscriptionDeliveryRepository{queries: queries}
}

func (sd *SubscriptionDeliveryRepository) CreateSubscriptionDelivery(ctx context.Context, delivery *repository.SubscriptionDelivery) (*repository.SubscriptionDelivery, error) {
	params := generated.CreateSubscriptionDeliveryParams{
		DeliveredOn: delivery.DeliveredOn,
		Description: pgtype.Text{Valid: false},
	}

	if exists, _ := sd.queries.UserSubscriptionExists(ctx, int64(delivery.UserSubscriptionID)); !exists {
		return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user subscription with id %d not found", delivery.UserSubscriptionID)
	}
	params.UserSubscriptionID = int64(delivery.UserSubscriptionID)

	if delivery.Description != nil {
		params.Description = pgtype.Text{
			Valid:  true,
			String: *delivery.Description,
		}
	}

	generatedDelivery, err := sd.queries.CreateSubscriptionDelivery(ctx, params)
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
			return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating subscription delivery: %s", err.Error())
	}

	delivery.ID = uint32(generatedDelivery.ID)
	delivery.CreatedAt = generatedDelivery.CreatedAt

	return delivery, nil
}

func (sd *SubscriptionDeliveryRepository) GetSubscriptonDeliveryByUserSubscriptionID(ctx context.Context, userSubscriptionID int64) ([]*repository.SubscriptionDelivery, error) {
	generatedDeliveries, err := sd.queries.GetSubscriptionDeliveryByUserSubscriptionID(ctx, userSubscriptionID)
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.NOT_FOUND_ERROR {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no deliveries found for user subscription id %d", userSubscriptionID)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching subscription deliveries: %s", err.Error())
	}

	deliveries := make([]*repository.SubscriptionDelivery, len(generatedDeliveries))
	for i, d := range generatedDeliveries {
		deliveries[i] = &repository.SubscriptionDelivery{
			ID:                 uint32(d.ID),
			Description:        &d.Description.String,
			UserSubscriptionID: uint32(d.UserSubscriptionID),
			DeliveredOn:        d.DeliveredOn,
			DeletedAt:          &d.DeletedAt.Time,
			CreatedAt:          d.CreatedAt,
		}
	}

	return deliveries, nil
}

func (sd *SubscriptionDeliveryRepository) UpdateSubscriptionDelivery(ctx context.Context, delivery *repository.UpdateSubscriptionDelivery) (*repository.SubscriptionDelivery, error) {
	params := generated.UpdateSubscriptionDeliveryParams{
		ID:          int64(delivery.ID),
		Description: pgtype.Text{Valid: false},
		DeliveredOn: pgtype.Timestamptz{Valid: false},
	}

	if delivery.DeliveredOn != nil {
		params.DeliveredOn = pgtype.Timestamptz{
			Valid: true,
			Time:  *delivery.DeliveredOn,
		}
	}

	if delivery.Description != nil {
		params.Description = pgtype.Text{
			Valid:  true,
			String: *delivery.Description,
		}
	}

	generatedDelivery, err := sd.queries.UpdateSubscriptionDelivery(ctx, params)
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.NOT_FOUND_ERROR {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "subscription delivery with id %d not found", delivery.ID)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error updating subscription delivery: %s", err.Error())
	}

	return &repository.SubscriptionDelivery{
		ID:                 uint32(generatedDelivery.ID),
		Description:        &generatedDelivery.Description.String,
		UserSubscriptionID: uint32(generatedDelivery.UserSubscriptionID),
		DeliveredOn:        generatedDelivery.DeliveredOn,
		DeletedAt:          &generatedDelivery.DeletedAt.Time,
		CreatedAt:          generatedDelivery.CreatedAt,
	}, nil
}

func (sd *SubscriptionDeliveryRepository) ListSubscriptionDeliveries(ctx context.Context, filter *repository.SubscriptionDeliveryFilter) ([]*repository.SubscriptionDelivery, *pkg.Pagination, error) {
	userDeliveries, err := sd.queries.ListSubscriptionDelivery(ctx, generated.ListSubscriptionDeliveryParams{
		Limit:  int32(filter.Pagination.PageSize),
		Offset: pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
	})
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching subscription deliveries: %s", err.Error())
	}

	totalCount, err := sd.queries.ListCountSubscriptionDelivery(ctx)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error counting subscription deliveries: %s", err.Error())
	}

	userSubacriptionList := make([]*repository.SubscriptionDelivery, len(userDeliveries))
	for i, userSub := range userDeliveries {
		userSubacriptionList[i] = &repository.SubscriptionDelivery{
			ID:                 uint32(userSub.ID),
			Description:        &userSub.Description.String,
			UserSubscriptionID: uint32(userSub.UserSubscriptionID),
			DeliveredOn:        userSub.DeliveredOn,
			DeletedAt:          &userSub.DeletedAt.Time,
			CreatedAt:          userSub.CreatedAt,
		}
	}

	return userSubacriptionList, pkg.CalculatePagination(uint32(totalCount), filter.Pagination.PageSize, filter.Pagination.Page), nil
}

func (sd *SubscriptionDeliveryRepository) DeleteSubscriptionDelivery(ctx context.Context, id int64) error {
	if err := sd.queries.DeleteSubscriptionDelivery(ctx, id); err != nil {
		if pkg.PgxErrorCode(err) == pkg.NOT_FOUND_ERROR {
			return pkg.Errorf(pkg.NOT_FOUND_ERROR, "subscription delivery with id %d not found", id)
		}
		return pkg.Errorf(pkg.INTERNAL_ERROR, "error deleting subscription delivery: %s", err.Error())
	}
	return nil
}
