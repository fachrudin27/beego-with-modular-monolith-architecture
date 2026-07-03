package main

import (
	"context"
	"fmt"

	"firstbeegoapi/internal/ordering/app"
	orderingapi "firstbeegoapi/internal/ordering/delivery/api"
	orderingpostgres "firstbeegoapi/internal/ordering/infra/postgres"
	orderingrepo "firstbeegoapi/internal/ordering/infra/postgres/repository"
	"firstbeegoapi/internal/shared"
	_ "firstbeegoapi/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	zapConfig := `{"development":false,"encoding":"json","level":"info"}`

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		zapConfig = `{"development":true,"encoding":"console","level":"debug"}`
	}

	if err := shared.InitZapLogger(zapConfig); err != nil {
		panic(fmt.Errorf("failed to configure zap logger: %w", err))
	}
	defer shared.SyncZapLogger()

	if err := shared.ValidateJWTConfig(beego.BConfig.RunMode); err != nil {
		panic(fmt.Errorf("invalid jwt configuration: %w", err))
	}

	if err := orderingpostgres.Init(context.Background()); err != nil {
		panic(fmt.Errorf("failed to initialize ordering postgres: %w", err))
	}
	defer func() {
		_ = orderingpostgres.Close()
	}()

	orderingDB, err := orderingpostgres.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get ordering postgres: %w", err))
	}

	orderingRepository := orderingrepo.NewOrderingRepository(orderingDB)
	orderingapi.SetOrderingService(app.NewOrderingService(orderingRepository))

	beego.Run()
}
