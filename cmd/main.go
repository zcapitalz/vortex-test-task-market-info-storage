package main

import (
	"market-info-storage/internal/app"
	"market-info-storage/internal/config"
)

func main() {
	cfg := config.MustNew()
	app.Run(cfg)
}
