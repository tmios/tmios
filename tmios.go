package main

import (
	"tmios/internal/cmp"
	"tmios/internal/config"
	"tmios/internal/http"
	"tmios/pkg/api"
)

func main() {
	err := cmp.NewCmp(
		config.NewConfig(
			config.WithConf(config.DefaultConfigFile, true),
			config.WithMysql(),
		),
		http.NewHttp(api.WithTest()),
	).Run()
	if err != nil {
		return
	}
}
