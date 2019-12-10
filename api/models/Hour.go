package models

import "github.com/jinzhu/gorm"

type LastHour struct {
	Hour   int `gorm:"not null;" json:"hour"`
	Minute int `gorm:"not null;" json:"minute"`
}

func FindLastHour(db *gorm.DB) ([]LastHour, error) {
	var err error
	hour := []LastHour{}
	err = db.Debug().Model(&LastHour{}).Limit(10).Find(&hour).Error
	if err != nil {
		return []LastHour{}, err
	}
	return hour, nil
}

func (t *LastHour) UpdateLastHour(db *gorm.DB) (string, error) {

	db = db.Debug().Model(&LastHour{}).Take(&LastHour{}).UpdateColumns(
		map[string]interface{}{
			"hour":  t.Hour,
			"minute":  t.Minute,
		},
	)
	if db.Error != nil {
		return "Error", db.Error
	}
	
	return "Ok", nil
}
