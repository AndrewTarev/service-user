package http

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"service-user/docs"
	"service-user/internal/app/delivery/middleware"
	"service-user/internal/app/service"
	"service-user/internal/app/utils"
)

type UserProfileHandler interface {
	CreateProfile(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	DeleteProfile(c *gin.Context)
}

type Handler struct {
	auth *utils.JWTManager
	UserProfileHandler
}

func NewHandler(services *service.Service, auth *utils.JWTManager) *Handler {
	return &Handler{
		auth:               auth,
		UserProfileHandler: NewProfileHandler(*services),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(middleware.ErrorHandler())

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := router.Group("/api/v1")
	{
		profile := apiV1.Group("/user-profile")
		profile.Use(middleware.AuthMiddleware(h.auth))
		{
			profile.POST("/", h.CreateProfile)
			profile.GET("/", h.GetProfile)
			profile.PATCH("/", h.UpdateProfile)
			profile.DELETE("/", h.DeleteProfile)
		}
	}

	return router
}
