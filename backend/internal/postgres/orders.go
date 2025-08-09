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

var _ repository.OrderRepository = (*OrderRepository)(nil)

type OrderRepository struct {
	queries *generated.Queries
	db      *Store
}

func NewOrderRepository(db *Store) *OrderRepository {
	return &OrderRepository{
		db:      db,
		queries: generated.New(db.pool),
	}
}

func (or *OrderRepository) CreateOrder(ctx context.Context, order *repository.Order, orderItems map[uint32]int32) (*repository.Order, error) {
	err := or.db.ExecTx(ctx, func(q *generated.Queries) error {
		createOrderParams := generated.CreateOrderParams{
			UserName:        order.UserName,
			UserPhoneNumber: order.UserPhoneNumber,
			PaymentStatus:   order.PaymentStatus,
			Status:          order.Status,
			ShippingAddress: pgtype.Text{Valid: false},
		}

		if order.ShippingAddress != nil {
			createOrderParams.ShippingAddress = pgtype.Text{
				Valid:  true,
				String: *order.ShippingAddress,
			}
		}

		totalAmount := 0.0
		orderItemParams := []generated.CreateOrderItemParams{}
		orderItemResult := []repository.OrderItem{}
		for productId, quantity := range orderItems {
			product, err := q.GetProductByID(ctx, int64(productId))
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return pkg.Errorf(pkg.NOT_FOUND_ERROR, "product with ID %d not found", productId)
				}
				return pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching product by id: %s", err.Error())
			}

			if product.StockQuantity < int64(quantity) {
				return pkg.Errorf(pkg.INVALID_ERROR, "product with id %d has stock_quantity of %d and trying to make an order of stock_quantity %d. Need to add stock first", productId, product.StockQuantity, quantity)
			}

			amount := pkg.PgTypeNumericToFloat64(product.Price) * float64(quantity)
			totalAmount += amount
			orderItemParams = append(orderItemParams, generated.CreateOrderItemParams{
				ProductID: int64(productId),
				Quantity:  quantity,
				Amount:    pkg.Float64ToPgTypeNumeric(amount),
			})

			_, err = q.UpdateProduct(ctx, generated.UpdateProductParams{
				ID: int64(productId),
				StockQuantity: pgtype.Int8{
					Valid: true,
					Int64: product.StockQuantity - int64(quantity),
				},
			})
			if err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update product with id %d quantity: %s", productId, err.Error())
			}

			orderItemResult = append(orderItemResult, repository.OrderItem{
				ProductID: productId,
				Quantity:  quantity,
				Amount:    amount,
				CurrentProductDetails: &repository.Product{
					ID:       uint32(product.ID),
					Name:     product.Description,
					Price:    pkg.PgTypeNumericToFloat64(product.Price),
					ImageUrl: product.ImageUrl,
				},
				OrderData: nil,
			})
		}

		createOrderParams.TotalAmount = pkg.Float64ToPgTypeNumeric(totalAmount)
		orderId, err := q.CreateOrder(ctx, createOrderParams)
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create order: %s", err.Error())
		}

		for idx, param := range orderItemParams {
			param.OrderID = orderId
			orderItemResult[idx].OrderID = uint32(orderId)
			if _, err := q.CreateOrderItem(ctx, param); err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create order item: %s", err.Error())
			}
		}

		order.ID = uint32(orderId)
		order.TotalAmount = totalAmount
		order.OrderItemsData = orderItemResult

		return nil
	})

	return order, err
}

func (or *OrderRepository) GetOrderByID(ctx context.Context, id int64) (*repository.Order, error) {
	order, err := or.queries.GetOrderByFullDataID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "order with ID %d not found", id)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching order by id: %s", err.Error())
	}

	orderItemData := []repository.OrderItem{}
	if err := json.Unmarshal(order.OrderItemData, &orderItemData); err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error unmashaling order_item_data to order: %s", err.Error())
	}

	rslt := &repository.Order{
		ID:              uint32(order.ID),
		UserName:        order.UserName,
		UserPhoneNumber: order.UserPhoneNumber,
		TotalAmount:     pkg.PgTypeNumericToFloat64(order.TotalAmount),
		PaymentStatus:   order.PaymentStatus,
		Status:          order.Status,
		ShippingAddress: nil,
		DeletedAt:       nil,
		CreatedAt:       order.CreatedAt,
		OrderItemsData:  orderItemData,
	}

	if order.ShippingAddress.Valid {
		rslt.ShippingAddress = &order.ShippingAddress.String
	}

	if order.DeletedAt.Valid {
		rslt.DeletedAt = &order.DeletedAt.Time
	}

	return rslt, nil
}

