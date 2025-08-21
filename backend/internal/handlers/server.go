package handlers

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/flexGURU/flower-haven/backend/internal/postgres"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	ln     net.Listener
	srv    *http.Server

	config     pkg.Config
	tokenMaker pkg.JWTMaker
	repo       *postgres.PostgresRepo
}

func NewServer(config pkg.Config, tokenMaker pkg.JWTMaker, repo *postgres.PostgresRepo) *Server {
	if config.ENVIRONMENT == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	s := &Server{
		router: r,
		ln:     nil,

		config:     config,
		tokenMaker: tokenMaker,
		repo:       repo,
	}

	s.setUpRoutes()

	return s
}

func (s *Server) setUpRoutes() {
	s.router.Use(CORSmiddleware(s.config.FRONTEND_URL))
	v1 := s.router.Group("/api/v1")

	v1Auth := s.router.Group("/api/v1")
	authRoute := v1Auth.Use(authMiddleware(s.tokenMaker))

	// health check
	s.router.GET("/health-check", s.healthCheckHandler)

	// User routes
	v1.POST("/user/login", s.login)
	v1.GET("/user/logout", s.logout)
	v1.POST("/user/refresh-token", s.refreshToken)

	v1.POST("/users", s.createUserHandler)
	authRoute.GET("/users/:id", s.getUserHandler)
	authRoute.GET("/users", s.listUsersHandler)
	authRoute.PUT("/users/:id", s.updateUserHandler)

	// v1.GET("/users/:id/orders", s.getUserOrdersHandler)
	authRoute.GET("/users/:id/subscriptions", s.getUserSubscriptionsHandler)

	// Category routes
	authRoute.POST("/categories", s.createCategoryHandler)
	v1.GET("/categories/:id", s.getCategoryHandler)
	v1.GET("/categories", s.listCategoriesHandler)
	authRoute.PUT("/categories/:id", s.updateCategoryHandler)
	authRoute.DELETE("/categories/:id", s.deleteCategoryHandler)

	// Product routes
	authRoute.POST("/products", s.createProductHandler)
	v1.GET("/products/:id", s.getProductHandler)
	v1.GET("/products", s.listProductsHandler)
	authRoute.PUT("/products/:id", s.updateProductHandler)
	authRoute.DELETE("/products/:id", s.deleteProductHandler)

	authRoute.GET("/products/:id/order-items", s.listProductOrderItemsHandler)

	// Subscription routes
	authRoute.POST("/subscriptions", s.createSubscriptionHandler)
	v1.GET("/subscriptions/:id", s.getSubscriptionHandler)
	v1.GET("/subscriptions", s.listSubscriptionsHandler)
	authRoute.PUT("/subscriptions/:id", s.updateSubscriptionHandler)
	authRoute.DELETE("/subscriptions/:id", s.deleteSubscriptionHandler)

	// User Subscription routes
	authRoute.POST("/user-subscriptions", s.createUserSubscriptionHandler)
	authRoute.GET("/user-subscriptions/:id", s.getUserSubscriptionHandler)
	authRoute.GET("/user-subscriptions", s.listUserSubscriptionsHandler)
	authRoute.PUT("/user-subscriptions/:id", s.updateUserSubscriptionHandler)
	authRoute.DELETE("/user-subscriptions/:id", s.deleteUserSubscriptionHandler)

	// Subscription Deliveries routes
	authRoute.POST("/subscription-deliveries", s.createSubscriptionDeliveryHandler)
	authRoute.GET("/subscription-deliveries/:id", s.getSubscriptionDeliveryByUserSubscriptionIDHandler)
	authRoute.GET("/subscription-deliveries", s.listSubscriptionDeliveriesHandler)
	authRoute.PUT("/subscription-deliveries/:id", s.updateSubscriptionDeliveryHandler)
	authRoute.DELETE("/subscription-deliveries/:id", s.deleteSubscriptionDeliveryHandler)

	// Order routes
	authRoute.POST("/orders", s.createOrderHandler)
	authRoute.GET("/orders/:id", s.getOrderHandler)
	authRoute.GET("/orders", s.listOrdersHandler)
	authRoute.PUT("/orders/:id", s.updateOrderHandler)
	authRoute.DELETE("/orders/:id", s.deleteOrderHandler)

	// Payment routes
	authRoute.POST("/payments", s.createPaymentHandler)
	authRoute.GET("/payments/:id", s.getPaymentHandler)
	authRoute.PUT("/payments/:id", s.updatePaymentHandler)
	authRoute.GET("/payments", s.listPaymentsHandler)

	// helpers routes
	authRoute.GET("/dashboard", s.getDashboardDataHandler)

	s.srv = &http.Server{
		Addr:         s.config.SERVER_ADDRESS,
		Handler:      s.router.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) Start() error {
	var err error
	if s.ln, err = net.Listen("tcp", s.config.SERVER_ADDRESS); err != nil {
		return err
	}

	go func(s *Server) {
		err := s.srv.Serve(s.ln)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}(s)

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Shutting down http server...")

	return s.srv.Shutdown(ctx)
}

func (s *Server) GetPort() int {
	if s.ln == nil {
		return 0
	}

	return s.ln.Addr().(*net.TCPAddr).Port
}

func errorResponse(err error) gin.H {
	return gin.H{
		"status_code": pkg.ErrorCode(err),
		"message":     pkg.ErrorMessage(err),
	}
}
