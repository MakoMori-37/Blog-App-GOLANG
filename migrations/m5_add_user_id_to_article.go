package migrations

import (
	"goBlogApp/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m5AddUserIdtoArticle() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "5",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Article{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Model(&models.Article{}).DropColumn("user_id").Error
		},
	}
}
