package controllers

import (
	"goBlogApp/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Categories struct {
	DB *gorm.DB
}

type CategoryResponse struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Articles []struct {
		ID uint `json:"id"`
		Title string `json:"title"`
	} `json:"articles"`
}
type AllCategoryResponse struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type createCategoryForm struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc" binding:"required"`
}

type updateCategoryForm struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (c *Categories) GetList(ctx *gin.Context) {
	var categories []models.Category
	c.DB.Order("id desc").Find(&categories)

	var response []AllCategoryResponse
	copier.Copy(&response, &categories)
	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func (c *Categories) GetDetail(ctx *gin.Context) {
	category, err := c.FindCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var response CategoryResponse
	copier.Copy(&response, &category)
	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func (c *Categories) Create(ctx *gin.Context) {
	var form createCategoryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	copier.Copy(&category, &form)
	if err := c.DB.Create(&category).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var response CategoryResponse
	copier.Copy(&response, &category)
	ctx.JSON(http.StatusCreated, gin.H{"data": response})
}

func (c *Categories) Update(ctx *gin.Context) {
	var form updateCategoryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	category, err := c.FindCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := c.DB.Model(&category).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var response CategoryResponse
	copier.Copy(&response, &category)
	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func (c *Categories) Delete(ctx *gin.Context) {
	category, err := c.FindCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.DB.Unscoped().Delete(&category)
	ctx.Status(http.StatusNoContent)
}

func (c *Categories) FindCategoryByID(ctx *gin.Context) (*models.Category, error) {
	var categories models.Category
	id := ctx.Param("id")

	if err := c.DB.Preload("Articles").First(&categories, id).Error; err != nil {
		return nil, err
	}

	return &categories, nil
}
