package container

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/wensiet/morchy-api/pkg/entity"
)

const (
	GetContainerQuery    = "SELECT id, node_id, image, status FROM container WHERE id = $1"
	ListContainersQuery  = "SELECT id, node_id, image, status FROM container"
	AddContainerQuery    = "INSERT INTO container (id, node_id, image, status) VALUES ($1, $2, $3, $4)"
	UpdateContainerQuery = "UPDATE container SET image = $1, status = $2 WHERE id = $3"
	DeleteContainerQuery = "DELETE FROM container WHERE id = $1"
)

type IService interface {
	GetContainer(ctx context.Context, id uuid.UUID) (*entity.Container, error)
	ListContainers(ctx context.Context) ([]entity.Container, error)
	AddContainer(ctx context.Context, nodeID uuid.UUID, image string) (*entity.Container, error)
	UpdateContainer(ctx context.Context, container *entity.Container) error
	RemoveContainer(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	dbPool *pgxpool.Pool
}

func NewService(dbPool *pgxpool.Pool) *Service {
	return &Service{dbPool: dbPool}
}

func (s *Service) GetContainer(ctx context.Context, id uuid.UUID) (*entity.Container, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var container entity.Container

	err = conn.QueryRow(ctx, GetContainerQuery, id).Scan(&container.ID, &container.NodeID, &container.Image, &container.Status)
	if err != nil {
		return nil, err
	}

	return &container, nil
}

func (s *Service) ListContainers(ctx context.Context) ([]entity.Container, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var containers []entity.Container
	rows, err := conn.Query(ctx, ListContainersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var container entity.Container
		err = rows.Scan(&container.ID, &container.NodeID, &container.Image, &container.Status)
		if err != nil {
			return nil, err
		}
		containers = append(containers, container)
	}

	return containers, nil
}

func (s *Service) AddContainer(ctx context.Context, nodeID uuid.UUID, image string) (*entity.Container, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	container := entity.NewContainer(nodeID, image)
	_, err = conn.Exec(ctx, AddContainerQuery, container.ID, container.NodeID, container.Image, container.Status)
	if err != nil {
		return nil, err
	}
	return container, nil
}

func (s *Service) RemoveContainer(ctx context.Context, id uuid.UUID) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, DeleteContainerQuery, id)
	return err
}

func (s *Service) UpdateContainer(ctx context.Context, container *entity.Container) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, UpdateContainerQuery, container.ID, container.Image, container.Status)
	return err
}
