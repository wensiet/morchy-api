package entity

import "github.com/google/uuid"

type Container struct {
	ID     uuid.UUID `json:"id"`
	NodeID uuid.UUID `json:"node_id"`
	Image  string    `json:"image"`
}
