package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/inkclip/backend/config"
	db "github.com/inkclip/backend/db/sqlc"
	docs "github.com/inkclip/backend/docs"
	"github.com/inkclip/backend/mail"
	"github.com/inkclip/backend/token"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type Server struct {
	config     config.Config
	store      db.Store
	tokenMaker token.Maker
	mailClient mail.Client
	router     *gin.Engine
}

func NewServer(config config.Config, store db.Store, mailClient mail.Client) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		mailClient: mailClient,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	// TODO: serverに持たせる
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	router := gin.Default()

	router.Use(jsonMiddleware()).Use(corsMiddleware()).Use(loggerMiddleware(logger))

	docs.SwaggerInfo.BasePath = "/"

	router.POST("/register", server.register)
	router.POST("/verify", server.verify)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/users/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/me", server.getMe)
	authRoutes.GET("/users/:id", server.getUser)
	authRoutes.GET("/users", server.listUser)

	authRoutes.POST("/webs", server.createWeb)
	authRoutes.GET("/webs/:id", server.getWeb)
	authRoutes.GET("/webs", server.listWeb)
	authRoutes.DELETE("/webs/:id", server.deleteWeb)

	authRoutes.POST("/notes", server.createNote)
	authRoutes.GET("/notes/:id", server.getNote)
	authRoutes.GET("/notes", server.listNote)
	authRoutes.DELETE("/notes/:id", server.deleteNote)
	authRoutes.PUT("/notes/:id", server.putNote)
	router.GET("/public_notes/:id", server.getPublicNote)

	// TODO: only env is dev
	if server.config.Env == "dev" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
