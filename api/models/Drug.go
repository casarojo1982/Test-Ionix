package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Drug struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	Approved  bool      `gorm:"size:255;not null;" json:"approved"`
	Min_dose  uint32    `gorm:"not null" json:"min_dose"`
	Max_dose  uint32    `gorm:"not null" json:"max_dose"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (d *Drug) Prepare() {
	d.ID = 0
	d.Name = html.EscapeString(strings.TrimSpace(d.Name))
	d.Approved = d.Approved
	d.Min_dose = d.Min_dose
	d.Max_dose = d.Max_dose
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
}

func (d *Drug) Validate() error {

	if d.Name == "" {
		return errors.New("Required Name")
	}

	return nil
}

func (d *Drug) SaveDrug(db *gorm.DB) (*Drug, error) {
	var err error
	err = db.Debug().Model(&Drug{}).Create(&d).Error
	if err != nil {
		return &Drug{}, err
	}
	return d, nil
}

func (d *Drug) FindAllDrug(db *gorm.DB) (*[]Drug, error) {
	var err error
	drugs := []Drug{}
	err = db.Debug().Model(&Drug{}).Limit(100).Find(&drugs).Error
	if err != nil {
		return &[]Drug{}, err
	}
	return &drugs, nil
}

func (d *Drug) FindDrugByID(db *gorm.DB, id uint64) (*Drug, error) {
	var err error
	err = db.Debug().Model(&Drug{}).Where("id = ?", id).Take(&d).Error
	if err != nil {
		return &Drug{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Drug{}, errors.New("Drug Not Found")
	}
	return d, nil
}

func (d *Drug) UpdateADrug(db *gorm.DB, id uint64) (*Drug, error) {

	var err error
	db = db.Debug().Model(&Drug{}).Where("id = ?", id).Take(&Drug{}).UpdateColumns(
		map[string]interface{}{
			"name":       d.Name,
			"approved":   d.Approved,
			"min_dose":   d.Approved,
			"max_dose":   d.Approved,
			"updated_at": time.Now(),
		},
	)
	err = db.Debug().Model(&Drug{}).Where("id = ?", id).Take(&d).Error
	if err != nil {
		return &Drug{}, err
	}
	// This is the display the updated user
	err = db.Debug().Model(&Drug{}).Where("id = ?", id).Take(&d).Error
	if err != nil {
		return &Drug{}, err
	}
	return d, nil
}

func (d *Drug) DeleteADrug(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&Drug{}).Where("id = ?", id).Take(&Drug{}).Delete(&Drug{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
