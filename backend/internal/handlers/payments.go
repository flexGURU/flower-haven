package handlers

import (
	"net/http"
	"time"

	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createPaymentReq struct {
	OrderID        *uint32 `json:"order_id"`
	SubscriptionId *uint32 `json:"subscription_id"`
	Description    *string `json:"description"`
	PaymentMethod  string  `json:"payment_method" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
	PaidAt         *string `json:"paid_at"`
}

func (s *Server) createPaymentHandler(ctx *gin.Context) {
	var req createPaymentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	payment := &repository.Payment{
		OrderID:            req.OrderID,
		UserSubscriptionID: req.SubscriptionId,
		Description:        req.Description,
		PaymentMethod:      req.PaymentMethod,
		Amount:             req.Amount,
	}

	if req.PaidAt != nil {
		payment.PaidAt = pkg.StringToTime(*req.PaidAt)
	} else {
		payment.PaidAt = time.Now()
	}

	newPayment, err := s.repo.PaymentRepository.CreatePayment(ctx, payment)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": newPayment})
}

func (s *Server) getPaymentHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid payment ID: %s", err.Error())))
		return
	}

	payment, err := s.repo.PaymentRepository.GetPaymentByID(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": payment})
}

type updatePaymentReq struct {
	Description   *string  `json:"description,omitempty"`
	PaymentMethod *string  `json:"payment_method"`
	Amount        *float64 `json:"amount"`
	PaidAt        *string  `json:"paid_at"`
}

func (s *Server) updatePaymentHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid payment ID: %s", err.Error())))
		return
	}

	var req updatePaymentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	params := &repository.UpdatePayment{
		ID:            id,
		Description:   nil,
		PaymentMethod: nil,
		Amount:        nil,
		PaidAt:        nil,
	}

	if req.Description != nil {
		params.Description = req.Description
	}
	if req.PaymentMethod != nil {
		params.PaymentMethod = req.PaymentMethod
	}
	if req.Amount != nil {
		params.Amount = req.Amount
	}
	if req.PaidAt != nil {
		parsedTime := pkg.StringToTime(*req.PaidAt)
		params.PaidAt = &parsedTime
	}

	updatedPayment, err := s.repo.PaymentRepository.UpdatePayment(ctx, params)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedPayment})
}

func (s *Server) listPaymentsHandler(ctx *gin.Context) {
	filter := &repository.PaymentFilter{
		Pagination:    &pkg.Pagination{},
		PaymentMethod: nil,
		StartDate:     nil,
		EndDate:       nil,
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

	if paymentMethod := ctx.Query("payment_method"); paymentMethod != "" {
		filter.PaymentMethod = &paymentMethod
	}

	startDateStr := ctx.DefaultQuery("start_date", "2006-01-02")
	startDate := pkg.StringToTime(startDateStr)
	filter.StartDate = &startDate

	endDateStr := ctx.DefaultQuery("end_date", time.Now().Format("2030-01-02"))
	endDate := pkg.StringToTime(endDateStr)
	endDate = endDate.Add(24 * time.Hour)
	filter.EndDate = &endDate

	payments, pagination, err := s.repo.PaymentRepository.ListPayments(ctx, filter)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": payments, "pagination": pagination})
}
