package models

import (
	"errors"
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/diegoclair/ApiGolang/api/coinmarketcap"
	"github.com/jinzhu/gorm"
)

type Sale struct {
	ID            uint64    `gorm:"primary_key;auto_increment" json:"id"`
	BitcoinAmount string    `gorm:"size:255;not null;" json:"bitcoin_amount"`
	Author        User      `json:"author"`
	AuthorID      uint32    `gorm:"not null" json:"author_id"`
	BitcoinPrice  float64   `gorm:"not null;" json:"bitcoin_price"`
	TotalBitcoin  float64   `gorm:"not null;" json:"total_bitcoin"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (s *Sale) Prepare(db *gorm.DB) {
	lastHr, lastMin := getLastHour(db)

	hr, min, price, newHour := coinmarketcap.GetBitcoinPrice(lastHr, lastMin)

	f, _ := strconv.ParseFloat(s.BitcoinAmount, 64)

	s.ID = 0
	s.BitcoinAmount = html.EscapeString(strings.TrimSpace(s.BitcoinAmount))
	s.BitcoinPrice = price
	s.TotalBitcoin = price * f
	s.Author = User{}
	s.CreatedAt = time.Now()

	fmt.Println(newHour, hr, min)
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

func (s *Sale) Validate() error {

	if s.BitcoinAmount == "" {
		return errors.New("Required BitcoinAmount")
	}
	if s.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (s *Sale) SaveSale(db *gorm.DB) (*Sale, error) {
	var err error
	err = db.Debug().Model(&Sale{}).Create(&s).Error
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
