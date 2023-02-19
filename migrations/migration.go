package migrations

import (
	"goBlogApp/config"
	"log"

	"gopkg.in/gormigrate.v1"
)

func Migrate() {
	db := config.GetDB()
	m := gormigrate.New(
		db,
		gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			m1CreatearticleTable(),
			m2CreateCategoryTable(),
			m3AddCategoryIdtoArticle(),
			m4CreateUsertable(),
			m5AddUserIdtoArticle(),
		},
	)

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}