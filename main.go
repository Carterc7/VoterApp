package main

import (
	"html/template"
	"log"
	"os"
	"voterapp/controllers"
	"voterapp/db"
	"voterapp/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set in environment")
	}

	db.Connect(mongoURI)

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

	router := gin.Default()
	tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*"))
	router.SetHTMLTemplate(tmpl)
	router.Static("/assets", "./assets")

	router.GET("/", controllers.ShowHomePage)
	router.GET("/poll", controllers.SearchPoll)
	router.GET("/poll/:id", controllers.ShowPoll)
	router.POST("/vote", controllers.Vote)
	router.GET("/poll/:id/results", controllers.ShowResults)
	router.GET("/register", controllers.ShowRegisterForm)
	router.POST("/register", controllers.Register)
	router.GET("/login", controllers.ShowLoginForm)
	router.POST("/login", controllers.Login)
	router.GET("/logout", controllers.Logout)

	auth := router.Group("/")
	auth.Use(middlewares.RequireLogin())
	auth.GET("/poll/new", controllers.ShowCreatePollForm)
	auth.POST("/poll/new", controllers.CreatePoll)
	auth.GET("/mypolls", controllers.ShowUserPolls)
	auth.POST("/poll/:id/delete", controllers.DeletePoll)

	router.Run(":8082")
}
