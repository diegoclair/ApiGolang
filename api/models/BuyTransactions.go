package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Buy struct {
	ID				uint64    `gorm:"primary_key;auto_increment" json:"id"`
	BitcoinAmount	string    `gorm:"size:255;not null;" json:"bitcoin_amount"`
	Author			User      `json:"author"`
	AuthorID		uint32    `gorm:"not null" json:"author_id"`
	CreatedAt		time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (p *Buy) Prepare() {
	p.ID = 0
	p.BitcoinAmount = html.EscapeString(strings.TrimSpace(p.BitcoinAmount))
	p.Author = User{}
	p.CreatedAt = time.Now()
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