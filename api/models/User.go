package models

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"github.sandbox.com/forum-app/forum-server/api/security"
	"html"
	"log"
	"os"
	"strings"
	"time"
)

type User struct {
	ID         uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Username   string    `gorm:"size:255;not null;unique" json:"username"`
	Email      string    `gorm:"size:100;not null;unique" json:"email"`
	Password   string    `gorm:"size:100;not null;" json:"password"`
	AvatarPath string    `gorm:"size:255;null;" json:"avatar_path"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

//BeforeSave function is to check the password before User's account is created
func (u *User) BeforeSave() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

//Prepare function is meant to do a checkup prior to User signup
func (u *User) Prepare() {
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

//AfterFind function
func (u *User) AfterFind() (err error) {
	if err != nil {
		return err
	}
	if u.AvatarPath != "" {
		u.AvatarPath = os.Getenv("DO_SPACES_URL") + u.AvatarPath
	}
	return nil
}

func (u *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["invalid_emails"] = err.Error()
			}
		}

	case "login":
		if u.Password == "" {
			err = errors.New("Required Password")
			errorMessages["required_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["invalid_email"] = err.Error()
			}
		}

	case "forgotpassword":
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["invalid_email"] = err.Error()
			}
		}

	default:
		if u.Username == "" {
			err = errors.New("Required Username")
			errorMessages["required_username"] = err.Error()
		}
		if u.Password == "" {
			err = errors.New("Reqired Password")
			errorMessages["required_password"] = err.Error()
		}
		if u.Password != "" && len(u.Password) < 10 {
			err = errors.New("Password should be at least 1- characters")
			errorMessages["invalid_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Reqired Email")
			errorMessages["required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["invalid_email"] = err.Error()
			}
		}
	}
	return errorMessages
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

//FindUserByID allows one to find users by their ID
func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

// UpdateAUser function allows the user to update their info
func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {
	if u.Password != "" {
		//To has the password
		err := u.BeforeSave()
		if err != nil {
			log.Fatal(err)
		}

		db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
			map[string]interface{}{
				"password":   u.Password,
				"email":      u.Email,
				"updated_at": time.Now(),
			},
		)
	}

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"email":      u.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}

	// This is to display the updated user
	err := db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *User) UpdatePassword(db *gorm.DB) error {

	// To has the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&User{}).Where("email = ?", u.Email).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":   u.Password,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
