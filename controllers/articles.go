package controllers

import (
	"goBlogApp/models"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Articles struct {
	DB *gorm.DB
}

type createArticleForm struct {
	Title   string                `form:"title" binding:"required"`
	Body    string                `form:"body" binding:"required"`
	Excerpt string                `form:"excerpt" binding:"required"`
	CategoryID uint                `form:"categoryId" binding:"required"`
	Image   *multipart.FileHeader `form:"image" binding:"required"`
}
type updateArticleForm struct {
	Title   string                `form:"title"`
	Body    string                `form:"body"`
	Excerpt string                `form:"excerpt"`
	CategoryID uint               `form:"categoryId"`
	Image   *multipart.FileHeader `form:"image"`
}
type articleResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Excerpt    string `json:"excerpt"`
	Image      string `json:"image"`
	CategoryID uint   `json:"categoryId"`
	Category   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	User struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"user"`
}
type craeteOrUpdateResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Excerpt    string `json:"excerpt"`
	Image      string `json:"image"`
	UserID uint `json:"userId"`
	CategoryID uint   `json:"categoryId"`
}

type articlesPaging struct {
	Item   []articleResponse `json:"item"`
	Paging *pagingResult     `json:"paging"` //ตรงนี้ที่ใส่เป็น pointer เพราะอยากให้มันประหยัด mem ไม่ต้องไป copy
}

func (a *Articles) GetList(c *gin.Context) {
	var articles []models.Article

	query := a.DB.Preload("User").Preload("Category").Order("id desc")
	categoryId := c.Query("categoryId")
	if categoryId != "" {
		query.Where("category_id = ?", categoryId)
	}
	term := c.Query("term")
	if term != "" {
		query.Where("title ILIKE ?", "%" + term + "%")
	}

	a.DB.Find(&articles)

	pagination := pagination{c: c, query: query, records: &articles}
	paging := pagination.paginate()

	response := []articleResponse{}
	copier.Copy(&response, &articles)

	c.JSON(http.StatusOK, gin.H{"data": articlesPaging{
		Item:   response,
		Paging: paging,
	}})
}

func (a *Articles) GetDetail(c *gin.Context) {
	article, err := a.FindArticleByID(c)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := articleResponse{}
	copier.Copy(&response, &article)

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (a *Articles) Create(c *gin.Context) {

	var form createArticleForm

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	//from this
	// article := models.Article{
	// 	Title: form.Title,
	// ..................
	// }

	//to this
	var article models.Article
	user, _ := c.Get("sub")
	copier.Copy(&article, &form)
	article.User = *user.(*models.User)

	if err := a.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	a.setArticleImage(c, &article)

	response := craeteOrUpdateResponse{}
	copier.Copy(&response, &article)

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

func (a *Articles) Update(c *gin.Context) {
	var form updateArticleForm

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	article, err := a.FindArticleByID(c)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := a.DB.Model(&article).Update(&form).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	a.setArticleImage(c, article)

	var response craeteOrUpdateResponse
	copier.Copy(&response, &article)

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (a *Articles) Delete(c *gin.Context) {
	article, err := a.FindArticleByID(c)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	a.DB.Delete(&article)
	c.Status(http.StatusNoContent)
}

func (a *Articles) setArticleImage(c *gin.Context, article *models.Article) error {
	file, err := c.FormFile("image")
	if err != nil || file == nil {
		return err
	}

	if article.Image != "" {
		article.Image = strings.Replace(article.Image, os.Getenv("HOST"), "", 1)
		//หา Directory ปัจจุบัน
		pwd, _ := os.Getwd()
		os.Remove(pwd + article.Image)
	}

	//check if update
	path := "uploads/articles/" + strconv.Itoa(int(article.ID))
	os.MkdirAll(path, 0755)
	filename := path + "/" + file.Filename
	if err := c.SaveUploadedFile(file, filename); err != nil {
		return err
	}

	article.Image = os.Getenv("HOST") + "/" + filename
	a.DB.Save(article)

	return nil
}

func (a *Articles) FindArticleByID(c *gin.Context) (*models.Article, error) {
	var article models.Article
	id := c.Param("id")

	if err := a.DB.Preload("User").Preload("Category").First(&article, id).Error; err != nil {
		return nil, err
	}

	return &article, nil
}

