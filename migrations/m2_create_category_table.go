package migrations

import (
	"goBlogApp/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m2CreateCategoryTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2",
		Migrate: func(tx *gorm.DB) error{
			return tx.AutoMigrate(&models.Category{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("categories").Error
		},
	}
}