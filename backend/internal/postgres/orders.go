package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

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

func (or *OrderRepository) CreateOrder(ctx context.Context, order *repository.Order, orderItems []repository.OrderItem) (*repository.Order, error) {
	err := or.db.ExecTx(ctx, func(q *generated.Queries) error {
		// create order details
		createOrderParams := generated.CreateOrderParams{
			UserName:        order.UserName,
			UserPhoneNumber: order.UserPhoneNumber,
			PaymentStatus:   order.PaymentStatus,
			Status:          order.Status,
			DeliveryDate:    order.DeliveryDate,
			TimeSlot:        order.TimeSlot,
			ByAdmin:         order.ByAdmin,
			ShippingAddress: pgtype.Text{Valid: false},
		}

		if order.ShippingAddress != nil {
			createOrderParams.ShippingAddress = pgtype.Text{
				Valid:  true,
				String: *order.ShippingAddress,
			}
		}

		totalAmount := 0.0
		clientSubscriptionParams := map[int]generated.CreateSubscriptionParams{}
		clientUserSubscriptionParams := map[int]generated.CreateUserSubscriptionParams{}
		orderItemParams := make([]generated.CreateOrderItemParams, len(orderItems))
		for idx, item := range orderItems {
			product, err := q.GetProductByID(ctx, int64(item.ProductID))
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return pkg.Errorf(pkg.NOT_FOUND_ERROR, "product with ID %d not found", item.ProductID)
				}
				return pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching product by id: %s", err.Error())
			}

			if item.StemID != 0 {
				// check if stem exists
				stem, err := q.GetProductStemByID(ctx, int64(item.StemID))
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						return pkg.Errorf(pkg.NOT_FOUND_ERROR, "stem with ID %d not found", item.StemID)
					}
					return pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching stem by id: %s", err.Error())
				}
				product.Price = stem.Price
			}

			if product.StockQuantity < int64(item.Quantity) {
				return pkg.Errorf(pkg.INVALID_ERROR, "product with id %d has stock_quantity of %d and trying to make an order of stock_quantity %d. Need to add stock first", item.ProductID, product.StockQuantity, item.Quantity)
			}

			amount := pkg.PgTypeNumericToFloat64(product.Price) * float64(item.Quantity)
			totalAmount += amount

			if amount != item.Amount {
				return pkg.Errorf(pkg.INVALID_ERROR, "product with id %d has price of %.2f and trying to make an order with amount %.2f. Amount should be %.2f", item.ProductID, pkg.PgTypeNumericToFloat64(product.Price), item.Amount, amount)
			}

			_, err = q.UpdateProduct(ctx, generated.UpdateProductParams{
				ID: int64(item.ProductID),
				StockQuantity: pgtype.Int8{
					Valid: true,
					Int64: product.StockQuantity - int64(item.Quantity),
				},
			})
			if err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update product with id %d quantity: %s", item.Quantity, err.Error())
			}

			if item.PaymentMethod == "subscription" {
				// create subscription params by_admin to false
				createSubParams := generated.CreateSubscriptionParams{
					Name:        order.UserName,
					Description: fmt.Sprintf("Subscription made by %s for product %s of quantity %d on %s", order.UserName, product.Name, item.Quantity, time.Now().Format("2006-01-02 15:04:05")),
					ProductIds:  []int32{int32(item.ProductID)},
					AddOns:      []int32{},
					Price:       pkg.Float64ToPgTypeNumeric(amount),
					StemIds:     []int32{int32(item.StemID)},
					ByAdmin:     false,
				}
				clientSubscriptionParams[idx] = createSubParams

				// create user_subscription params with frequency and day_of_week
				createUserSubParams := generated.CreateUserSubscriptionParams{
					UserID:    pgtype.Int8{Valid: false},
					StartDate: time.Now(),
					EndDate:   time.Now().AddDate(0, 3, 0), // default to 3 months
					DayOfWeek: int16(order.DeliveryDate.Weekday()),
					Frequency: item.Frequency,
				}
				clientUserSubscriptionParams[idx] = createUserSubParams

				// add the orderItem params To CreateOrderItemParams
				createOrerItemParam := generated.CreateOrderItemParams{
					ProductID:     int64(item.ProductID),
					Quantity:      item.Quantity,
					Amount:        pkg.Float64ToPgTypeNumeric(amount),
					PaymentMethod: item.PaymentMethod,
					Frequency:     pgtype.Text{Valid: true, String: item.Frequency},
					StemID:        pgtype.Int8{Valid: false},
				}
				if item.StemID != 0 {
					createOrerItemParam.StemID = pgtype.Int8{Valid: true, Int64: int64(item.StemID)}
				}
				orderItemParams[idx] = createOrerItemParam
			} else {
				createOrerItemParam := generated.CreateOrderItemParams{
					ProductID:     int64(item.ProductID),
					Quantity:      item.Quantity,
					Amount:        pkg.Float64ToPgTypeNumeric(amount),
					PaymentMethod: item.PaymentMethod,
					Frequency:     pgtype.Text{Valid: false},
					StemID:        pgtype.Int8{Valid: false},
				}
				if item.StemID != 0 {
					createOrerItemParam.StemID = pgtype.Int8{Valid: true, Int64: int64(item.StemID)}
				}
				orderItemParams[idx] = createOrerItemParam
			}
		}

		createOrderParams.TotalAmount = pkg.Float64ToPgTypeNumeric(totalAmount)
		orderId, err := q.CreateOrder(ctx, createOrderParams)
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create order: %s", err.Error())
		}

		// create subscriptions with orderId as parent id
		subScriptionIds := map[int]int64{}
		for idx, subParams := range clientSubscriptionParams {
			subParams.ParentOrderID = pgtype.Int8{Valid: true, Int64: orderId}
			subscriptionId, err := q.CreateSubscription(ctx, subParams)
			if err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create subscription: %s", err.Error())
			}
			subScriptionIds[idx] = subscriptionId
		}

		// create user subsctiprions
		for idx, userSubParams := range clientUserSubscriptionParams {
			userSubParams.SubscriptionID = subScriptionIds[idx]
			_, err := q.CreateUserSubscription(ctx, userSubParams)
			if err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create user subscription: %s", err.Error())
			}
		}

		// create order items
		for _, param := range orderItemParams {
			param.OrderID = orderId
			if _, err := q.CreateOrderItem(ctx, param); err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create order item: %s", err.Error())
			}
		}

		order.ID = uint32(orderId)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return or.GetOrderByID(ctx, int64(order.ID))
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
