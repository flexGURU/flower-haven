package postgres

import (
	"context"
	"strings"

	"github.com/flexGURU/flower-haven/backend/internal/postgres/generated"
	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ repository.PaymentRepository = (*PaymentRepository)(nil)

type PaymentRepository struct {
	queries *generated.Queries
}

func NewPaymentRepository(queries *generated.Queries) *PaymentRepository {
	return &PaymentRepository{queries: queries}
}

func (pr *PaymentRepository) CreatePayment(ctx context.Context, payment *repository.Payment) (*repository.Payment, error) {
	params := generated.CreatePaymentParams{
		Amount:             pkg.Float64ToPgTypeNumeric(payment.Amount),
		PaymentMethod:      payment.PaymentMethod,
		PaidAt:             payment.PaidAt,
		OrderID:            pgtype.Int8{Valid: false},
		Description:        pgtype.Text{Valid: false},
		UserSubscriptionID: pgtype.Int8{Valid: false},
	}

	if payment.OrderID != nil {
		if exists, _ := pr.queries.OrderExists(ctx, int64(*payment.OrderID)); !exists {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "order with id %d not found", *payment.OrderID)
		}
		params.OrderID = pgtype.Int8{Valid: true, Int64: int64(*payment.OrderID)}
	}
	if payment.UserSubscriptionID != nil {
		if exists, _ := pr.queries.UserSubscriptionExists(ctx, int64(*payment.UserSubscriptionID)); !exists {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user subscription with id %d not found", *payment.UserSubscriptionID)
		}
		params.UserSubscriptionID = pgtype.Int8{Valid: true, Int64: int64(*payment.UserSubscriptionID)}
	}
	if payment.Description != nil {
		params.Description = pgtype.Text{Valid: true, String: *payment.Description}
	}

	generatedPayment, err := pr.queries.CreatePayment(ctx, params)
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
			return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating payment: %s", err.Error())
	}

	payment.ID = uint32(generatedPayment.ID)
	payment.CreatedAt = generatedPayment.CreatedAt

	return payment, nil
}

func (pr *PaymentRepository) GetPaymentByID(ctx context.Context, id int64) (*repository.Payment, error) {
	generatedPayment, err := pr.queries.GetPaymentByID(ctx, id)
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.NOT_FOUND_ERROR {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "payment with id %d not found", id)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching payment by id: %s", err.Error())
	}

	payment := &repository.Payment{
		ID:                 uint32(generatedPayment.ID),
		Description:        nil,
		OrderID:            nil,
		UserSubscriptionID: nil,
		PaymentMethod:      generatedPayment.PaymentMethod,
		Amount:             pkg.PgTypeNumericToFloat64(generatedPayment.Amount),
		PaidAt:             generatedPayment.PaidAt,
		CreatedAt:          generatedPayment.CreatedAt,
	}

	if generatedPayment.Description.Valid {
		payment.Description = &generatedPayment.Description.String
	}
	if generatedPayment.OrderID.Valid {
		orderID := uint32(generatedPayment.OrderID.Int64)
		payment.OrderID = &orderID
	}
	if generatedPayment.UserSubscriptionID.Valid {
		userSubscriptionID := uint32(generatedPayment.UserSubscriptionID.Int64)
		payment.UserSubscriptionID = &userSubscriptionID
	}

	return payment, nil
}

