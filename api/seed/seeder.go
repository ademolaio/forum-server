package seed

import (
	"github.com/jinzhu/gorm"
	"github.sandbox.com/forum-app/forum-server/api/models"
	"log"
)

var users = []models.User{
	models.User{
		Username: "Paul",
		Email:    "apostle@example.com",
		Password: "password",
	},

	models.User{
		Username: "Peter",
		Email:    "disciple@example.com",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "The Zealot",
		Content: "I saw our Lord and Saviour Jesus on my way to Demasucs",
	},
	models.Post{
		Title:   "The Walk",
		Content: "I walk with Jesus through his entire ministry",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}, &models.Like{}, &models.Comment{}).Error
	if err != nil {
		log.Fatalf("Cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("Cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("Attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed posts table: %v", err)
		}
	}
}
