package handlers

import (
	"net/http"

	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createSubscriptionReq struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	ProductIds  []uint32 `json:"product_ids" binding:"required"`
	AddOns      []uint32 `json:"add_ons"`
	Price       float64  `json:"price" binding:"required"`
}

func (s *Server) createSubscriptionHandler(ctx *gin.Context) {
	var req createSubscriptionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	subscription := &repository.Subscription{
		Name:        req.Name,
		Description: req.Description,
		ProductIds:  req.ProductIds,
		AddOns:      req.AddOns,
		Price:       req.Price,
	}

	if len(subscription.ProductIds) == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "at least one product ID is required")))
		return
	}

	newSubscription, err := s.repo.SubscriptionRepository.CreateSubscription(ctx, subscription)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": newSubscription})
}

func (s *Server) getSubscriptionHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid subscription ID: %s", err.Error())))
		return
	}

	subscription, err := s.repo.SubscriptionRepository.GetSubscriptionByID(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": subscription})
}

func (s *Server) listSubscriptionsHandler(ctx *gin.Context) {
	filter := &repository.SubscriptionFilter{
		Pagination: &pkg.Pagination{},
		Search:     nil,
		PriceFrom:  nil,
		PriceTo:    nil,
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

	if priceFrom := ctx.Query("price_from"); priceFrom != "" {
		priceFromFloat, err := pkg.StringToFloat64(priceFrom)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		filter.PriceFrom = &priceFromFloat
	}

	if priceTo := ctx.Query("price_to"); priceTo != "" {
		priceToFloat, err := pkg.StringToFloat64(priceTo)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		filter.PriceTo = &priceToFloat
	}

	subscriptions, pagination, err := s.repo.SubscriptionRepository.ListSubscriptions(ctx, filter)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       subscriptions,
		"pagination": pagination,
	})
}

func (s *Server) updateSubscriptionHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid subscription ID: %s", err.Error())))
		return
	}

	var req repository.UpdateSubscription
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}
	req.ID = id

	if req.ProductIds != nil && len(*req.ProductIds) == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "at least one product ID is required")))
		return
	}

	updatedSubscription, err := s.repo.SubscriptionRepository.UpdateSubscription(ctx, &req)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedSubscription})
}

func (s *Server) deleteSubscriptionHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid subscription ID: %s", err.Error())))
		return
	}

	err = s.repo.SubscriptionRepository.DeleteSubscription(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription deleted successfully"})
}
