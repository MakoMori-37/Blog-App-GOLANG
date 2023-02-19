package routes

import (
	"goBlogApp/config"
	"goBlogApp/controllers"
	"goBlogApp/middleware"

	"github.com/gin-gonic/gin"
)

func Serve(r *gin.Engine) {
	db := config.GetDB()
	v1 := r.Group("/api/v1")
	authenticate := middleware.Authenticate().MiddlewareFunc()
	authorize := middleware.Authorize()

	authController := controllers.Auth{DB: db}
	authGroup := v1.Group("auth")
	{
		authGroup.POST("/sign-up", authController.Signup)
		authGroup.POST("/sign-in", middleware.Authenticate().LoginHandler)
		authGroup.GET("/profile", authenticate, authController.GetProfile)
		authGroup.PATCH("/profile", authenticate, authController.UpdateProfile)
	}
	
	usersController := controllers.Users{DB: db}
	usersGroup := v1.Group("users")
	usersGroup.Use(authenticate)
	{
		usersGroup.GET("", usersController.GetList)
		usersGroup.POST("", usersController.Create)
		usersGroup.GET("/:id", usersController.GetDetail)
		usersGroup.PATCH("/:id", usersController.Update)
		usersGroup.DELETE("/:id", usersController.Delete)
		usersGroup.PATCH("/:id/promote", usersController.Promote)
		usersGroup.PATCH("/:id/demote", usersController.Demote)
	}

	articleController := controllers.Articles{DB: db}
	articlesGroup := v1.Group("articles")
	articlesGroup.GET("", articleController.GetList)
	articlesGroup.GET("/:id", articleController.GetDetail)
	articlesGroup.Use(authenticate)
	{
		articlesGroup.POST("", authenticate, articleController.Create)
		articlesGroup.PATCH("/:id", articleController.Update)
		articlesGroup.DELETE("/:id", articleController.Delete)
	}

	categoryController := controllers.Categories{DB: db}
	categoryGroup := v1.Group("categories")
	categoryGroup.Use(authenticate)
	{
		categoryGroup.GET("", categoryController.GetList)
		categoryGroup.GET("/:id", categoryController.GetDetail)
		categoryGroup.POST("", categoryController.Create)
		categoryGroup.PATCH("/:id", categoryController.Update)
		categoryGroup.DELETE("/:id", categoryController.Delete)
	}

	dashboardController := controllers.Dashboard{DB: db}
	dashboardGroup := v1.Group("dashboard")
	dashboardGroup.Use(authenticate, authorize)
	{
		dashboardGroup.GET("", dashboardController.GetInfo)
	}
}
