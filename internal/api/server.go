package api

import (
	"log"

	"github.com/gin-gonic/gin"
	ginSwaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router  *gin.Engine
	handler *Handler
}

func NewServer(handler *Handler) *Server {
	server := &Server{
		router:  gin.Default(),
		handler: handler,
	}

	server.setupRoutes()
	server.router.GET("/swagger/*any", ginSwagger.WrapHandler(ginSwaggerFiles.Handler))
	return server
}

func (s *Server) setupRoutes() {
	s.router.POST("/accounts", s.handler.CreateAccount)
	s.router.GET("/accounts/:account_id", s.handler.GetAccount)
	s.router.POST("/transactions", s.handler.CreateTransaction)
}

func (s *Server) Start(addr string) error {
	log.Printf("Server starting on %s", addr)
	return s.router.Run(addr)
}
