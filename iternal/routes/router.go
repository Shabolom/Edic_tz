package routes

import (
	"Edos_Docer/iternal/api"
	"Edos_Docer/iternal/middleware"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	user := api.NewUserApi()

	r := gin.New()
	r.Use(gin.Logger())

	authRequired := r.Group("/")
	authRequired.Use(middleware.Authorization())
	authRequired.Use(middleware.Logger())
	authRequired.Use(middleware.GZIP())

	r.POST("/api/register", user.Register)
	r.POST("/api/login", user.Login)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	{
		authRequired.DELETE("/api/delete_acc", user.SelfDelete)
		authRequired.POST("/api/post_csv", user.PostCSV)
		authRequired.POST("/api/swap_login_password", user.Swap)
		authRequired.GET("/api/get_size_of_elements/:value", user.GetBatch)
	}

	return r
}
