package node

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/wensiet/morchy-api/internal/usecase"
	"github.com/wensiet/morchy-api/pkg/entity"
)

const (
	GetNodeQueryWithContainers = `
		SELECT n.id, n.status, c.id, c.node_id, c.image, c.status
		FROM node n
		LEFT JOIN container c ON n.id = c.node_id
		WHERE n.id = $1`
	ListNodesQueryWithContainers = `
		SELECT n.id, n.status, c.id, c.node_id, c.image, c.status
		FROM node n
		LEFT JOIN container c ON n.id = c.node_id`
	AddNodeQuery    = "INSERT INTO node(id, status) VALUES($1, $2)"
	UpdateNodeQuery = "UPDATE node SET status = $1 WHERE id = $2"
	DeleteNodeQuery = "DELETE FROM node WHERE id = $1"
)

type IService interface {
	GetNode(ctx context.Context, id uuid.UUID) (*entity.Node, error)
	ListNodes(ctx context.Context) ([]*entity.Node, error)
	AddNode(ctx context.Context) (*entity.Node, error)
	UpdateNode(ctx context.Context, node *entity.Node) error
	DeleteNode(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	dbPool *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{
		dbPool: db,
	}
}

func (s *Service) GetNode(ctx context.Context, id uuid.UUID) (*entity.Node, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var node entity.Node
	rows, err := conn.Query(ctx, GetNodeQueryWithContainers, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var container entity.Container
		var image, status sql.NullString

		err := rows.Scan(&node.ID, &node.Status, &container.ID, &container.NodeID, &image, &status)
		if err != nil {
			return nil, err
		}

		if container.ID != uuid.Nil && image.Valid && status.Valid {
			container.Image = image.String
			container.Status = entity.ContainerStatus(status.String)
			err := container.Status.Validate()
			if err != nil {
				return nil, err
			}
			node.Containers = append(node.Containers, container)
		}
	}

	if node.ID == uuid.Nil {
		return nil, usecase.NodeNotFoundErr
	}

	return &node, nil
}

func (s *Service) ListNodes(ctx context.Context) ([]*entity.Node, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var nodes []*entity.Node
	rows, err := conn.Query(ctx, ListNodesQueryWithContainers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var node entity.Node
		err := rows.Scan(&node.ID, &node.Status)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &node)
	}
	return nodes, nil
}

func (s *Service) AddNode(ctx context.Context) (*entity.Node, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	node := entity.NewNode()
	_, err = conn.Exec(ctx, AddNodeQuery, node.ID, node.Status)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (s *Service) UpdateNode(ctx context.Context, node *entity.Node) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	err = node.Status.Validate()
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, UpdateNodeQuery, node.Status, node.ID)
	return err
}

func (s *Service) DeleteNode(ctx context.Context, id uuid.UUID) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, DeleteNodeQuery, id)
	return err
}
