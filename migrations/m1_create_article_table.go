package migrations

import (
	"goBlogApp/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1CreatearticleTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1",
		Migrate: func(tx *gorm.DB) error{
			return tx.AutoMigrate(&models.Article{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("articles").Error
		},
	}
}