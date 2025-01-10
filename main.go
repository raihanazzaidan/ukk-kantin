package main

import (
	// "backend_golang/controllers"
	// "backend_golang/middlewares"
	"backend_golang/setup"

	"github.com/gin-gonic/gin"
)

func main() {
	//Declare New Gin Route System
	router := gin.New()
	// router.Use(middleware.CORSMiddleware())
	//Run Database Setup
	setup.ConnectDatabase()




	router.Run(":6969")
}