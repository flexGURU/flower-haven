package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createUserReq struct {
	Name        string  `json:"name" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	PhoneNumber string  `json:"phone_number" binding:"required"`
	Password    string  `json:"password" binding:"required"`
	Address     *string `json:"address,omitempty"`
	IsAdmin     string  `json:"is_admin"     binding:"required"`
}

func (s *Server) createUserHandler(ctx *gin.Context) {
	var req createUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	hashPassword, err := pkg.GenerateHashPassword(req.Password, s.config.PASSWORD_COST)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	defaultRefreshToken := "default_token_value"
	userPayload := &repository.User{
		Name:         req.Name,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		Password:     &hashPassword,
		IsAdmin:      strings.ToLower(req.IsAdmin) == "true",
		RefreshToken: &defaultRefreshToken,
		Address:      req.Address,
	}

	user, err := s.repo.UserRepository.CreateUser(ctx, userPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Name,
		user.Email,
		s.config.REFRESH_TOKEN_DURATION,
		user.IsAdmin,
	)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	accessToken, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Name,
		user.Email,
		s.config.TOKEN_DURATION,
		user.IsAdmin,
	)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user, "auth": gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}})
}

func (s *Server) getUserHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	user, err := s.repo.UserRepository.GetUserByID(ctx, int64(id))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func (s *Server) listUsersHandler(ctx *gin.Context) {
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

	filter := &repository.UserFilter{
		Pagination: &pkg.Pagination{
			Page:     pageNo,
			PageSize: pageSize,
		},
		Search:   nil,
		IsAdmin:  nil,
		IsActive: nil,
	}

	if search := ctx.Query("search"); search != "" {
		filter.Search = &search
	}

	if isAdmin := ctx.Query("is_admin"); isAdmin != "" {
		isAdminBool, err := pkg.StringToBool(isAdmin)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		filter.IsAdmin = &isAdminBool
	}

	if isActive := ctx.Query("is_active"); isActive != "" {
		isActiveBool, err := pkg.StringToBool(isActive)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		filter.IsActive = &isActiveBool
	}

	users, pagination, err := s.repo.UserRepository.ListUsers(ctx, filter)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users, "pagination": pagination})
}

func (s *Server) updateUserHandler(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	var req repository.UpdateUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}
	req.ID = id
	// not updating in this endpoint
	req.Password = nil
	req.RefreshToken = nil

	if req.IsAdmin != nil {
		authPayload, ok := ctx.Get(authorizationPayloadKey)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "authorization payload not found")))
			return
		}

		payload, ok := authPayload.(*pkg.Payload)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid authorization payload")))
			return
		}

		if !payload.IsAdmin {
			ctx.JSON(http.StatusForbidden, errorResponse(pkg.Errorf(pkg.FORBIDDEN_ERROR, "only admin can update is_admin field")))
			return
		}
	}

	if req.IsActive != nil {
		authPayload, ok := ctx.Get(authorizationPayloadKey)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "authorization payload not found")))
			return
		}

		payload, ok := authPayload.(*pkg.Payload)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid authorization payload")))
			return
		}

		if !payload.IsAdmin {
			ctx.JSON(http.StatusForbidden, errorResponse(pkg.Errorf(pkg.FORBIDDEN_ERROR, "only admin can update is_active field")))
			return
		}
	}

	user, err := s.repo.UserRepository.UpdateUser(ctx, &req)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

type loginReq struct {
	Email    string `binding:"required" json:"email"`
	Password string `binding:"required" json:"password"`
}

func (s *Server) login(ctx *gin.Context) {
	var req loginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	user, err := s.repo.UserRepository.GetUserInternal(ctx, 0, req.Email)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	if err := pkg.ComparePasswordAndHash(*user.Password, req.Password); err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	refreshToken, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Name,
		user.Email,
		s.config.REFRESH_TOKEN_DURATION,
		user.IsAdmin,
	)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	accessToken, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Name,
		user.Email,
		s.config.TOKEN_DURATION,
		user.IsAdmin,
	)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	// ctx.SetCookie(
	// 	"refreshToken",
	// 	refreshToken,
	// 	int(s.config.REFRESH_TOKEN_DURATION),
	// 	"/",
	// 	"",
	// 	true,
	// 	true,
	// )

	updatedUser, err := s.repo.UserRepository.UpdateUser(ctx, &repository.UpdateUser{
		ID:           user.ID,
		RefreshToken: &refreshToken,
	})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          updatedUser,
	})
}

func (s *Server) logout(ctx *gin.Context) {
	ctx.SetCookie("refreshToken", "", -1, "/", "", true, true)
	ctx.JSON(http.StatusOK, gin.H{"data": "success"})
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *Server) refreshToken(ctx *gin.Context) {
	var req refreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	payload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	accesstoken, err := s.tokenMaker.CreateToken(
		payload.UserID,
		payload.Name,
		payload.Email,
		s.config.TOKEN_DURATION,
		payload.IsAdmin,
	)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":            accesstoken,
		"access_token_expires_at": time.Now().Add(s.config.TOKEN_DURATION),
	})
}
