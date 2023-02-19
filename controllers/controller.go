package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type pagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}

type pagination struct {
	c *gin.Context
	query *gorm.DB
	records interface{}
}

func (p *pagination) paginate() *pagingResult {
	page, _ := strconv.Atoi(p.c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(p.c.DefaultQuery("limit", "12"))

	//ใช้ go routine เพื่อ obtimize การเก็บค่า count
	ch := make(chan int)
	go p.countRecords(ch)

	//find records
	//limit offset
	//limit => 10
	//page => 1, 1-10 offset => 0
	//page => 2, 11-20 offset => 1
	offset := (page - 1) * limit
	p.query.Limit(limit).Offset(offset).Find(p.records)

	count := <- ch
	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	var nextPage int
	if page == totalPage {
		nextPage = totalPage
	} else {
		// nextPage = page - 1
		nextPage = page + 1
	}

	return &pagingResult{
		Page:      page,
		Limit:     limit,
		Count:     count,
		PrevPage:  page - 1,
		NextPage:  nextPage,
		TotalPage: totalPage,
	}
}

func (p *pagination) countRecords(ch chan int)  {
	var count int
	p.query.Model(p.records).Count(&count)

	ch <- count
}
