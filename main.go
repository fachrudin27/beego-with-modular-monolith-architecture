package main

import (
	"fmt"

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

	beego.Run()
}
