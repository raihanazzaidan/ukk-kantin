package main

import (
	"backend_golang/controllers"
	"backend_golang/middlewares"
	"backend_golang/setup"

	"github.com/gin-gonic/gin"
)

func main() {
	//Declare New Gin Route System
	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	//Run Database Setup
	setup.ConnectDatabase()

	// APIs that doesn't need middleware
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	// APIs that need middleware
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/logout", controllers.Logout)
	protected.GET("/user", controllers.GetCurrentUser)

	// User APIs
	protected.GET("/user/all", controllers.GetAllUser)
	protected.GET("/user/:id", controllers.GetUserById)
	// protected.POST("/user/add", controllers.AddUser)
	protected.PUT("/user/update/:id", controllers.UpdateUser)
	protected.PUT("/user/passupdate/:id", controllers.ResetPassword)
	protected.DELETE("/user/delete/:id", controllers.DeleteUser)

	// Siswa APIs
	protected.GET("/siswa/all", controllers.GetAllSiswa)
	protected.GET("/siswa/:id", controllers.GetSiswaById)
	protected.POST("/siswa/add", controllers.AddSiswa)
	protected.PUT("/siswa/update/:id", controllers.UpdateSiswa)
	protected.DELETE("/siswa/delete/:id", controllers.DeteleSiswa)

	// Stan APIs
	protected.GET("/stan/all", controllers.GetAllStan)
	protected.GET("/stan/:id", controllers.GetStanById)
	protected.POST("/stan/add", controllers.AddStan)
	protected.PUT("/stan/update/:id", controllers.UpdateStan)
	protected.DELETE("/stan/delete/:id", controllers.DeleteStan)

	// Menu APIs
	protected.GET("/menu/all", controllers.GetAllMenu)
	protected.GET("/menu/stan/:stan_id", controllers.GetMenuByStanId)
	protected.POST("/menu/add", controllers.AddMenu)
	protected.PUT("/menu/update/:id", controllers.UpdateMenu)
	protected.DELETE("/menu/delete/:id", controllers.DeleteMenu)

	// Diskon APIs
	protected.GET("/diskon/all", controllers.GetAllDiskon)
	protected.GET("/diskon/:id", controllers.GetDiskonById)
	protected.POST("/diskon/add", controllers.AddDiskon)
	protected.PUT("/diskon/update/:id", controllers.UpdateDiskon)
	protected.DELETE("/diskon/delete/:id", controllers.DeleteDiskon)



	// Port
	router.Run(":6969")
}
