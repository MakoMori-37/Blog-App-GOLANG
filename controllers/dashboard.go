package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Dashboard struct {
	DB *gorm.DB
}

type dashboardArticle struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Image   string `json:"image"`
}

type dashboardResponse struct {
	LatestArticles []dashboardArticle `json:"atestArticles"`
	UsersCount     []struct{
		Role  string `json:"role"`
		Count uint   `json:"count"`
	} `json:"usersCount"`
	CategoriesCount uint `json:"categoriesCount"`
	ArticlesCount   uint `json:"articlesCount"`
}

func (d *Dashboard) GetInfo(ctx *gin.Context) {
	res := dashboardResponse{}
	d.DB.Table("articles").Order("id desc").Limit(5).Find(&res.LatestArticles)
	d.DB.Table("articles").Count(&res.ArticlesCount)
	d.DB.Table("categories").Count(&res.CategoriesCount)
	d.DB.Table("users").Select("role, count(*)").Group("role").Scan(&res.UsersCount)

	ctx.JSON(http.StatusOK, gin.H{"data": &res})
}