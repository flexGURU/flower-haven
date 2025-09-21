package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/flexGURU/flower-haven/backend/internal/postgres/generated"
	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ repository.UserSubscriptionRepository = (*UserSubscriptionRepository)(nil)

type UserSubscriptionRepository struct {
	queries *generated.Queries
}

func NewUserSubscriptionRepository(queries *generated.Queries) *UserSubscriptionRepository {
	return &UserSubscriptionRepository{queries: queries}
}

func (usr *UserSubscriptionRepository) CreateUserSubscription(ctx context.Context, subscription *repository.UserSubscription) (*repository.UserSubscription, error) {
	params := generated.CreateUserSubscriptionParams{
		UserID:         pgtype.Int8{Valid: true, Int64: int64(subscription.UserID)},
		SubscriptionID: int64(subscription.SubscriptionID),
		DayOfWeek:      subscription.DayOfWeek,
		StartDate:      subscription.StartDate,
		EndDate:        subscription.EndDate,
	}

	if exists, _ := usr.queries.UserExists(ctx, int64(subscription.UserID)); !exists {
		return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user with id %d not found", subscription.UserID)
	}

	if exists, _ := usr.queries.SubscriptionExists(ctx, int64(subscription.SubscriptionID)); !exists {
		return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "subscription with id %d not found", subscription.SubscriptionID)
	}

	userSubscriptionId, err := usr.queries.CreateUserSubscription(ctx, params)
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
			return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating user subscription: %s", err.Error())
	}

	// send message about the subscriptions

	return usr.GetUserSubscriptionByID(ctx, userSubscriptionId)
}

func (usr *UserSubscriptionRepository) GetUserSubscriptionByID(ctx context.Context, id int64) (*repository.UserSubscription, error) {
	userSubscription, err := usr.queries.GetUserSubscriptionByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user_subscription with ID %d not found", id)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching user_subscription by id: %s", err.Error())
	}

	return generatedUserSubToRepoUserSub(generated.UserSubscription{
		ID:             userSubscription.ID,
		UserID:         userSubscription.UserID,
		SubscriptionID: userSubscription.SubscriptionID,
		DayOfWeek:      userSubscription.DayOfWeek,
		Status:         userSubscription.Status,
		StartDate:      userSubscription.StartDate,
		EndDate:        userSubscription.EndDate,
		DeletedAt:      userSubscription.DeletedAt,
		CreatedAt:      userSubscription.CreatedAt,
	}, userSubscription.UserData, userSubscription.SubscriptionData, userSubscription.PaymentData)
}

func (usr *UserSubscriptionRepository) GetUsersSubscriptionsByUserID(ctx context.Context, userId int64, filter *repository.UserSubscriptionFilter) ([]*repository.UserSubscription, *pkg.Pagination, error) {
	userSubscriptions, err := usr.queries.GetUserSubscriptionsByUserID(ctx, generated.GetUserSubscriptionsByUserIDParams{
		UserID: pgtype.Int8{Valid: true, Int64: userId},
		Limit:  int32(filter.Pagination.PageSize),
		Offset: pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
	})
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing user's subscriptions with id %d: %s", userId, err.Error())
	}

	totalCount, err := usr.queries.GetCountUserSubscriptionsByUserID(ctx, pgtype.Int8{Valid: true, Int64: userId})
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error counting user's subscriptions with id %d: %s", userId, err.Error())
	}

	userSubacriptionList := make([]*repository.UserSubscription, len(userSubscriptions))
	for idx, userSub := range userSubscriptions {
		userSubacriptionList[idx], err = generatedUserSubToRepoUserSub(generated.UserSubscription{
			ID:             userSub.ID,
			UserID:         userSub.UserID,
			SubscriptionID: userSub.SubscriptionID,
			DayOfWeek:      userSub.DayOfWeek,
			Status:         userSub.Status,
			StartDate:      userSub.StartDate,
			EndDate:        userSub.EndDate,
			DeletedAt:      userSub.DeletedAt,
			CreatedAt:      userSub.CreatedAt,
		}, nil, userSub.SubscriptionData, nil)
		if err != nil {
			return nil, nil, err
		}
	}

	return userSubacriptionList, pkg.CalculatePagination(uint32(totalCount), filter.Pagination.PageSize, filter.Pagination.Page), nil
}

