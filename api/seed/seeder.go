package seed

import (
	"log"

	"github.com/diegoclair/ApiGolang/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		FullName:  "Steven victor",
		Email:     "steven@gmail.com",
		Password:  "password",
		BirthDate: "1991/07/03",
	},
	models.User{
		FullName:  "Martin Luther",
		Email:     "luther@gmail.com",
		Password:  "123456",
		BirthDate: "1993/10/03",
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

var hour = []models.LastHour{
	models.LastHour{
		Hora:   0,
		Minuto: 0,
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Buy{}, &models.Sale{}, &models.User{}, &models.LastHour{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err.Error())
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Buy{}, &models.Sale{}, &models.LastHour{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Buy{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	err = db.Debug().Model(&models.Sale{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.LastHour{}).Create(&hour).Error
	if err != nil {
		log.Fatalf("cannot seed last hour table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		buys[i].AuthorID = users[i].ID
		sales[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Buy{}).Create(&buys[i]).Error
		if err != nil {
			log.Fatalf("cannot seed buys table: %v", err)
		}

		err = db.Debug().Model(&models.Sale{}).Create(&sales[i]).Error
		if err != nil {
			log.Fatalf("cannot seed sales table: %v", err)
		}
	}
}
