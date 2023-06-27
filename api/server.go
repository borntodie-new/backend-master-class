package api

import (
	"fmt"
	db "github.com/borntodie-new/backend-master-class/db/sqlc"
	"github.com/borntodie-new/backend-master-class/token"
	"github.com/borntodie-new/backend-master-class/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server HTTP request for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// 注册自定义的验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

// Start runs the HTTP server on a specific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func (s *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)
	router.POST("/token/renew_access", s.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))
	{
		authRoutes.POST("/accounts", s.createAccount)
		authRoutes.GET("/accounts/:id", s.getAccount)
		authRoutes.GET("/accounts", s.listAccount)
		authRoutes.POST("/transfers", s.createTransfer)
	}

	s.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
