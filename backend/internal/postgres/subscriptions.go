package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/flexGURU/flower-haven/backend/internal/postgres/generated"
	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ repository.SubscriptionRepository = (*SubscriptionRepository)(nil)

type SubscriptionRepository struct {
	queries *generated.Queries
}

func NewSubscriptionRepository(queries *generated.Queries) *SubscriptionRepository {
	return &SubscriptionRepository{queries: queries}
}

func (sr *SubscriptionRepository) CreateSubscription(ctx context.Context, subscription *repository.Subscription) (*repository.Subscription, error) {
	params := generated.CreateSubscriptionParams{
		Name:        subscription.Name,
		Description: subscription.Description,
		Price:       pkg.Float64ToPgTypeNumeric(subscription.Price),
		ProductIds:  make([]int32, 0, len(subscription.ProductIds)),
		AddOns:      make([]int32, 0, len(subscription.AddOns)),
	}

	for _, productId := range subscription.ProductIds {
		if exist, _ := sr.queries.ProductExists(ctx, int64(productId)); !exist {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "product with id %d not found", productId)
		}
		params.ProductIds = append(params.ProductIds, int32(productId))
	}

	for _, addOnId := range subscription.AddOns {
		if exist, _ := sr.queries.ProductExists(ctx, int64(addOnId)); !exist {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "add-on with id %d not found", addOnId)
		}
		params.AddOns = append(params.AddOns, int32(addOnId))
	}

	subscriptionId, err := sr.queries.CreateSubscription(ctx, params)
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
			return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating subscription: %s", err.Error())
	}

	return sr.GetSubscriptionByID(ctx, subscriptionId)
}

func (sr *SubscriptionRepository) GetSubscriptionByID(ctx context.Context, id int64) (*repository.Subscription, error) {
	generatedSubscription, err := sr.queries.GetSubscriptionByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "subscription with ID %d not found", id)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching subscription by id: %s", err.Error())
	}

	return generatedSubToRepoSub(generated.Subscription{
		ID:          generatedSubscription.ID,
		Name:        generatedSubscription.Name,
		Description: generatedSubscription.Description,
		ProductIds:  generatedSubscription.ProductIds,
		AddOns:      generatedSubscription.AddOns,
		Price:       generatedSubscription.Price,
		DeletedAt:   generatedSubscription.DeletedAt,
		CreatedAt:   generatedSubscription.CreatedAt,
	}, generatedSubscription.ProductsData, generatedSubscription.AddOnsData)
}

func (sr *SubscriptionRepository) UpdateSubscription(ctx context.Context, subscription *repository.UpdateSubscription) (*repository.Subscription, error) {
	params := generated.UpdateSubscriptionParams{
		ID:          int64(subscription.ID),
		Name:        pgtype.Text{Valid: false},
		Description: pgtype.Text{Valid: false},
		ProductIds:  nil,
		AddOns:      nil,
		Price:       pgtype.Numeric{Valid: false},
	}

	if subscription.Name != nil {
		params.Name = pgtype.Text{
			Valid:  true,
			String: *subscription.Name,
		}
	}

	if subscription.Description != nil {
		params.Description = pgtype.Text{
			Valid:  true,
			String: *subscription.Description,
		}
	}

	if subscription.ProductIds != nil {
		params.ProductIds = make([]int32, len(*subscription.ProductIds))
		for idx, productId := range *subscription.ProductIds {
			params.ProductIds[idx] = int32(productId)
		}
	}

	if subscription.AddOns != nil {
		params.AddOns = make([]int32, len(*subscription.AddOns))
		for idx, addOnId := range *subscription.AddOns {
			params.AddOns[idx] = int32(addOnId)
		}
	}

	if subscription.Price != nil {
		params.Price = pkg.Float64ToPgTypeNumeric(*subscription.Price)
	}

	subscriptionId, err := sr.queries.UpdateSubscription(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "subscription with ID %d not found", subscription.ID)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error updating subscription by id: %s", err.Error())
	}

	return sr.GetSubscriptionByID(ctx, subscriptionId)
}

