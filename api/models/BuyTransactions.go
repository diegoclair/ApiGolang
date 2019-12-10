package models

import (
	"errors"
	"fmt"
	"html"
	"reflect"
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

func getLastHour(db *gorm.DB) (int, int) {

	hour, _ := FindLastHour(db)

	last_time := reflect.ValueOf(hour[0])
	fmt.Println("last_time", last_time)
	l_hr := last_time.Field(0).Interface().(int)
	l_min := last_time.Field(1).Interface().(int)

	return l_hr, l_min
}

func (p *Buy) Prepare(db *gorm.DB) {

	l_hr, l_min := getLastHour(db)

	hr, min, price, newHour := coinmarketcap.GetBitcoinPrice(l_hr, l_min)

	f, _ := strconv.ParseFloat(p.BitcoinAmount, 64)

	p.ID = 0
	p.BitcoinAmount = html.EscapeString(strings.TrimSpace(p.BitcoinAmount))
	p.BitcoinPrice = price
	p.TotalBitcoin = price * f
	p.Author = User{}
	p.CreatedAt = time.Now()

	if newHour {
		var hours = []LastHour{
			LastHour{
				Hour:   hr,
				Minute: min,
			},
		}
		hours[0].UpdateLastHour(db)
		fmt.Println(newHour, hr, min)
	}
}

func (p *Buy) Validate() error {

	if p.BitcoinAmount == "" {
		return errors.New("Required Bitcoin Amount")
	}
	if p.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (p *Buy) SaveBuy(db *gorm.DB) (*Buy, error) {
	var err error
	err = db.Debug().Model(&Buy{}).Create(&p).Error
	if err != nil {
		return &Buy{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Buy{}, err
		}
	}
	return p, nil
}

func (p *Buy) FindAllBuys(db *gorm.DB) (*[]Buy, error) {
	var err error
	buys := []Buy{}
	err = db.Debug().Model(&Buy{}).Limit(100).Find(&buys).Error
	if err != nil {
		return &[]Buy{}, err
	}
	if len(buys) > 0 {
		for i, _ := range buys {
			err := db.Debug().Model(&User{}).Where("id = ?", buys[i].AuthorID).Take(&buys[i].Author).Error
			if err != nil {
				return &[]Buy{}, err
			}
		}
	}
	return &buys, nil
}

func (p *Buy) FindBuyByID(db *gorm.DB, pid uint64) (*Buy, error) {
	var err error
	err = db.Debug().Model(&Buy{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Buy{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Buy{}, err
		}
	}
	return p, nil
}

func (p *Buy) UpdateABuy(db *gorm.DB) (*Buy, error) {

	var err error

	err = db.Debug().Model(&Buy{}).Where("id = ?", p.ID).Updates(Buy{BitcoinAmount: p.BitcoinAmount}).Error
	if err != nil {
		return &Buy{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Buy{}, err
		}
	}
	return p, nil
}

func (p *Buy) DeleteABuy(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Buy{}).Where("id = ? and author_id = ?", pid, uid).Take(&Buy{}).Delete(&Buy{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Buy not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
