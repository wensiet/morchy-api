package entity

import (
	"errors"
	"github.com/google/uuid"
)

var InvalidContainerStatusErr = errors.New("invalid container status")

type ContainerStatus string

// AddContainer godoc
// entity.AddContainer struct
type AddContainer struct {
	NodeID uuid.UUID `json:"node_id"`
	Image  string    `json:"image"`
}

// Container godoc
// entity.Container struct
type Container struct {
	ID     uuid.UUID       `json:"id"`
	NodeID uuid.UUID       `json:"node_id"`
	Image  string          `json:"image"`
	Status ContainerStatus `json:"status"`
}

func NewContainer(nodeID uuid.UUID, image string) *Container {
	return &Container{
		ID:     uuid.New(),
		NodeID: nodeID,
		Image:  image,
		Status: ContainerStatusPending,
	}
}

const (
	ContainerStatusRunning ContainerStatus = "running"
	ContainerStatusFailed  ContainerStatus = "failed"
	ContainerStatusPending ContainerStatus = "pending"
)

func (cs ContainerStatus) Validate() error {
	switch cs {
	case ContainerStatusRunning, ContainerStatusPending, ContainerStatusFailed:
		return nil
	default:
		return InvalidContainerStatusErr
	}
}
