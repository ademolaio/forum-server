package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Like struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	PostID    uint64    `gorm:"not null" json:"post_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// SaveLike the functions that allows user(s) to have persist their like on a post
func (l *Like) SaveLike(db *gorm.DB) (*Like, error) {

	// Check if the author user has liked this post before:
	err := db.Debug().Model(&Like{}).Where("post_id = ? AND user_id = ?", l.PostID, l.UserID).Take(&l).Error
	if err != nil {
		if err.Error() == "record not found" {
			// The user has not liked this post before, so lets save incoming like:
			err = db.Debug().Model(&Like{}).Create(&l).Error
			if err != nil {
				return &Like{}, err
			}
		}
	} else {
		// The user has liked before, so create a custom error message
		err = errors.New("double like")
		return &Like{}, err
	}
	return l, nil
}

// DeleteLike allows user(s) to remove a like on a post
func (l *Like) DeleteLike(db *gorm.DB) (*Like, error) {
	var err error
	var deletedLike *Like

	err = db.Debug().Model(Like{}).Where("id = ?", l.ID).Take(&l).Error
	if err != nil {
		return &Like{}, err
	} else {

		// If the like exist, save it in deleted like and delete it
		deletedLike = l
		db = db.Debug().Model(&Like{}).Where("id = ?", l.ID).Take(&Like{}).Delete(&Like{})
		if db.Error != nil {
			fmt.Println("Can't delete like: ", db.Error)
			return &Like{}, db.Error
		}
	}
	return deletedLike, nil
}

// GetLikesInfo allows one to get the Info of Likes
func (l *Like) GetLikesInfo(db *gorm.DB, pid uint64) (*[]Like, error) {
	likes := []Like{}
	err := db.Debug().Model(&Like{}).Where("post_id ?", pid).Find(&likes).Error
	if err != nil {
		return &[]Like{}, err
	}
	return &likes, err
}

// DeleteUsersLikes When a post is deleted, we also delete the likes that the post had
func (l *Like) DeleteUsersLikes(db *gorm.DB, uid uint32) (int64, error) {
	likes := []Like{}
	db = db.Debug().Model(&Like{}).Where("user_id = ?", uid).Find(&likes).Delete(&likes)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (l *Like) DeletePostLikes(db *gorm.DB, pid uint64) (int64, error) {
	likes := []Like{}
	db = db.Debug().Model(&Like{}).Where("post_id = ?", pid).Find(&likes).Delete(&likes)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
