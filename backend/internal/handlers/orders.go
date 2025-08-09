package handlers

import (
	"net/http"
	"strings"

	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createOrderReq struct {
	UserName        string  `json:"user_name" binding:"required"`
	UserPhoneNumber string  `json:"user_phone_number" binding:"required"`
	PaymentStatus   string  `json:"payment_status" binding:"required"`
	Status          string  `json:"status" binding:"required"`
	ShippingAddress *string `json:"shipping_address,omitempty" `
	Items           []struct {
		ProductID uint32 `json:"product_id" binding:"required"`
		Quantity  int32  `json:"quantity" binding:"required"`
	} `json:"items" binding:"required"`
}

func (s *Server) createOrderHandler(ctx *gin.Context) {
	var req createOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	order := &repository.Order{
		UserName:        req.UserName,
		UserPhoneNumber: req.UserPhoneNumber,
		PaymentStatus:   strings.ToLower(req.PaymentStatus) == "true",
		Status:          req.Status,
		ShippingAddress: req.ShippingAddress,
	}

	orderItem := map[uint32]int32{}
	for _, item := range req.Items {
		orderItem[item.ProductID] = item.Quantity
	}

	newOrder, err := s.repo.OrderRepository.CreateOrder(ctx, order, orderItem)
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
