package models

import (
	"github.com/jinzhu/gorm"
	"html"
	"strings"
)

type ResetPassword struct {
	gorm.Model
	Email string `gorm:"size:100;not null;" json:"email"`
	Token string `gorm:"size:255;not null;" json:"token"`
}

// Prepare function for password reset
func (resetPassword *ResetPassword) Prepare() {
	resetPassword.Token = html.EscapeString(strings.TrimSpace(resetPassword.Token))
	resetPassword.Email = html.EscapeString(strings.TrimSpace(resetPassword.Email))

}

// SaveDetails function to save the changes during Password Reset
func (resetPassword *ResetPassword) SaveDetails(db *gorm.DB) (*ResetPassword, error) {
	var err error
	err = db.Debug().Create(&resetPassword).Error
	if err != nil {
		return &ResetPassword{}, err
	}
	return resetPassword, nil
}

// DeleteDetails function to delete changes during Password Resent
func (resetPassword *ResetPassword) DeleteDetails(db *gorm.DB) (int64, error) {

	db = db.Debug().Model(&ResetPassword{}).Where("id = ?", resetPassword.ID).Take(&ResetPassword{}).Delete(&ResetPassword{})
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil

}
