package handlers

import (
	"net/http"

	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createUserSubscriptionReq struct {
	UserId         uint32 `json:"user_id" binding:"required"`
	SubscriptionId uint32 `json:"subscription_id" binding:"required"`
	StartDate      string `json:"start_date" binding:"required"`
	EndDate        string `json:"end_date" binding:"required"`
	DayOfWeek      int16  `json:"day_of_week" binding:"required"`
}

func (s *Server) createUserSubscriptionHandler(ctx *gin.Context) {
	var req createUserSubscriptionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	subscription := &repository.UserSubscription{
		UserID:         req.UserId,
		SubscriptionID: req.SubscriptionId,
		DayOfWeek:      req.DayOfWeek,
		StartDate:      pkg.StringToTime(req.StartDate),
		EndDate:        pkg.StringToTime(req.EndDate),
	}

	if subscription.DayOfWeek < 0 || subscription.DayOfWeek > 6 {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "day of week must be between 0 (Sunday) and 6 (Saturday)")))
		return
	}

	newSubscription, err := s.repo.UserSubscriptionRepository.CreateUserSubscription(ctx, subscription)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": newSubscription})
}

type updateUserSubscriptionReq struct {
	DayOfWeek *int16  `json:"day_of_week"`
	Status    *string `json:"status"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

func (s *Server) updateUserSubscriptionHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid user subscription ID: %s", err.Error())))
		return
	}

	var req updateUserSubscriptionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	if req.DayOfWeek != nil && (*req.DayOfWeek < 0 || *req.DayOfWeek > 6) {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "day of week must be between 0 (Sunday) and 6 (Saturday)")))
		return
	}

	subscription := &repository.UpdateUserSubscription{
		ID:        id,
		DayOfWeek: nil,
		Status:    nil,
		StartDate: nil,
		EndDate:   nil,
	}

	if req.DayOfWeek != nil {
		subscription.DayOfWeek = req.DayOfWeek
	}
	if req.Status != nil {
		status, err := pkg.StringToBool(*req.Status)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid status value: %s", err.Error())))
			return
		}
		subscription.Status = &status
	}
	if req.StartDate != nil {
		d := pkg.StringToTime(*req.StartDate)
		subscription.StartDate = &d
	}
	if req.EndDate != nil {
		d := pkg.StringToTime(*req.EndDate)
		subscription.EndDate = &d
	}

	updatedSubscription, err := s.repo.UserSubscriptionRepository.UpdateUserSubscription(ctx, subscription)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedSubscription})
}

func (s *Server) listUserSubscriptionsHandler(ctx *gin.Context) {
	filter := &repository.UserSubscriptionFilter{
		Pagination: &pkg.Pagination{},
		Status:     nil,
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

	if status := ctx.Query("status"); status != "" {
		statusBool, err := pkg.StringToBool(status)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		filter.Status = &statusBool
	}

	subscriptions, pagination, err := s.repo.UserSubscriptionRepository.ListUserSubscriptions(ctx, filter)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": subscriptions, "pagination": pagination})
}

func (s *Server) getUserSubscriptionsHandler(ctx *gin.Context) {
	userId, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid user ID: %s", err.Error())))
		return
	}

	filter := &repository.UserSubscriptionFilter{
		Pagination: &pkg.Pagination{},
		Status:     nil,
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

	subscriptions, pagination, err := s.repo.UserSubscriptionRepository.GetUsersSubscriptionsByUserID(ctx, int64(userId), filter)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": subscriptions, "pagination": pagination})
}

func (s *Server) getUserSubscriptionHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid user ID: %s", err.Error())))
		return
	}

	subscriptions, err := s.repo.UserSubscriptionRepository.GetUserSubscriptionByID(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": subscriptions})
}

func (s *Server) deleteUserSubscriptionHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid subscription ID: %s", err.Error())))
		return
	}

	err = s.repo.UserSubscriptionRepository.DeleteUserSubscription(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User subscription deleted successfully"})
}
