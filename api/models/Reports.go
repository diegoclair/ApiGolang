package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Reports struct {
	Buys 	Buy		`json:"buys"`
	Sales	Sale	`json:"sales"` 
}

func FindReportsByUserID(db *gorm.DB, uid uint32) (, error) {
	var err error
	if err = db.Joins("JOIN buys on buys.author_id=users.id").
	Joins("JOIN sales on sales.author_id=users.id").Where("users.id=?",uid).
	Group("users.id").Find(&artists).Error; err != nil {
		log.Fatal(err)
	}
	/* err = db.Debug().Joins("users").Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	} */
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}


/*Get the list of all artists who acted in movie “Nayagan
if err = db.Joins("JOIN artist_movies on artist_movies.artist_id=artists.id").
	Joins("JOIN movies on artist_movies.movie_id=movies.id").Where("movies.title=?", "Nayagan").
	Group("artists.id").Find(&artists).Error; err != nil {
		log.Fatal(err)
}
*/