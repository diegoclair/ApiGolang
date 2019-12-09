package main

import (
	"github.com/diegoclair/ApiGolang/api"
	"github.com/diegoclair/ApiGolang/api/coinmarketcap"
)

func main() {
	api.Run()
	coinmarketcap.Run()
}
