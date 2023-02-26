package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	v1 "github.com/SaidovZohid/competition-project/api/v1"
	"github.com/SaidovZohid/competition-project/config"
	logging "github.com/SaidovZohid/competition-project/pkg/logger"
	"github.com/SaidovZohid/competition-project/pkg/token"
	"github.com/SaidovZohid/competition-project/storage"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/SaidovZohid/competition-project/api/docs" // for swagwger
)

type RouterOptions struct {
	Cfg        *config.Config
	Storage    storage.StorageI
	InMemory   storage.InMemoryStorageI
	TokenMaker token.Maker
	Logger     *logging.Logger
}

// @Security ApiKeyAuth
// New @title           Swagger doc for test taking website api
// @version         1.0
// @description     This is a api Swagger Doc.
// @BasePath  /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(opt *RouterOptions) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "*")
	router.Use(cors.New(corsConfig))

	handlerV1 := v1.New(&v1.HandlerV1Options{
		Cfg:        opt.Cfg,
		Storage:    opt.Storage,
		InMemory:   opt.InMemory,
		TokenMaker: opt.TokenMaker,
		Logger:     opt.Logger,
	})

	router.Static("/qr-codes", "./qr-codes")

	apiV1 := router.Group("/v1")
	apiV1.POST("/urls/make-short-url", handlerV1.AuthMiddleware, handlerV1.MakeShortUrl)
	apiV1.GET("/urls/:id", handlerV1.AuthMiddleware, handlerV1.GetUrl)
	apiV1.PUT("/urls/:id", handlerV1.AuthMiddleware, handlerV1.UpdateUrl)
	apiV1.DELETE("/urls/:id", handlerV1.AuthMiddleware, handlerV1.DeleteUrl)

	apiV1.GET("/users/:id", handlerV1.GetUser)
	apiV1.GET("/users", handlerV1.GetAllUsers)
	apiV1.PUT("/users/:id", handlerV1.AuthMiddleware, handlerV1.UpdateUser)
	apiV1.DELETE("/users/:id", handlerV1.AuthMiddleware, handlerV1.DeleteUser)
	apiV1.GET("/users/email/:email", handlerV1.GetUserByEmail)

	apiV1.POST("/auth/register", handlerV1.Register)
	apiV1.POST("/auth/verify", handlerV1.Verify)
	apiV1.POST("/auth/login", handlerV1.Login)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
