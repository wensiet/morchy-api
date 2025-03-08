package infrastructure

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPGPool(ctx context.Context, user, password, host, port, dbName string, maxConn int32) (*pgxpool.Pool, error) {
	conn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		user, password, host, port, dbName,
	)
	conf, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, err
	}

	conf.MaxConns = maxConn
	conf.ConnConfig.PreferSimpleProtocol = true

	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
