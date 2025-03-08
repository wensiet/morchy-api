package entity

import (
	"errors"
	"github.com/google/uuid"
)

var InvalidNodeStatusErr = errors.New("invalid node status")

type NodeStatus string

type Node struct {
	ID         uuid.UUID   `json:"id"`
	Status     NodeStatus  `json:"status"`
	Containers []Container `json:"containers"`
}

func NewNode() *Node {
	return &Node{
		ID:         uuid.New(),
		Status:     NewNodeStatus,
		Containers: []Container{},
	}
}

const (
	NewNodeStatus     NodeStatus = "new"
	RunningNodeStatus NodeStatus = "running"
	FailedNodeStatus  NodeStatus = "failed"
)

func (ns NodeStatus) Validate() error {
	switch ns {
	case NewNodeStatus, RunningNodeStatus, FailedNodeStatus:
		return nil
	default:
		return InvalidNodeStatusErr
	}
}
