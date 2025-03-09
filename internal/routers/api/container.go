package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wensiet/morchy-api/internal/usecase/container"
	"github.com/wensiet/morchy-api/pkg/entity"
)

type IContainerRouter interface {
	GetContainer(c *gin.Context)
	ListContainers(c *gin.Context)
	AddContainer(c *gin.Context)
	UpdateContainer(c *gin.Context)
	DeleteContainer(c *gin.Context)
}

type ContainerRouter struct {
	containerService container.IService
}

func NewContainerRouter(containerService container.IService) ContainerRouter {
	return ContainerRouter{containerService: containerService}
}

// GetContainer godoc
//
//	@Summary		Get container by id
//	@Description	Allows to get a container by its ID
//	@Tags			Container
//	@Accept			json
//	@Param			resource_id	path	string	true	"Container's ID"
//	@Produce		json
//	@Success		200	{object}	entity.Container
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/v1/container/{resource_id} [get]
func (cr *ContainerRouter) GetContainer(c *gin.Context) {
	idParam := c.Param("resource_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	containerModel, err := cr.containerService.GetContainer(c, id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
			"kind":    "error",
		})
		return
	}
	c.JSON(200, containerModel)
}

// ListContainers godoc
//
//	@Summary		List all containers
//	@Description	Retrieves a list of all containers
//	@Tags			Container
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		entity.Container
//	@Failure		500	{object}	map[string]string
//	@Router			/api/v1/container [get]
func (cr *ContainerRouter) ListContainers(c *gin.Context) {
	containers, err := cr.containerService.ListContainers(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, containers)
}

// AddContainer godoc
//
//	@Summary		Add a new container
//	@Description	Creates a new container
//	@Tags			Container
//	@Accept			json
//	@Produce		json
//	@Param			container	body		entity.AddContainer	true	"New container data"
//	@Success		201			{object}	entity.Container
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/api/v1/container [post]
func (cr *ContainerRouter) AddContainer(c *gin.Context) {
	var req *entity.AddContainer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	containerModel, err := cr.containerService.AddContainer(c, req.NodeID, req.Image)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, containerModel)
}

// UpdateContainer godoc
//
//	@Summary		Update an existing container
//	@Description	Updates the details of an existing container
//	@Tags			Container
//	@Accept			json
//	@Produce		json
//	@Param			container	body	entity.Container	true	"Updated container data"
//	@Success		204
//	@Failure		422	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/v1/container [put]
func (cr *ContainerRouter) UpdateContainer(c *gin.Context) {
	var containerModel *entity.Container

	if err := c.ShouldBindJSON(&containerModel); err != nil {
		c.JSON(422, gin.H{"error": err.Error()})
		return
	}

	if err := cr.containerService.UpdateContainer(c, containerModel); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(204, gin.H{})
}

// DeleteContainer godoc
//
//	@Summary		Delete a container
//	@Description	Deletes a container by its ID
//	@Tags			Container
//	@Accept			json
//	@Produce		json
//	@Param			resource_id	path	string	true	"Container's ID"
//	@Success		204
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/v1/container/{resource_id} [delete]
func (cr *ContainerRouter) DeleteContainer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = cr.containerService.RemoveContainer(c, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(204, gin.H{})
}
