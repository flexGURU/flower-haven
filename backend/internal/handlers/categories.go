package handlers

import (
	"net/http"

	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createCategoryReq struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	ImageUrl    []string `json:"image_url" binding:"required"`
}

func (s *Server) createCategoryHandler(ctx *gin.Context) {
	var req createCategoryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	category := &repository.Category{
		Name:        req.Name,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
	}

	newCategory, err := s.repo.CategoryRepository.CreateCategory(ctx, category)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": newCategory})
}

func (s *Server) getCategoryHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid category ID: %s", err.Error())))
		return
	}

	category, err := s.repo.CategoryRepository.GetCategoryByID(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": category})
}

func (s *Server) listCategoriesHandler(ctx *gin.Context) {
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

	filter := &repository.CategoryFilter{
		Pagination: &pkg.Pagination{Page: pageNo, PageSize: pageSize},
		Search:     nil,
	}
	if search := ctx.Query("search"); search != "" {
		filter.Search = &search
	}

	categories, pagination, err := s.repo.CategoryRepository.ListCategories(ctx, filter)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       categories,
		"pagination": pagination,
	})
}

func (s *Server) updateCategoryHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid category ID: %s", err.Error())))
		return
	}

	var req repository.UpdateCategory
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}
	req.ID = id

	updatedCategory, err := s.repo.CategoryRepository.UpdateCategory(ctx, &req)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedCategory})
}

func (s *Server) deleteCategoryHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid category ID: %s", err.Error())))
		return
	}

	err = s.repo.CategoryRepository.DeleteCategory(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
