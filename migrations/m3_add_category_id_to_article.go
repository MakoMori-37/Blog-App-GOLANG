package migrations

import (
	"goBlogApp/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m3AddCategoryIdtoArticle() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3",
		Migrate: func(tx *gorm.DB) error{
			err := tx.AutoMigrate(&models.Article{}).Error
			var articles []models.Article
			tx.Find(&articles)
			for _, article := range articles {
				article.CategoryID = 1
				tx.Save(&article)
			}

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Model(&models.Article{}).DropColumn("category_id").Error
		},
	}
}