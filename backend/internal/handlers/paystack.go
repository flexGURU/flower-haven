package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/gin-gonic/gin"
)

type initializePaystackPaymentReq struct {
	Email  string  `json:"email" binding:"required,email"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (s *Server) initializePaystackPayment(ctx *gin.Context) {
	var req initializePaystackPaymentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	amount := int64(req.Amount * 100) // convert to kobo

	accessCode, reference, err := s.ps.InitializePayment(req.Email, amount)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	if err := s.repo.PaystackRepository.CreatePayment(ctx, req.Email, amount, reference); err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_code": accessCode, "reference": reference})
}

type paystackWebhookReq struct {
	Event string      `json:"event" binding:"required"`
	Data  interface{} `json:"data" binding:"required"`
}

func (s *Server) handlePaystackWebhook(ctx *gin.Context) {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "could not read request body")))
		return
	}

	// Validate signature
	signature := ctx.GetHeader("x-paystack-signature")
	mac := hmac.New(sha512.New, []byte(s.config.PAYSTACK_SECRET_KEY))
	mac.Write(bodyBytes)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	if signature != expectedMAC {
		ctx.JSON(http.StatusUnauthorized, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid webhook signature")))
		return
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req paystackWebhookReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	payload, err := json.Marshal(req.Data)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	if err := s.repo.PaystackRepository.LogPaystackEvent(ctx, req.Event, payload); err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "event logged and payment status updated"})
}

func (s *Server) getPaystackPayment(ctx *gin.Context) {
	reference := ctx.Param("reference")
	if reference == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "reference is required")))
		return
	}

	payment, err := s.repo.PaystackRepository.GetPaymentByReference(ctx, reference)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": payment})
}

func (s *Server) listPaystackPayments(ctx *gin.Context) {
	pageNoStr := ctx.DefaultQuery("page", "1")
	pageNo, err := pkg.StringToUint32(pageNoStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	pageSizeStr := ctx.DefaultQuery("limit", "10")
	pageSize, err := pkg.StringToUint32(pageSizeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	status := ctx.Query("status")

	payments, pagination, err := s.repo.PaystackRepository.ListPaystackPayments(ctx, status, &pkg.Pagination{
		Page:     pageNo,
		PageSize: pageSize,
	})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": payments, "pagination": pagination})
}

func (s *Server) listPaystackEvents(ctx *gin.Context) {
	pageNoStr := ctx.DefaultQuery("page", "1")
	pageNo, err := pkg.StringToUint32(pageNoStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	pageSizeStr := ctx.DefaultQuery("limit", "10")
	pageSize, err := pkg.StringToUint32(pageSizeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	eventType := ctx.Query("event")

	events, pagination, err := s.repo.PaystackRepository.ListPaystackEvents(ctx, eventType, &pkg.Pagination{
		Page:     pageNo,
		PageSize: pageSize,
	})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": events, "pagination": pagination})
}
