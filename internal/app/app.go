package app

import (
	"context"
	_ "github.com/wensiet/morchy-api/docs"
	"github.com/wensiet/morchy-api/internal/config"
	"github.com/wensiet/morchy-api/internal/infrastructure"
	"github.com/wensiet/morchy-api/internal/routers"
	"github.com/wensiet/morchy-api/internal/usecase/container"
	"github.com/wensiet/morchy-api/internal/usecase/node"
	"log"
)

const DefaultDBPoolMaxConn = 30

func Run() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	pgPool, err := infrastructure.NewPGPool(
		ctx,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		DefaultDBPoolMaxConn,
	)
	if err != nil {
		log.Fatal(err)
	}

	nodeService := node.NewService(pgPool)
	containerService := container.NewService(pgPool)

	router := routers.InitRouter(nodeService, containerService)

	err = router.Run()
	if err != nil {
		panic(err)
	}
}
