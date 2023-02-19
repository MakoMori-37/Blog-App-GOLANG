package controllers

import (
	"goBlogApp/models"
	"mime/multipart"
	"net/http"


	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Auth struct {
	DB *gorm.DB
}

type authForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type updateProfileForm struct {
	Email  string                `form:"email"`
	Name   string                `form:"name"`
	Avatar *multipart.FileHeader `form:"avatar"`
}

type authResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func (a *Auth) Signup(ctx *gin.Context) {
	var form authForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptPassword()
	if err := a.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var response authResponse
	copier.Copy(&response, &user)
	ctx.JSON(http.StatusCreated, gin.H{"data": response})
}

func (a *Auth) GetProfile(ctx *gin.Context) {
	sub, _ := ctx.Get("sub")
	user := sub.(*models.User)

	var response userResponse
	copier.Copy(&response, &user)
	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func (a *Auth) UpdateProfile(ctx *gin.Context) {
	var form updateProfileForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	sub, _ := ctx.Get("sub")
	user := sub.(*models.User)

	setUserImage(ctx, user)
	if err := a.DB.Model(user).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var response userResponse
	copier.Copy(&response, user)
	ctx.JSON(http.StatusOK, gin.H{"data": response})
}
