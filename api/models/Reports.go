package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Reports struct {
	Buys  []Buy  `json:"buys"`
	Sales []Sale `json:"sales"`
}

func (r *Reports) FindReportsByUserID(db *gorm.DB, uid uint32) (*Reports, error) {

	//========================================= BUYS ==============================================
	buys := []Buy{}
	var err error
	err = db.Debug().Model(&Buy{}).Where("author_id = ?", uid).Scan(&buys).Error
	if err != nil {
		return &Reports{}, err
	}
	if len(buys) > 0 {
		for i := range buys {
			if uid != 0 {
				err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&buys[i].Author).Error
				if err != nil {
					return &Reports{}, err
				}
			}
		}
	}
	r.Buys = buys

	//========================================= SALES ==============================================
	sales := []Sale{}
	err = db.Debug().Model(&Sale{}).Where("author_id = ?", uid).Scan(&sales).Error
	if err != nil {
		return &Reports{}, err
	}
	if len(sales) > 0 {
		for i := range sales {
			if uid != 0 {
				err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&sales[i].Author).Error
				if err != nil {
					return &Reports{}, err
				}
			}
		}
	}
	r.Sales = sales
	return r, nil
}

func (r *Reports) FindReportsByDate(db *gorm.DB, date string) (*Reports, error) {
	fmt.Println(date)
	//========================================= BUYS ==============================================
	buys := []Buy{}
	var err error
	err = db.Debug().Model(&Buy{}).Where("date(created_at) = ?", date).Scan(&buys).Error
	if err != nil {
		return &Reports{}, err
	}
	if len(buys) > 0 {
		for i := range buys {
			err = db.Debug().Model(&User{}).Where("id = ?", buys[i].AuthorID).Take(&buys[i].Author).Error
			if err != nil {
				return &Reports{}, err
			}
		}
	}
	r.Buys = buys

	//========================================= SALES ==============================================
	sales := []Sale{}
	err = db.Debug().Model(&Sale{}).Where("date(created_at) = ?", date).Scan(&sales).Error
	if err != nil {
		return &Reports{}, err
	}
	if len(sales) > 0 {
		for i := range sales {
			err = db.Debug().Model(&User{}).Where("id = ?", sales[i].AuthorID).Take(&sales[i].Author).Error
			if err != nil {
				return &Reports{}, err
			}
		}
	}
	r.Sales = sales
	return r, nil
}
