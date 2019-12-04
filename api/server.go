package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/diegoclair/ApiGolang/api/controllers"
	"github.com/diegoclair/ApiGolang/api/seed"
)

var server = controllers.Server{}

func Run() {

	var errBuy error
	errBuy = godotenv.LoadBuy()
	if errBuy != nil {
		log.Fatalf("Error getting env, not comming through %v", errBuy)
	} else {
		fmt.Println("We are getting the env values")
	}

	var errSale error
	errSale = godotenv.LoadSale()
	if errSale != nil {
		log.Fatalf("Error getting env, not comming through %v", errSale)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	seed.Load(server.DB)

	server.Run(":8080")

}