func (usr *UserSubscriptionRepository) UpdateUserSubscription(ctx context.Context, subscription *repository.UpdateUserSubscription) (*repository.UserSubscription, error) {
	params := generated.UpdateUserSubscriptionParams{
		ID:        int64(subscription.ID),
		StartDate: pgtype.Timestamptz{Valid: false},
		EndDate:   pgtype.Timestamptz{Valid: false},
		DayOfWeek: pgtype.Int2{Valid: false},
		Status:    pgtype.Bool{Valid: false},
	}

	if subscription.StartDate != nil {
		params.StartDate = pgtype.Timestamptz{
			Valid: true,
			Time:  *subscription.StartDate,
		}
	}
	if subscription.EndDate != nil {
		params.EndDate = pgtype.Timestamptz{
			Valid: true,
			Time:  *subscription.EndDate,
		}
	}
	if subscription.DayOfWeek != nil {
		params.DayOfWeek = pgtype.Int2{
			Valid: true,
			Int16: *subscription.DayOfWeek,
		}
	}
	if subscription.Status != nil {
		params.Status = pgtype.Bool{
			Valid: true,
			Bool:  *subscription.Status,
		}
	}

	userSubscriptionId, err := usr.queries.UpdateUserSubscription(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user_subscription with ID %d not found", subscription.ID)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error updating subscription by id: %s", err.Error())
	}

	return usr.GetUserSubscriptionByID(ctx, userSubscriptionId)
}

func (usr *UserSubscriptionRepository) ListUserSubscriptions(ctx context.Context, filter *repository.UserSubscriptionFilter) ([]*repository.UserSubscription, *pkg.Pagination, error) {
	paramsListUserSub := generated.ListUserSubscriptionsParams{
		Limit:  int32(filter.Pagination.PageSize),
		Offset: pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
		Status: pgtype.Bool{Valid: false},
	}

	paramsCountUserSub := pgtype.Bool{Valid: false}

	if filter.Status != nil {
		paramsListUserSub.Status = pgtype.Bool{
			Valid: true,
			Bool:  *filter.Status,
		}

		paramsCountUserSub.Valid = true
		paramsCountUserSub.Bool = *filter.Status
	}

	generatedUserSub, err := usr.queries.ListUserSubscriptions(ctx, paramsListUserSub)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing user_subscriptions %s", err.Error())
	}

	totalCount, err := usr.queries.ListCountUserSubscriptions(ctx, paramsCountUserSub)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error counting user_subscriptions: %s", err.Error())
	}

	userSubacriptionList := make([]*repository.UserSubscription, len(generatedUserSub))
	for idx, userSub := range generatedUserSub {
		userSubacriptionList[idx], err = generatedUserSubToRepoUserSub(generated.UserSubscription{
			ID:             userSub.ID,
			UserID:         userSub.UserID,
			SubscriptionID: userSub.SubscriptionID,
			DayOfWeek:      userSub.DayOfWeek,
			Status:         userSub.Status,
			StartDate:      userSub.StartDate,
			EndDate:        userSub.EndDate,
			DeletedAt:      userSub.DeletedAt,
			CreatedAt:      userSub.CreatedAt,
		}, userSub.UserData, userSub.SubscriptionData, nil)
		if err != nil {
			return nil, nil, err
		}
	}

	return userSubacriptionList, pkg.CalculatePagination(uint32(totalCount), filter.Pagination.PageSize, filter.Pagination.Page), nil
}

func (usr *UserSubscriptionRepository) DeleteUserSubscription(ctx context.Context, id int64) error {
	if err := usr.queries.DeleteUserSubscription(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.Errorf(pkg.NOT_FOUND_ERROR, "user_payment with ID %d not found", id)
		}
		return pkg.Errorf(pkg.INTERNAL_ERROR, "error deleting subscription by id: %s", err.Error())
	}

	return nil
}

func generatedUserSubToRepoUserSub(genUserSub generated.UserSubscription, userData, subData, paymentData []byte) (*repository.UserSubscription, error) {
	userSuscription := &repository.UserSubscription{
		ID:               uint32(genUserSub.ID),
		UserID:           0,
		SubscriptionID:   uint32(genUserSub.SubscriptionID),
		DayOfWeek:        genUserSub.DayOfWeek,
		Status:           genUserSub.Status,
		StartDate:        genUserSub.StartDate,
		EndDate:          genUserSub.EndDate,
		CreatedAt:        genUserSub.CreatedAt,
		DeletedAt:        nil,
		UserData:         nil,
		SubscriptionData: nil,
		PaymentData:      nil,
	}

	if genUserSub.UserID.Valid {
		userSuscription.UserID = uint32(genUserSub.UserID.Int64)
	}

	if genUserSub.DeletedAt.Valid {
		userSuscription.DeletedAt = &genUserSub.DeletedAt.Time
	}

	if userData != nil {
		var uData repository.User
		if err := json.Unmarshal(userData, &uData); err != nil {
			return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed unmarshalling user data to user_subscription: %s", err.Error())
		}
		userSuscription.UserData = &uData
	}

	if subData != nil {
		var sData repository.Subscription
		if err := json.Unmarshal(subData, &sData); err != nil {
			return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed unmarshalling subscription data to user_subscription: %s", err.Error())
		}
		userSuscription.SubscriptionData = &sData
	}

	if paymentData != nil {
		var pData []repository.Payment
		if err := json.Unmarshal(paymentData, &pData); err != nil {
			return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed unmarshalling payment data to user_subscription: %s", err.Error())
		}
		userSuscription.PaymentData = pData
	}

	return userSuscription, nil
}
