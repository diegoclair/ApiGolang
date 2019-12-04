package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/diegoclair/ApiGolang/api/models"
)

var users = []models.User{
	models.User{
		FullName: "Steven victor",
		Email:    "steven@gmail.com",
		Password: 	"password",
		BirthDate: 	"1991/07/03",
	},
	models.User{
		FullName: 	"Martin Luther",
		Email:    	"luther@gmail.com",
		Password: 	"123456",
		BirthDate: 	"1993/10/03",
	},
}

var buys = []models.Buy{
	models.Buy{
		BitcoinAmount: "0.125",
	},
	models.Buy{
		BitcoinAmount: "2.025",
	},
}

var sales = []models.Sale{
	models.Sale{
		BitcoinAmount: "0.59",
	},
	models.Sale{
		BitcoinAmount: "0.005447",
	},
}

func Load(db *gorm.DB) {

	//Buy----------------------------------------------------------------------------------
	errBuy := db.Debug().DropTableIfExists(&models.Buy{}, &models.User{}).Error
	if errBuy != nil {
		log.Fatalf("cannot drop table: %v", errBuy)
	}
	errBuy = db.Debug().AutoMigrate(&models.User{}, &models.Buy{}).Error
	if errBuy != nil {
		log.Fatalf("cannot migrate table: %v", errBuy)
	}

	errBuy = db.Debug().Model(&models.Buy{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if errBuy != nil {
		log.Fatalf("attaching foreign key error: %v", errBuy)
	}

	for i, _ := range users {
		errBuy = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if errBuy != nil {
			log.Fatalf("cannot seed users table: %v", errBuy)
		}
		buys[i].AuthorID = users[i].ID

		errBuy = db.Debug().Model(&models.Buy{}).Create(&buys[i]).Error
		if errBuy != nil {
			log.Fatalf("cannot seed buys table: %v", errBuy)
		}
	}

	//Sale----------------------------------------------------------------------------------
	errSale := db.Debug().DropTableIfExists(&models.Sale{}, &models.User{}).Error
	if errSale != nil {
		log.Fatalf("cannot drop table: %v", errSale)
	}
	errSale = db.Debug().AutoMigrate(&models.User{}, &models.Sale{}).Error
	if errSale != nil {
		log.Fatalf("cannot migrate table: %v", errSale)
	}

	errSale = db.Debug().Model(&models.Sale{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if errSale != nil {
		log.Fatalf("attaching foreign key error: %v", errSale)
	}

	for i, _ := range users {
		errSale = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if errSale != nil {
			log.Fatalf("cannot seed users table: %v", errSale)
		}
		sales[i].AuthorID = users[i].ID

		errSale = db.Debug().Model(&models.Sale{}).Create(&sales[i]).Error
		if errSale != nil {
			log.Fatalf("cannot seed sales table: %v", errSale)
		}
	}
}
