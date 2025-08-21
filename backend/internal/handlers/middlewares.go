package handlers

import (
	"net/http"
	"strings"

	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey        = "Authorization"
	authorizationHeaderBearerType = "bearer"
	authorizationPayloadKey       = "payload"
)

func authMiddleware(maker pkg.JWTMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if authHeader == "" {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(pkg.Errorf(pkg.AUTHENTICATION_ERROR, "No header was passed")),
			)

			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(
					pkg.Errorf(pkg.AUTHENTICATION_ERROR, "Invalid or Missing Bearer Token"),
				),
			)

			return
		}

		authType := fields[0]
		if strings.ToLower(authType) != authorizationHeaderBearerType {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(
					pkg.Errorf(pkg.AUTHENTICATION_ERROR, "Authentication Type Not Supported"),
				),
			)

			return
		}

		token := fields[1]

		payload, err := maker.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(pkg.Errorf(pkg.AUTHENTICATION_ERROR, "Access Token Not Valid")),
			)

			return
		}

		ctx.Set(authorizationPayloadKey, payload)

		ctx.Next()
	}
}

func CORSmiddleware(frontendUrls []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")

		for _, allowedOrigin := range frontendUrls {
			if origin == allowedOrigin {
				ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
