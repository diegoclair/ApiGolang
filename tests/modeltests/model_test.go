package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/diegoclair/ApiGolang/api/controllers"
	"github.com/diegoclair/ApiGolang/api/models"
)

var server = controllers.Server{}
var userInstance = models.User{}
var buyInstance = models.Buy{}
var saleInstance = models.Sale{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	refreshUserTable()

	user := models.User{
		FullName:  "Pet",
		Email:     "pet@gmail.com",
		Password:  "password",
		BirthDate: "1994/01/27",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() error {

	users := []models.User{
		models.User{
			FullName:  "Steven victor",
			Email:     "steven@gmail.com",
			Password:  "password",
			BirthDate: "1989/05/20",
		},
		models.User{
			FullName:  "Kenny Morris",
			Email:     "kenny@gmail.com",
			Password:  "password",
			BirthDate: "1977/12/17",
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func refreshUserBuyAndSaleTable() error {

	err := server.DB.DropTableIfExists(&models.User{}, &models.Buy{}, &models.Sale{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Buy{}, &models.Sale{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserOneBuyAndOneSale() (models.Buy,models.Sale, error) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		return models.Buy{},models.Sale{}, err
	}
	user := models.User{
		FullName:  "Sam Phil",
		Email:     "sam@gmail.com",
		Password:  "password",
		BirthDate: "1991/07/03",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Buy{},models.Sale{}, err
	}
	buy := models.Buy{
		BitcoinAmount:   "0.000247",
		AuthorID: user.ID,
	}
	sale := models.Sale{
		BitcoinAmount:   "0.000247",
		AuthorID: user.ID,
	}
	errBuy := server.DB.Model(&models.Buy{}).Create(&buy).Error
	errSale := server.DB.Model(&models.Sale{}).Create(&sale).Error
	if errBuy != nil {
		return models.Buy{},models.Sale{}, errBuy
	}
	if errSale != nil {
		return models.Buy{},models.Sale{}, errSale
	}
	return buy,sale, nil
}

func seedUsersBuysAndSales() ([]models.User, []models.Buy, []models.Sale, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Buy{}, []models.Sale{}, err
	}
	var users = []models.User{
		models.User{
			FullName: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
			BirthDate: 	"1990/02/28",
		},
		models.User{
			FullName: "Magu Frank",
			Email:    "magu@gmail.com",
			Password: "password",
			BirthDate: 	"1994/03/09",
		},
	}
	var buys = []models.Buy{
		models.Buy{
			BitcoinAmount: "0.59",
		},
		models.Buy{
			BitcoinAmount: "2.025",
		},
	}
	var sales = []models.Sale{
		models.Sale{
			BitcoinAmount: "0.005447",
		},
		models.Sale{
			BitcoinAmount: "0.125",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		buys[i].AuthorID = users[i].ID
		sales[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Buy{}).Create(&buys[i]).Error
		if err != nil {
			log.Fatalf("cannot seed buys table: %v", err)
		}
		err = server.DB.Model(&models.Sale{}).Create(&sales[i]).Error
		if err != nil {
			log.Fatalf("cannot seed sales table: %v", err)
		}
	}
	return users, buys,sales, nil
}