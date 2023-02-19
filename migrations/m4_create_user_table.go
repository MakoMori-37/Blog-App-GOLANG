package migrations

import (
	"goBlogApp/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m4CreateUsertable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID:"4",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.User{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("users").Error
		},
	}
}