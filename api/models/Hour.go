package models

import "github.com/jinzhu/gorm"

type LastHour struct {
	Hora   int `gorm:"not null;" json:"hora"`
	Minuto int `gorm:"not null;" json:"minuto"`
}

func FindLastHour(db *gorm.DB) (*[]LastHour, error) {
	var err error
	hour := []LastHour{}
	err = db.Debug().Model(&LastHour{}).Limit(100).Find(&hour).Error
	if err != nil {
		return &[]LastHour{}, err
	}
	return &hour, nil
}
