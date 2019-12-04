package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Sale struct {
	ID				uint64    `gorm:"primary_key;auto_increment" json:"id"`
	BitcoinAmount	string    `gorm:"size:255;not null;" json:"bitcoin_amount"`
	Author			User      `json:"author"`
	AuthorID		uint32    `gorm:"not null" json:"author_id"`
	CreatedAt		time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (p *Sale) Prepare() {
	p.ID = 0
	p.BitcoinAmount = html.EscapeString(strings.TrimSpace(p.BitcoinAmount))
	p.Author = User{}
	p.CreatedAt = time.Now()
}

func (p *Sale) Validate() error {

	if p.BitcoinAmount == "" {
		return errors.New("Required BitcoinAmount")
	}
	if p.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (p *Sale) SaveSale(db *gorm.DB) (*Sale, error) {
	var err error
	err = db.Debug().Model(&Sale{}).Create(&p).Error
	if err != nil {
		return &Sale{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Sale{}, err
		}
	}
	return p, nil
}

func (p *Sale) FindAllSales(db *gorm.DB) (*[]Sale, error) {
	var err error
	sales := []Sale{}
	err = db.Debug().Model(&Sale{}).Limit(100).Find(&sales).Error
	if err != nil {
		return &[]Sale{}, err
	}
	if len(sales) > 0 {
		for i, _ := range sales {
			err := db.Debug().Model(&User{}).Where("id = ?", sales[i].AuthorID).Take(&sales[i].Author).Error
			if err != nil {
				return &[]Sale{}, err
			}
		}
	}
	return &sales, nil
}

func (p *Sale) FindSaleByID(db *gorm.DB, pid uint64) (*Sale, error) {
	var err error
	err = db.Debug().Model(&Sale{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Sale{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Sale{}, err
		}
	}
	return p, nil
}

func (p *Sale) UpdateASale(db *gorm.DB) (*Sale, error) {

	var err error

	err = db.Debug().Model(&Sale{}).Where("id = ?", p.ID).Updates(Sale{BitcoinAmount: p.BitcoinAmount}).Error
	if err != nil {
		return &Sale{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Sale{}, err
		}
	}
	return p, nil
}

func (p *Sale) DeleteASale(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Sale{}).Where("id = ? and author_id = ?", pid, uid).Take(&Sale{}).Delete(&Sale{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Sale not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}