func (pr *PaymentRepository) UpdatePayment(ctx context.Context, payment *repository.UpdatePayment) (*repository.Payment, error) {
	params := generated.UpdatePaymentParams{
		ID:            int64(payment.ID),
		Description:   pgtype.Text{Valid: false},
		PaymentMethod: pgtype.Text{Valid: false},
		Amount:        pgtype.Numeric{Valid: false},
		PaidAt:        pgtype.Timestamptz{Valid: false},
	}

	if payment.Description != nil {
		params.Description = pgtype.Text{
			Valid:  true,
			String: *payment.Description,
		}
	}
	if payment.PaymentMethod != nil {
		params.PaymentMethod = pgtype.Text{
			Valid:  true,
			String: *payment.PaymentMethod,
		}
	}
	if payment.Amount != nil {
		params.Amount = pkg.Float64ToPgTypeNumeric(*payment.Amount)
	}
	if payment.PaidAt != nil {
		params.PaidAt = pgtype.Timestamptz{
			Valid: true,
			Time:  *payment.PaidAt,
		}
	}

	paymentId, err := pr.queries.UpdatePayment(ctx, params)
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.NOT_FOUND_ERROR {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "payment with id %d not found", payment.ID)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error updating payment: %s", err.Error())
	}

	return pr.GetPaymentByID(ctx, paymentId)
}

func (pr *PaymentRepository) ListPayments(ctx context.Context, filter *repository.PaymentFilter) ([]*repository.Payment, *pkg.Pagination, error) {
	paramsListPayments := generated.ListPaymentsParams{
		Limit:         int32(filter.Pagination.PageSize),
		Offset:        pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
		PaymentMethod: pgtype.Text{Valid: false},
		StartDate:     pgtype.Timestamptz{Valid: false},
		EndDate:       pgtype.Timestamptz{Valid: false},
	}

	paramsCountPayments := generated.ListCountPaymentsParams{
		PaymentMethod: pgtype.Text{Valid: false},
		StartDate:     pgtype.Timestamptz{Valid: false},
		EndDate:       pgtype.Timestamptz{Valid: false},
	}

	if filter.PaymentMethod != nil {
		search := strings.ToLower(*filter.PaymentMethod)
		paramsListPayments.PaymentMethod = pgtype.Text{
			Valid:  true,
			String: "%" + search + "%",
		}
		paramsCountPayments.PaymentMethod = pgtype.Text{
			Valid:  true,
			String: "%" + search + "%",
		}
	}
	if filter.StartDate != nil && filter.EndDate != nil {
		paramsListPayments.StartDate = pgtype.Timestamptz{
			Valid: true,
			Time:  *filter.StartDate,
		}
		paramsListPayments.EndDate = pgtype.Timestamptz{
			Valid: true,
			Time:  *filter.EndDate,
		}

		paramsCountPayments.StartDate = pgtype.Timestamptz{
			Valid: true,
			Time:  *filter.StartDate,
		}
		paramsCountPayments.EndDate = pgtype.Timestamptz{
			Valid: true,
			Time:  *filter.EndDate,
		}
	}
	generatedPayments, err := pr.queries.ListPayments(ctx, paramsListPayments)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing payments: %s", err.Error())
	}

	totalCount, err := pr.queries.ListCountPayments(ctx, paramsCountPayments)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error counting payments: %s", err.Error())
	}

	payments := make([]*repository.Payment, len(generatedPayments))
	for i, gp := range generatedPayments {
		payments[i] = &repository.Payment{
			ID:                 uint32(gp.ID),
			Description:        nil,
			OrderID:            nil,
			UserSubscriptionID: nil,
			PaymentMethod:      gp.PaymentMethod,
			Amount:             pkg.PgTypeNumericToFloat64(gp.Amount),
			PaidAt:             gp.PaidAt,
			CreatedAt:          gp.CreatedAt,
		}

		if gp.Description.Valid {
			payments[i].Description = &gp.Description.String
		}
		if gp.OrderID.Valid {
			orderID := uint32(gp.OrderID.Int64)
			payments[i].OrderID = &orderID
		}
		if gp.UserSubscriptionID.Valid {
			userSubscriptionID := uint32(gp.UserSubscriptionID.Int64)
			payments[i].UserSubscriptionID = &userSubscriptionID
		}
	}

	return payments, pkg.CalculatePagination(uint32(totalCount), filter.Pagination.PageSize, filter.Pagination.Page), nil
}
