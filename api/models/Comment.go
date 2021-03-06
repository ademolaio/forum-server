package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
	"time"
)

type Comment struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	PostID    uint64    `gorm:"not null" json:"post_id"`
	Body      string    `gorm:"text;not null" json:"body"`
	User      User      `json:"user"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare is a function that ensures all requirements are there
func (c *Comment) Prepare() {
	c.ID = 0
	c.Body = html.EscapeString(strings.TrimSpace(c.Body))
	c.User = User{}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

// Validate is a function that ensures a post exists before you comment
func (c *Comment) Validate(action string) map[string]string {
	var errorMesssages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if c.Body == "" {
			err = errors.New("Required Comment")
			errorMesssages["required_body"] = err.Error()
		}
	default:
		if c.Body == "" {
			err = errors.New("Required Comment")
			errorMesssages["required_body"] = err.Error()
		}
	}
	return errorMesssages
}

// SaveComment is function that allows you to save comment to a post
func (c *Comment) SaveComment(db *gorm.DB) (*Comment, error) {
	err := db.Debug().Create(&c).Error
	if err != nil {
		return &Comment{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.UserID).Take(&c.User).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return c, nil
}

// GetComments is function to get comments from a post
func (c *Comment) GetComments(db *gorm.DB, pid uint64) (*[]Comment, error) {
	comments := []Comment{}
	err := db.Debug().Model(&Comment{}).Where("post_id = ?", pid).Order("created_at desc").Find(&comments).Error
	if err != nil {
		return &[]Comment{}, err
	}
	if len(comments) > 0 {
		for i, _ := range comments {
			err := db.Debug().Model(&User{}).Where("id = ?", comments[i].UserID).Take(&comments[i].User).Error
			if err != nil {
				return &[]Comment{}, err
			}
		}
	}
	return &comments, err
}

// UpdateAComment is a function that enables user(s) to edit a comment
func (c *Comment) UpdateAComment(db *gorm.DB) (*Comment, error) {
	var err error

	err = db.Debug().Model(&Comment{}).Where("id = ?", c.ID).Updates(Comment{Body: c.Body, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Comment{}, err
	}

	fmt.Println("This is the comment body: ", c.Body)
	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.UserID).Take(&c.User).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return c, nil
}

// DeleteAComment is function that enables user(s) to delete a comment
func (c *Comment) DeleteAComment(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Comment{}).Where("id = ?", c.ID).Take(&Comment{}).Delete(&Comment{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// DeleteUserComments when a User is deleted, we also delete comments that the post had
func (c *Comment) DeleteUserComments(db *gorm.DB, uid uint32) (int64, error) {
	comments := []Comment{}
	db = db.Debug().Model(&Comment{}).Where("user_id = ?", uid).Find(&comments).Delete(&comments)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// DeletePostComments when a post is deleted, we also delete the comments that the post had
func (c *Comment) DeletePostComments(db *gorm.DB, pid uint64) (int64, error) {
	comments := []Comment{}
	db = db.Debug().Model(&Comment{}).Where("post_id = ?", pid).Find(&comments).Delete(&comments)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
