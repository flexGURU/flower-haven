package handlers

import (
	"net/http"
	"time"

	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createSubscriptionDeliveryReq struct {
	UserSubscriptionID uint32    `json:"user_subscription_id" binding:"required"`
	DeliveredOn        time.Time `json:"delivered_on" binding:"required"`
	Description        *string   `json:"description,omitempty"`
}

func (s *Server) createSubscriptionDeliveryHandler(ctx *gin.Context) {
	var req createSubscriptionDeliveryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	delivery := &repository.SubscriptionDelivery{
		UserSubscriptionID: req.UserSubscriptionID,
		DeliveredOn:        req.DeliveredOn,
		Description:        req.Description,
	}

	newDelivery, err := s.repo.SubscriptionDeliveryRepository.CreateSubscriptionDelivery(ctx, delivery)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": newDelivery})
}

func (s *Server) listSubscriptionDeliveriesHandler(ctx *gin.Context) {
	filter := &repository.SubscriptionDeliveryFilter{
		Pagination: &pkg.Pagination{},
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

	deliveries, pagination, err := s.repo.SubscriptionDeliveryRepository.ListSubscriptionDeliveries(ctx, filter)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": deliveries, "pagination": pagination})
}

func (s *Server) updateSubscriptionDeliveryHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid subscription delivery ID: %s", err.Error())))
		return
	}

	var req repository.UpdateSubscriptionDelivery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	req.ID = id

	updatedDelivery, err := s.repo.SubscriptionDeliveryRepository.UpdateSubscriptionDelivery(ctx, &req)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedDelivery})
}

func (s *Server) getSubscriptionDeliveryByUserSubscriptionIDHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid subscription delivery ID: %s", err.Error())))
		return
	}

	delivery, err := s.repo.SubscriptionDeliveryRepository.GetSubscriptonDeliveryByUserSubscriptionID(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": delivery})
}

func (s *Server) deleteSubscriptionDeliveryHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid subscription delivery ID: %s", err.Error())))
		return
	}

	err = s.repo.SubscriptionDeliveryRepository.DeleteSubscriptionDelivery(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription delivery deleted successfully"})
}
