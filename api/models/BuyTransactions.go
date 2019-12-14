package models

import (
	"errors"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/diegoclair/ApiGolang/api/coinmarketcap"
	"github.com/jinzhu/gorm"
)

type Buy struct {
	ID            uint64    `gorm:"primary_key;auto_increment" json:"id"`
	BitcoinAmount string    `gorm:"size:255;not null;" json:"bitcoin_amount"`
	Author        User      `json:"author"`
	AuthorID      uint32    `gorm:"not null" json:"author_id"`
	BitcoinPrice  float64   `gorm:"not null;" json:"bitcoin_price"`
	TotalBitcoin  float64   `gorm:"not null;" json:"total_bitcoin"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (b *Buy) Prepare(db *gorm.DB) {

	lastHr, lastMin := getLastHour(db)

	hr, min, price, newHour := coinmarketcap.GetBitcoinPrice(lastHr, lastMin)

	f, _ := strconv.ParseFloat(b.BitcoinAmount, 64)

	b.ID = 0
	b.BitcoinAmount = html.EscapeString(strings.TrimSpace(b.BitcoinAmount))
	b.BitcoinPrice = price
	b.TotalBitcoin = price * f
	b.Author = User{}
	b.CreatedAt = time.Now()

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

//Validate Bitcoin Fields
func (b *Buy) Validate() error {

	if b.BitcoinAmount == "" {
		return errors.New("Required Bitcoin Amount")
	}
	if b.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (b *Buy) SaveBuy(db *gorm.DB) (*Buy, error) {
	var err error
	err = db.Debug().Model(&Buy{}).Create(&b).Error
	if err != nil {
		return &Buy{}, err
	}
	if b.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", b.AuthorID).Take(&b.Author).Error
		if err != nil {
			return &Buy{}, err
		}
	}
	return b, nil
}

func (b *Buy) FindAllBuys(db *gorm.DB) (*[]Buy, error) {
	var err error
	buys := []Buy{}
	err = db.Debug().Model(&Buy{}).Limit(100).Find(&buys).Error
	if err != nil {
		return &[]Buy{}, err
	}
	if len(buys) > 0 {
		for i := range buys {
			err := db.Debug().Model(&User{}).Where("id = ?", buys[i].AuthorID).Take(&buys[i].Author).Error
			if err != nil {
				return &[]Buy{}, err
			}
		}
	}
	return &buys, nil
}

func (b *Buy) FindBuyByID(db *gorm.DB, pid uint64) (*Buy, error) {
	var err error
	err = db.Debug().Model(&Buy{}).Where("id = ?", pid).Take(&b).Error
	if err != nil {
		return &Buy{}, err
	}
	if b.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", b.AuthorID).Take(&b.Author).Error
		if err != nil {
			return &Buy{}, err
		}
	}
	return b, nil
}

func (b *Buy) UpdateABuy(db *gorm.DB) (*Buy, error) {

	var err error

	err = db.Debug().Model(&Buy{}).Where("id = ?", b.ID).Updates(Buy{BitcoinAmount: b.BitcoinAmount}).Error
	if err != nil {
		return &Buy{}, err
	}
	if b.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", b.AuthorID).Take(&b.Author).Error
		if err != nil {
			return &Buy{}, err
		}
	}
	return b, nil
}

func (b *Buy) DeleteABuy(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Buy{}).Where("id = ? and author_id = ?", pid, uid).Take(&Buy{}).Delete(&Buy{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Buy not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
