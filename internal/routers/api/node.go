package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wensiet/morchy-api/internal/usecase/node"
	"github.com/wensiet/morchy-api/pkg/entity"
)

type INodeRouter interface {
	GetNode(c *gin.Context)
	ListNodes(c *gin.Context)
	AddNode(c *gin.Context)
	UpdateNode(c *gin.Context)
	DeleteNode(c *gin.Context)
}

type NodeRouter struct {
	nodeService node.IService
}

func NewNodeRouter(nodeService node.IService) NodeRouter {
	return NodeRouter{
		nodeService: nodeService,
	}
}

// GetNode godoc
//
//	@Summary		Get node by id
//	@Description	Allows to get node with it ID
//	@Tags			Node
//	@Accept			json
//	@Param			resource_id	path	string	true	"Node's ID"
//	@Produce		json
//	@Success		200
//	@Router			/api/v1/node/{resource_id} [get]
func (nr *NodeRouter) GetNode(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("resource_id"))
	if err != nil {
		c.JSON(422, gin.H{
			"message": err.Error(),
		})
		return
	}

	nodeModel, err := nr.nodeService.GetNode(ctx, id)
	if err != nil {
		c.JSON(500, gin.H{
			"kind":   "ERROR",
			"reason": err.Error(),
		})
		return
	}

	c.JSON(200, nodeModel)
}

// ListNodes godoc
//
//	@Summary		List all nodes
//	@Description	Retrieves a list of all nodes
//	@Tags			Node
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		entity.Node
//	@Failure		500	{object}	map[string]string
//	@Router			/api/v1/node [get]
func (nr *NodeRouter) ListNodes(c *gin.Context) {
	ctx := c.Request.Context()

	nodes, err := nr.nodeService.ListNodes(ctx)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, nodes)
}

// AddNode godoc
//
//	@Summary		Add a new node
//	@Description	Creates a new node
//	@Tags			Node
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Node
//	@Failure		500	{object}	map[string]string
//	@Router			/api/v1/node [post]
func (nr *NodeRouter) AddNode(c *gin.Context) {
	ctx := c.Request.Context()
	nodeModel, err := nr.nodeService.AddNode(ctx)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, nodeModel)
}

// UpdateNode godoc
//
//	@Summary		Update a node
//	@Description	Updates an existing node
//	@Tags			Node
//	@Accept			json
//	@Produce		json
//	@Param			node	body	entity.Node	true	"Updated Node Data"
//	@Success		204
//	@Failure		422	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/v1/node [put]
func (nr *NodeRouter) UpdateNode(c *gin.Context) {
	ctx := c.Request.Context()
	var nodeModel entity.Node

	err := c.BindJSON(&nodeModel)
	if err != nil {
		c.JSON(422, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = nr.nodeService.UpdateNode(ctx, &nodeModel)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(204, gin.H{})
}

// DeleteNode godoc
//
//	@Summary		Delete a node
//	@Description	Deletes a node by its ID
//	@Tags			Node
//	@Accept			json
//	@Produce		json
//	@Param			resource_id	path	string	true	"Node's ID"
//	@Success		204
//	@Failure		422	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/v1/node/{resource_id} [delete]
func (nr *NodeRouter) DeleteNode(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("resource_id"))
	if err != nil {
		c.JSON(422, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = nr.nodeService.DeleteNode(ctx, id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(204, gin.H{})
}
