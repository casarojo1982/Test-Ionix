package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/victorsteven/fullstack/api/models"
)

var users = []models.User{
	models.User{
		Name:     "Robinson Rodriguez",
		Email:    "rjr@gmail.com",
		Password: "123456",
	},
	models.User{
		Name:     "Rebecca Ferguson",
		Email:    "becaferguson@gmail.com",
		Password: "123456",
	},
}

var drugs = []models.Drug{
	models.Drug{
		Name:     "Aspirina",
		Approved: true,
		Min_dose: 1,
		Max_dose: 2,
	},
	models.Drug{
		Name:     "Brugesic",
		Approved: true,
		Min_dose: 1,
		Max_dose: 3},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Drug{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Drug{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

	for i, _ := range drugs {
		err = db.Debug().Model(&models.Drug{}).Create(&drugs[i]).Error
		if err != nil {
			log.Fatalf("cannot seed drugs table: %v", err)
		}
	}
}
