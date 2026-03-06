package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rlapenok/rybakov_test/internal/config"
	"github.com/rlapenok/rybakov_test/internal/infra/transport/handlers"

	"github.com/rs/zerolog"
)

// Server is a struct that contains the http server
type Server struct {
	server *http.Server
}

// NewServer creates a new server
func NewServer(
	cfg *config.ServerConfig,
	logger *zerolog.Logger,
	authToken string,
	withdrawalHandler *handlers.WithdrawalHandler,
) *Server {
	// Set the environment mode
	gin.SetMode(gin.ReleaseMode)

	// Create a new engine
	engine := gin.New()

	engine.Use(loggerMiddleware(logger))
	engine.Use(errorsMiddleware())

	apiGroup := engine.Group("/api")
	versionGroup := apiGroup.Group("/v1")
	versionGroup.Use(authMiddleware(authToken))
	versionGroup.POST("/withdrawals", withdrawalHandler.CreateWithdrawal)
	versionGroup.GET("/withdrawals/:id", withdrawalHandler.GetWithdrawalByID)
	// Create a new server
	return &Server{server: &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: engine}}
}

// Start starts the server
func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return ErrStartServer.WithMessage(err.Error())
	}

	return nil
}

// Stop stops the server
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

/*
func instanceRoutes(group *gin.RouterGroup, greenClient instanceDomain.InstanceAPIClient) {
	getSettingsUseCase := instanceUseCase.NewGetSettingsUseCase(greenClient)
	getStateUseCase := instanceUseCase.NewGetStateUseCase(greenClient)

	handler := instanceTransport.NewHandler(getSettingsUseCase, getStateUseCase)

	instanceGroup := group.Group("/instance")
	instanceGroup.GET("/settings", handler.GetSettings)
	instanceGroup.GET("/state", handler.GetState)
}

func sendingRoutes(group *gin.RouterGroup, greenClient sendingDomain.SendingAPIClient) {
	sendMessageUseCase := sendingUseCase.NewSendMessageUseCase(greenClient)
	sendFileByURLUseCase := sendingUseCase.NewSendFileByURLUseCase(greenClient)
	handler := sendingTransport.NewHandler(sendMessageUseCase, sendFileByURLUseCase)

	sendingGroup := group.Group("/sending")
	sendingGroup.POST("/message", handler.SendMessage)
	sendingGroup.POST("/fileByUrl", handler.SendFileByURL)
}
*/
