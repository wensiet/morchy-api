package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wensiet/morchy-api/internal/routers/api"
	"github.com/wensiet/morchy-api/internal/usecase/container"
	"github.com/wensiet/morchy-api/internal/usecase/node"
)

func InitRouter(nodeService node.IService, containerService container.IService) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	nodeRoutes := api.NewNodeRouter(
		nodeService,
	)
	containerRoutes := api.NewContainerRouter(
		containerService,
	)

	apiv1 := r.Group("/api/v1")
	{
		nodeRouter := apiv1.Group("/node")
		{
			nodeRouter.GET("/:resource_id", nodeRoutes.GetNode)
			nodeRouter.GET("", nodeRoutes.ListNodes)
			nodeRouter.POST("", nodeRoutes.AddNode)
			nodeRouter.PUT("", nodeRoutes.UpdateNode)
			nodeRouter.DELETE("/:resource_id", nodeRoutes.DeleteNode)
		}
		containerRouter := apiv1.Group("/container")
		{
			containerRouter.GET("/:resource_id", containerRoutes.GetContainer)
			containerRouter.GET("", containerRoutes.ListContainers)
			containerRouter.POST("", containerRoutes.AddContainer)
			containerRouter.PUT("", containerRoutes.UpdateContainer)
			containerRouter.DELETE("/:resource_id", containerRoutes.DeleteContainer)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
