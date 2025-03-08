package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/wensiet/morchy-api/internal/routers/api"
	"github.com/wensiet/morchy-api/internal/usecase/node"
)

func InitRouter(nodeService node.IService) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	nodeRoutes := api.NewNodeRouter(
		nodeService,
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
	}

	return r
}
