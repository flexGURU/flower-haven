package handlers

import (
	"net/http"
	"time"

	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createOrderReq struct {
	UserName        string  `json:"user_name" binding:"required"`
	UserPhoneNumber string  `json:"user_phone_number" binding:"required"`
	PaymentStatus   bool    `json:"payment_status"`
	Status          string  `json:"status" binding:"required"`
	DeliveryDate    string  `json:"delivery_date" binding:"required"` // parse into time.Time
	TimeSlot        string  `json:"time_slot" binding:"required"`
	ShippingAddress *string `json:"shipping_address,omitempty"`
	ByAdmin         bool    `json:"by_admin"`

	Items []struct {
		ProductID     uint32  `json:"product_id" binding:"required"`
		StemID        *uint32 `json:"stem_id,omitempty"` // flowers only
		PaymentMethod string  `json:"payment_method" binding:"required,oneof=normal subscription"`
		Frequency     string  `json:"frequency,omitempty"` // required if subscription
		Quantity      int32   `json:"quantity" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
	} `json:"items" binding:"required"`
}

func (s *Server) createOrderHandler(ctx *gin.Context) {
	var req createOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	deliveryDate, err := time.Parse("2006-01-02", req.DeliveryDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid delivery_date format, expected YYYY-MM-DD")))
		return
	}

	order := &repository.Order{
		UserName:        req.UserName,
		UserPhoneNumber: req.UserPhoneNumber,
		PaymentStatus:   req.PaymentStatus,
		Status:          req.Status,
		DeliveryDate:    deliveryDate,
		TimeSlot:        req.TimeSlot,
		ByAdmin:         req.ByAdmin,
		ShippingAddress: req.ShippingAddress,
	}

	// Build order items
	var orderItems []repository.OrderItem
	for _, item := range req.Items {
		// Validation: subscription must have frequency
		if item.PaymentMethod == "subscription" && item.Frequency == "" {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "subscription items must include frequency")))
			return
		}

		orderItem := repository.OrderItem{
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			Amount:        item.Amount,
			PaymentMethod: item.PaymentMethod,
			Frequency:     item.Frequency,
		}

		if item.StemID != nil {
			orderItem.StemID = *item.StemID
		}

		orderItems = append(orderItems, orderItem)
	}

	newOrder, err := s.repo.OrderRepository.CreateOrder(ctx, order, orderItems)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": newOrder})
}

func (s *Server) listOrdersHandler(ctx *gin.Context) {
	filter := &repository.OrderFilter{
		Pagination:    &pkg.Pagination{},
		Search:        nil,
		PaymentStatus: nil,
		Status:        nil,
	}

	pageNoStr := ctx.DefaultQuery("page", "1")
	pageNo, err := pkg.StringToUint32(pageNoStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}
	filter.Pagination.Page = pageNo

	pageSizeStr := ctx.DefaultQuery("limit", "10")
	pageSize, err := pkg.StringToUint32(pageSizeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}
	filter.Pagination.PageSize = pageSize

	if search := ctx.Query("search"); search != "" {
		filter.Search = &search
	}

	if status := ctx.Query("status"); status != "" {
		filter.Status = &status
	}

	if paymentStatus := ctx.Query("payment_status"); paymentStatus != "" {
		paymentStatusBool, err := pkg.StringToBool(paymentStatus)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		filter.PaymentStatus = &paymentStatusBool
	}

	orders, pagination, err := s.repo.OrderRepository.ListOrders(ctx, filter)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       orders,
		"pagination": pagination,
	})
}

type updateOrderReq struct {
	UserName        *string `json:"user_name"`
	UserPhoneNumber *string `json:"user_phone_number"`
	PaymentStatus   *string `json:"payment_status"`
	Status          *string `json:"status"`
	ShippingAddress *string `json:"shipping_address,omitempty"`
}

func (s *Server) updateOrderHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid order ID: %s", err.Error())))
		return
	}

	var req updateOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	params := repository.UpdateOrder{
		ID:              id,
		UserName:        nil,
		UserPhoneNumber: nil,
		PaymentStatus:   nil,
		Status:          nil,
		ShippingAddress: nil,
	}

	if req.UserName != nil {
		params.UserName = req.UserName
	}
	if req.UserPhoneNumber != nil {
		params.UserPhoneNumber = req.UserPhoneNumber
	}
	if req.PaymentStatus != nil {
		paymentStatus, err := pkg.StringToBool(*req.PaymentStatus)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid payment status: %s", err.Error())))
			return
		}
		params.PaymentStatus = &paymentStatus
	}
	if req.Status != nil {
		params.Status = req.Status
	}
	if req.ShippingAddress != nil {
		params.ShippingAddress = req.ShippingAddress
	}

	updatedOrder, err := s.repo.OrderRepository.UpdateOrder(ctx, &params)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedOrder})
}

func (s *Server) getOrderHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid order ID: %s", err.Error())))
		return
	}

	order, err := s.repo.OrderRepository.GetOrderByID(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": order})
}

func (s *Server) deleteOrderHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid order ID: %s", err.Error())))
		return
	}

	err = s.repo.OrderRepository.DeleteOrder(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
