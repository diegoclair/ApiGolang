package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/diegoclair/ApiGolang/api/coinmarketcap"
	"github.com/jinzhu/gorm"
)

type Sale struct {
	ID                uint64    `gorm:"primary_key;auto_increment" json:"id"`
	BitcoinAmount     float64   `gorm:"not null;" json:"bitcoin_amount"`
	Author            User      `json:"author"`
	AuthorID          uint32    `gorm:"not null" json:"author_id"`
	BitcoinPrice      float64   `gorm:"not null;" json:"bitcoin_price"`
	TotalBitcoinPrice float64   `gorm:"not null;" json:"total_bitcoin_price"`
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (s *Sale) Prepare(db *gorm.DB) {
	lastHr, lastMin := GetLastHour(db)

	hr, min, price, newHour := coinmarketcap.GetBitcoinPrice(lastHr, lastMin)

	f := s.BitcoinAmount

	s.ID = 0
	s.BitcoinAmount = s.BitcoinAmount
	s.BitcoinPrice = price
	s.TotalBitcoinPrice = price * f
	s.Author = User{}
	s.CreatedAt = time.Now()

	if newHour {
		var hours = []LastHour{
			LastHour{
				Hour:   hr,
				Minute: min,
			},
		}
		hours[0].UpdateLastHour(db)
	}
}

func (s *Sale) Validate(db *gorm.DB) error {

	if s.BitcoinAmount == 0 {
		return errors.New("Required BitcoinAmount")
	}
	if s.AuthorID < 1 {
		return errors.New("Required Author")
	}

	var err error
	var user = User{}

	err = db.Debug().Model(&User{}).Where("id = ?", s.AuthorID).Take(&user).Error
	if err != nil {
		return errors.New("Error to get user, func Validade, file SaleTransactions")
	}
	if user.TotalBitcoinAmount < s.BitcoinAmount {
		var str = "Insufficient Biticoin amount. Your current amount is: " + fmt.Sprintf("%f", user.TotalBitcoinAmount)
		return errors.New(str)
	}

	return nil
}

func (s *Sale) SaveSale(db *gorm.DB) (*Sale, error) {
	var err error
	//to create sale
	err = db.Debug().Model(&Sale{}).Create(&s).Error
	if err != nil {
		return &Sale{}, err
	}

	//to update user bitcoinAmount and balance
	var user = User{}
	err = db.Debug().Model(&User{}).Where("id = ?", s.AuthorID).Take(&user).Error
	if err != nil {
		return s, errors.New("Error to get user, func SaveSale, file SaleTransactions")
	}

	db = db.Debug().Model(&User{}).Where("id = ?", s.AuthorID).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"balance":              user.Balance + s.TotalBitcoinPrice,
			"total_bitcoin_amount": user.TotalBitcoinAmount - s.BitcoinAmount,
		},
	)
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return s, errors.New("User not found")
		}
		return s, db.Error
	}

	//to get user data to return in the sale
	if s.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", s.AuthorID).Take(&s.Author).Error
		if err != nil {
			return &Sale{}, err
		}
	}

	return s, nil
}

func (s *Sale) FindAllSales(db *gorm.DB) (*[]Sale, error) {
	var err error
	sales := []Sale{}
	err = db.Debug().Model(&Sale{}).Limit(100).Find(&sales).Error
	if err != nil {
		return &[]Sale{}, err
	}
	if len(sales) > 0 {
		for i := range sales {
			err := db.Debug().Model(&User{}).Where("id = ?", sales[i].AuthorID).Take(&sales[i].Author).Error
			if err != nil {
				return &[]Sale{}, err
			}
		}
	}
	return &sales, nil
}

func (s *Sale) FindSaleByID(db *gorm.DB, pid uint64) (*Sale, error) {
	var err error
	err = db.Debug().Model(&Sale{}).Where("id = ?", pid).Take(&s).Error
	if err != nil {
		return &Sale{}, err
	}
	if s.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", s.AuthorID).Take(&s.Author).Error
		if err != nil {
			return &Sale{}, err
		}
	}
	return s, nil
}

func (s *Sale) UpdateASale(db *gorm.DB) (*Sale, error) {

	var err error

	err = db.Debug().Model(&Sale{}).Where("id = ?", s.ID).Updates(Sale{BitcoinAmount: s.BitcoinAmount}).Error
	if err != nil {
		return &Sale{}, err
	}
	if s.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", s.AuthorID).Take(&s.Author).Error
		if err != nil {
			return &Sale{}, err
		}
	}
	return s, nil
}

func (s *Sale) DeleteASale(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Sale{}).Where("id = ? and author_id = ?", pid, uid).Take(&Sale{}).Delete(&Sale{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Sale not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