func (or *OrderRepository) UpdateOrder(ctx context.Context, order *repository.UpdateOrder) (*repository.Order, error) {
	params := generated.UpdateOrderParams{
		ID:              int64(order.ID),
		UserName:        pgtype.Text{Valid: false},
		UserPhoneNumber: pgtype.Text{Valid: false},
		PaymentStatus:   pgtype.Bool{Valid: false},
		Status:          pgtype.Text{Valid: false},
		ShippingAddress: pgtype.Text{Valid: false},
	}

	if order.UserName != nil {
		params.UserName = pgtype.Text{
			Valid:  true,
			String: *order.UserName,
		}
	}

	if order.UserPhoneNumber != nil {
		params.UserPhoneNumber = pgtype.Text{
			Valid:  true,
			String: *order.UserPhoneNumber,
		}
	}

	if order.PaymentStatus != nil {
		params.PaymentStatus = pgtype.Bool{
			Valid: true,
			Bool:  *order.PaymentStatus,
		}
	}

	if order.Status != nil {
		params.Status = pgtype.Text{
			Valid:  true,
			String: *order.Status,
		}
	}
	if order.ShippingAddress != nil {
		params.ShippingAddress = pgtype.Text{
			Valid:  true,
			String: *order.ShippingAddress,
		}
	}

	orderId, err := or.queries.UpdateOrder(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "order with ID %d not found", order.ID)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error updating order: %s", err.Error())
	}

	return or.GetOrderByID(ctx, orderId)
}

func (or *OrderRepository) ListOrders(ctx context.Context, filter *repository.OrderFilter) ([]*repository.Order, *pkg.Pagination, error) {
	paramsListOrders := generated.ListOrderParams{
		Limit:         int32(filter.Pagination.PageSize),
		Offset:        pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
		Search:        pgtype.Text{Valid: false},
		PaymentStatus: pgtype.Bool{Valid: false},
		Status:        pgtype.Text{Valid: false},
	}

	paramsCountOrders := generated.ListCountOrderParams{
		Search:        pgtype.Text{Valid: false},
		PaymentStatus: pgtype.Bool{Valid: false},
		Status:        pgtype.Text{Valid: false},
	}

	if filter.Search != nil {
		search := strings.ToLower(*filter.Search)
		paramsListOrders.Search = pgtype.Text{
			Valid:  true,
			String: "%" + search + "%",
		}
		paramsCountOrders.Search = pgtype.Text{
			Valid:  true,
			String: "%" + search + "%",
		}
	}

	if filter.PaymentStatus != nil {
		paramsListOrders.PaymentStatus = pgtype.Bool{
			Valid: true,
			Bool:  *filter.PaymentStatus,
		}
		paramsCountOrders.PaymentStatus = pgtype.Bool{
			Valid: true,
			Bool:  *filter.PaymentStatus,
		}
	}

	if filter.Status != nil {
		status := strings.ToLower(*filter.Status)
		paramsListOrders.Status = pgtype.Text{
			Valid:  true,
			String: "%" + status + "%",
		}
		paramsCountOrders.Status = pgtype.Text{
			Valid:  true,
			String: "%" + status + "%",
		}
	}

	generatedOrders, err := or.queries.ListOrder(ctx, paramsListOrders)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing orders: %s", err.Error())
	}

	totalCount, err := or.queries.ListCountOrder(ctx, paramsCountOrders)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error counting orders: %s", err.Error())
	}

	orders := make([]*repository.Order, len(generatedOrders))
	for i, order := range generatedOrders {
		orders[i] = &repository.Order{
			ID:              uint32(order.ID),
			UserName:        order.UserName,
			UserPhoneNumber: order.UserPhoneNumber,
			TotalAmount:     pkg.PgTypeNumericToFloat64(order.TotalAmount),
			PaymentStatus:   order.PaymentStatus,
			Status:          order.Status,
			ShippingAddress: &order.ShippingAddress.String,
			DeletedAt:       &order.DeletedAt.Time,
			CreatedAt:       order.CreatedAt,
		}
	}

	return orders, pkg.CalculatePagination(uint32(totalCount), filter.Pagination.PageSize, filter.Pagination.Page), nil
}

func (or *OrderRepository) DeleteOrder(ctx context.Context, id int64) error {
	if err := or.queries.DeleteOrder(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.Errorf(pkg.NOT_FOUND_ERROR, "order with ID %d not found", id)
		}
		return pkg.Errorf(pkg.INTERNAL_ERROR, "error deleting order by id: %s", err.Error())
	}

	return nil
}

func (or *OrderRepository) GetOrderItemsByProductID(ctx context.Context, productID int64, filter *repository.OrderFilter) ([]*repository.OrderItem, *pkg.Pagination, error) {
	orderItems, err := or.queries.GetOrderItemsByProductID(ctx, generated.GetOrderItemsByProductIDParams{
		ProductID: productID,
		Limit:     int32(filter.Pagination.PageSize),
		Offset:    pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "order items with product id %d not found", productID)
		}
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching order items by product id: %s", err.Error())
	}

	totalCount, err := or.queries.GetCountOrderItemsByProductID(ctx, productID)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error counting order items by product id: %s", err.Error())
	}

	orderItemList := make([]*repository.OrderItem, len(orderItems))
	for i, item := range orderItems {
		var orderData repository.Order
		if err := json.Unmarshal(item.OrderData, &orderData); err != nil {
			return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error unmarshaling order data: %s", err.Error())
		}

		orderItemList[i] = &repository.OrderItem{
			ID:                    uint32(item.ID),
			OrderID:               uint32(item.OrderID),
			ProductID:             uint32(item.ProductID),
			Quantity:              item.Quantity,
			Amount:                pkg.PgTypeNumericToFloat64(item.Amount),
			OrderData:             &orderData,
			CurrentProductDetails: nil,
		}
	}

	return orderItemList, pkg.CalculatePagination(uint32(totalCount), filter.Pagination.PageSize, filter.Pagination.Page), nil
}
