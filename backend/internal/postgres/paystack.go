package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/flexGURU/flower-haven/backend/internal/postgres/generated"
	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ repository.PaystackRepository = (*PaystackRepository)(nil)

type PaystackRepository struct {
	queries *generated.Queries
}

func NewPaystackRepository(queries *generated.Queries) *PaystackRepository {
	return &PaystackRepository{
		queries: queries,
	}
}

func (ps *PaystackRepository) CreatePayment(ctx context.Context, email string, amount int64, reference string) error {
	if err := ps.queries.CreatePaystackPayment(ctx, generated.CreatePaystackPaymentParams{
		Email:     email,
		Amount:    fmt.Sprintf("%d", amount),
		Reference: reference,
	}); err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create paystack payment: %s", err.Error())
	}
	return nil
}

func (ps *PaystackRepository) GetPaymentByReference(ctx context.Context, reference string) (repository.PaystackPayment, error) {
	payment, err := ps.queries.GetPaystackPaymentByReference(ctx, reference)
	if err != nil {
		return repository.PaystackPayment{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get paystack payment by reference: %s", err.Error())
	}
	return repository.PaystackPayment{
		ID:        payment.ID,
		Email:     payment.Email,
		Amount:    payment.Amount,
		Reference: payment.Reference,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}, nil
}

func (ps *PaystackRepository) UpdatePaymentStatus(ctx context.Context, reference string, status string) error {
	if err := ps.queries.UpdatePaystackPaymentStatus(ctx, generated.UpdatePaystackPaymentStatusParams{
		Status:    status,
		Reference: reference,
	}); err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update paystack payment status: %s", err.Error())
	}
	return nil
}

func (ps *PaystackRepository) ListPaystackPayments(ctx context.Context, status string, pagination *pkg.Pagination) ([]repository.PaystackPayment, *pkg.Pagination, error) {
	listParams := generated.ListPaystackPaymentsParams{
		Limit:  int32(pagination.PageSize),
		Offset: pkg.Offset(pagination.Page, pagination.PageSize),
		Status: pgtype.Text{Valid: false},
	}
	countParams := pgtype.Text{Valid: false}
	if status != "" {
		listParams.Status = pgtype.Text{Valid: true, String: status}
		countParams = pgtype.Text{Valid: true, String: status}
	}

	payments, err := ps.queries.ListPaystackPayments(ctx, listParams)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list paystack payments: %s", err.Error())
	}

	totalCount, err := ps.queries.ListCountPaystackPayments(ctx, countParams)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to count paystack payments: %s", err.Error())
	}

	result := make([]repository.PaystackPayment, len(payments))
	for i, p := range payments {
		result[i] = repository.PaystackPayment{
			ID:        p.ID,
			Email:     p.Email,
			Amount:    p.Amount,
			Reference: p.Reference,
			Status:    p.Status,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}

	return result, pkg.CalculatePagination(uint32(totalCount), pagination.PageSize, pagination.Page), nil
}

func (ps *PaystackRepository) LogPaystackEvent(ctx context.Context, event string, payload []byte) error {
	if err := ps.queries.CreatePaystackEvent(ctx, generated.CreatePaystackEventParams{
		Event: event,
		Data:  payload,
	}); err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to log paystack event: %s", err.Error())
	}
	return nil
}

func (ps *PaystackRepository) ListPaystackEvents(ctx context.Context, event string, pagination *pkg.Pagination) ([]repository.PaystackEvent, *pkg.Pagination, error) {
	listParams := generated.ListPaystackEventsParams{
		Limit:  int32(pagination.PageSize),
		Offset: pkg.Offset(pagination.Page, pagination.PageSize),
		Event:  pgtype.Text{Valid: false},
	}
	countParams := pgtype.Text{Valid: false}
	if event != "" {
		listParams.Event = pgtype.Text{Valid: true, String: event}
		countParams = pgtype.Text{Valid: true, String: event}
	}

	events, err := ps.queries.ListPaystackEvents(ctx, listParams)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list paystack events: %s", err.Error())
	}

	totalCount, err := ps.queries.ListCountPaystackEvents(ctx, countParams)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to count paystack events: %s", err.Error())
	}

	result := make([]repository.PaystackEvent, len(events))
	for i, e := range events {
		event := repository.PaystackEvent{
			ID:        e.ID,
			Event:     e.Event,
			CreatedAt: e.CreatedAt,
		}

		json.Unmarshal(e.Data, &event.Data)

		result[i] = event
	}

	return result, pkg.CalculatePagination(uint32(totalCount), pagination.PageSize, pagination.Page), nil
}
