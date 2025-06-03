package main

import (
	"html/template"
	"voterapp/controllers"
	"voterapp/db"
	"voterapp/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	mongoURI := "mongodb+srv://lawsoncards7:Superman7@voterapp.twlxnhd.mongodb.net/?retryWrites=true&w=majority&appName=VoterApp"
	db.Connect(mongoURI)

	// Register custom functions for use in templates
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"divf": func(a, b int) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b)
		},
		"mul": func(a float64, b int) float64 { return a * float64(b) },
	}

	// Initialize Gin router
	router := gin.Default()

	// Register the custom functions and load the templates
	tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*"))
	router.SetHTMLTemplate(tmpl)

	// Static assets
	router.Static("/assets", "./assets")

	// Public Routes
	router.GET("/", controllers.ShowHomePage)
	router.GET("/poll", controllers.SearchPoll)
	router.GET("/poll/:id", controllers.ShowPoll)
	router.POST("/vote", controllers.Vote)
	router.GET("/poll/:id/results", controllers.ShowResults)

	// Auth Routes
	router.GET("/register", controllers.ShowRegisterForm)
	router.POST("/register", controllers.Register)
	router.GET("/login", controllers.ShowLoginForm)
	router.POST("/login", controllers.Login)
	router.GET("/logout", controllers.Logout)

	// Protected Routes (Require Login)
	auth := router.Group("/")
	auth.Use(middlewares.RequireLogin())
	auth.GET("/poll/new", controllers.ShowCreatePollForm)
	auth.POST("/poll/new", controllers.CreatePoll)
	auth.GET("/mypolls", controllers.ShowUserPolls)
	auth.POST("/poll/:id/delete", controllers.DeletePoll)

	// Run server
	router.Run(":8082")
}
