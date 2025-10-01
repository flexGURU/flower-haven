package handlers

import (
	"net/http"

	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createProductReq struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description" binding:"required"`
	Price         float64  `json:"price" binding:"required"`
	ImageUrl      []string `json:"image_url" binding:"required"`
	CategoryId    uint32   `json:"category_id" binding:"required"`
	HasStems      bool     `json:"has_stems"`
	IsMessageCard bool     `json:"is_message_card"`
	IsFlowers     bool     `json:"is_flowers"`
	IsAddOn       bool     `json:"is_add_on"`
	StockQuantity int64    `json:"stock_quantity" binding:"required"`

	// extended fields
	Stems []repository.ProductStem `json:"stems,omitempty"`
}

func (s *Server) createProductHandler(ctx *gin.Context) {
	var req createProductReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	// if has stems make price to the first stem price
	if req.HasStems && len(req.Stems) <= 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "product has stems but no stems provided")))
		return
	}
	if req.HasStems {
		req.Price = req.Stems[0].Price
	}

	product := &repository.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		ImageUrl:      req.ImageUrl,
		HasStems:      req.HasStems,
		IsMessageCard: req.IsMessageCard,
		IsFlowers:     req.IsFlowers,
		IsAddOn:       req.IsAddOn,
		CategoryID:    req.CategoryId,
		StockQuantity: req.StockQuantity,
		Stems:         req.Stems,
	}

	newProduct, err := s.repo.ProductRepository.CreateProduct(ctx, product)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": newProduct})
}

func (s *Server) updateProductHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid product ID: %s", err.Error())))
		return
	}

	var req repository.UpdateProduct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	req.ID = id

	updatedProduct, err := s.repo.ProductRepository.UpdateProduct(ctx, &req)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedProduct})
}

func (s *Server) listProductsHandler(ctx *gin.Context) {
	filter := &repository.ProductFilter{
		Pagination:    &pkg.Pagination{},
		Search:        nil,
		PriceFrom:     nil,
		PriceTo:       nil,
		CategoryIDs:   nil,
		IsMessageCard: nil,
		IsFlowers:     nil,
		IsAddOn:       nil,
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

	if category := ctx.QueryArray("category_id"); len(category) > 0 {
		categoryIds := make([]int64, len(category))
		for i, cat := range category {
			catId, err := pkg.StringToUint32(cat)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid category ID: %s should be a number", err.Error())))
				return
			}
			categoryIds[i] = int64(catId)
		}
		filter.CategoryIDs = &categoryIds
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

	if isMessageCard := ctx.Query("is_message_card"); isMessageCard != "" {
		isMessageCardBool, err := pkg.StringToBool(isMessageCard)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		filter.IsMessageCard = &isMessageCardBool
	}

	if isFlowers := ctx.Query("is_flowers"); isFlowers != "" {
		isFlowersBool, err := pkg.StringToBool(isFlowers)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		filter.IsFlowers = &isFlowersBool
	}

	if isAddOn := ctx.Query("is_add_on"); isAddOn != "" {
		isAddOnBool, err := pkg.StringToBool(isAddOn)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		filter.IsAddOn = &isAddOnBool
	}

	products, pagination, err := s.repo.ProductRepository.ListProducts(ctx, filter)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products, "pagination": pagination})
}

func (s *Server) listAddOnProductsHandler(ctx *gin.Context) {
	products, err := s.repo.ProductRepository.ListAddOns(ctx)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products})
}

func (s *Server) listMessageCardProductsHandler(ctx *gin.Context) {
	products, err := s.repo.ProductRepository.ListMessageCards(ctx)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products})
}

func (s *Server) getProductHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid product ID: %s", err.Error())))
		return
	}

	product, err := s.repo.ProductRepository.GetProductByID(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": product})
}

func (s *Server) deleteProductHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid product ID: %s", err.Error())))
		return
	}

	err = s.repo.ProductRepository.DeleteProduct(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (s *Server) listProductOrderItemsHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid product ID: %s", err.Error())))
		return
	}

	filter := &repository.OrderFilter{
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

	orderItems, pagination, err := s.repo.OrderRepository.GetOrderItemsByProductID(ctx, int64(id), filter)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": orderItems, "pagination": pagination})
}
