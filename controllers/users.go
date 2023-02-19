package controllers

import (
	"goBlogApp/models"
	"goBlogApp/config"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Users struct {
	DB *gorm.DB
}

type createUserForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type updateUserForm struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=8"`
	Name     string `json:"name"`
}

type usersPaging struct {
	Items  []userResponse `json:"items"`
	Paging *pagingResult  `json:"paging"`
}

type userResponse struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Role   string `json:"role"`
}

func (u *Users) GetList(ctx *gin.Context) {
	var users []models.User
	query := u.DB.Order("id desc").Find(&users)

	term := ctx.Query("term")
	if term != "" {
		query.Where("name ILIKE ?", "%" + term + "%")
	}

	pagination := pagination{c: ctx, query: query, records: &users}
	paging := pagination.paginate()

	response := []userResponse{}
	copier.Copy(&response, &users)
	ctx.JSON(http.StatusOK, gin.H{
		"data": usersPaging{Items: response, Paging: paging},
	})
}

func (u *Users) GetDetail(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var response userResponse
	copier.Copy(&response, &user)
	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (u *Users) Create(ctx *gin.Context) {
	var form createUserForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptPassword()

	if err := u.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var response userResponse
	copier.Copy(&response, &user)
	ctx.JSON(http.StatusCreated, gin.H{
		"data": response,
	})
}

func (u *Users) Update(ctx *gin.Context) {
	var form updateUserForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if form.Password != "" {
		user.Password = user.GenerateEncryptPassword()
	}

	if err := u.DB.Model(&user).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var response userResponse
	copier.Copy(&response, &user)
	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (u *Users) Delete(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	u.DB.Delete(&user)

	ctx.Status(http.StatusNoContent)
}

func (u *Users) Promote(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Promote()
	u.DB.Save(user)

	var response userResponse
	copier.Copy(&response, &user)
	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (u *Users) Demote(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Demote()
	u.DB.Save(user)

	var response userResponse
	copier.Copy(&response, &user)
	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (u *Users) findUserByID(ctx *gin.Context) (*models.User, error) {
	id := ctx.Param("id")
	var user models.User

	if err := u.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func setUserImage(ctx *gin.Context, user *models.User) error {
	file, _ := ctx.FormFile("avatar")
	if file == nil {
		return nil
	}

	if user.Avatar != "" {
		user.Avatar = strings.Replace(user.Avatar, os.Getenv("Host"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + user.Avatar)
	}

	path := "uploads/users/" + strconv.Itoa(int(user.ID))
	os.MkdirAll(path, os.ModePerm)
	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}

	db := config.GetDB()
	user.Avatar = os.Getenv("HOST") + "/" + filename
	db.Save(user)

	return nil
}