func (sr *SubscriptionRepository) ListSubscriptions(ctx context.Context, filter *repository.SubscriptionFilter) ([]*repository.Subscription, *pkg.Pagination, error) {
	paramsListSubscriptions := generated.ListSubscriptionsParams{
		Limit:     int32(filter.Pagination.PageSize),
		Offset:    pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
		Search:    pgtype.Text{Valid: false},
		PriceFrom: pgtype.Float8{Valid: false},
		PriceTo:   pgtype.Float8{Valid: false},
	}

	paramsCountSubscriptions := generated.ListSubscriptionsCountParams{
		Search:    pgtype.Text{Valid: false},
		PriceFrom: pgtype.Float8{Valid: false},
		PriceTo:   pgtype.Float8{Valid: false},
	}

	if filter.Search != nil {
		search := strings.ToLower(*filter.Search)
		paramsListSubscriptions.Search = pgtype.Text{
			Valid:  true,
			String: "%" + search + "%",
		}
		paramsCountSubscriptions.Search = pgtype.Text{
			Valid:  true,
			String: "%" + search + "%",
		}
	}

	if filter.PriceFrom != nil && filter.PriceTo != nil {
		paramsListSubscriptions.PriceFrom = pgtype.Float8{
			Valid:   true,
			Float64: *filter.PriceFrom,
		}
		paramsListSubscriptions.PriceTo = pgtype.Float8{
			Valid:   true,
			Float64: *filter.PriceTo,
		}

		paramsCountSubscriptions.PriceFrom = pgtype.Float8{
			Valid:   true,
			Float64: *filter.PriceFrom,
		}
		paramsCountSubscriptions.PriceTo = pgtype.Float8{
			Valid:   true,
			Float64: *filter.PriceTo,
		}
	}

	generatedSubscriptions, err := sr.queries.ListSubscriptions(ctx, paramsListSubscriptions)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing subscriptions: %s", err.Error())
	}

	totalCount, err := sr.queries.ListSubscriptionsCount(ctx, paramsCountSubscriptions)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error counting users: %s", err.Error())
	}

	subscriptionList := make([]*repository.Subscription, len(generatedSubscriptions))
	for idx, generatedSubscription := range generatedSubscriptions {
		subscriptionList[idx], err = generatedSubToRepoSub(generated.Subscription{
			ID:          generatedSubscription.ID,
			Name:        generatedSubscription.Name,
			Description: generatedSubscription.Description,
			ProductIds:  generatedSubscription.ProductIds,
			AddOns:      generatedSubscription.AddOns,
			Price:       generatedSubscription.Price,
			DeletedAt:   generatedSubscription.DeletedAt,
			CreatedAt:   generatedSubscription.CreatedAt,
		}, generatedSubscription.ProductsData, generatedSubscription.AddOnsData)
		if err != nil {
			return nil, nil, err
		}
	}

	return subscriptionList, pkg.CalculatePagination(uint32(totalCount), filter.Pagination.PageSize, filter.Pagination.Page), nil
}

func (sr *SubscriptionRepository) DeleteSubscription(ctx context.Context, id int64) error {
	if err := sr.queries.DeleteSubscription(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.Errorf(pkg.NOT_FOUND_ERROR, "subscription with ID %d not found", id)
		}
		return pkg.Errorf(pkg.INTERNAL_ERROR, "error deleting subscription by id: %s", err.Error())
	}

	return nil
}

func generatedSubToRepoSub(genSub generated.Subscription, productsData, addOnsData []byte) (*repository.Subscription, error) {
	subscription := &repository.Subscription{
		ID:           uint32(genSub.ID),
		Name:         genSub.Name,
		Description:  genSub.Description,
		Price:        pkg.PgTypeNumericToFloat64(genSub.Price),
		CreatedAt:    genSub.CreatedAt,
		ProductIds:   make([]uint32, len(genSub.ProductIds)),
		AddOns:       make([]uint32, len(genSub.AddOns)),
		ProductsData: nil,
		AddOnsData:   nil,
	}

	if genSub.DeletedAt.Valid {
		subscription.DeletedAt = &genSub.DeletedAt.Time
	}

	for idx, productId := range genSub.ProductIds {
		subscription.ProductIds[idx] = uint32(productId)
	}

	for idx, addOn := range genSub.AddOns {
		subscription.AddOns[idx] = uint32(addOn)
	}

	if productsData != nil {
		var productData []repository.Product
		if err := json.Unmarshal(productsData, &productData); err != nil {
			return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed unmarshalling products data to subscription: %s", err.Error())
		}
		subscription.ProductsData = productData
	}

	if addOnsData != nil {
		var addOnData []repository.Product
		if err := json.Unmarshal(addOnsData, &addOnData); err != nil {
			return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed unmarshalling addOns data to subscription: %s", err.Error())
		}
		subscription.AddOnsData = addOnData
	}

	return subscription, nil
}
