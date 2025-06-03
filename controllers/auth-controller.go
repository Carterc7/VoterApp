package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"
	"voterapp/db"
	"voterapp/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func ShowRegisterForm(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func Register(c *gin.Context) {
	username := strings.TrimSpace(c.PostForm("username"))
	password := c.PostForm("password")

	if username == "" || len(password) < 4 {
		c.String(http.StatusBadRequest, "Invalid input")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: string(hashedPassword),
	}

	collection := db.GetUsersCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.String(http.StatusInternalServerError, "Registration failed")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func ShowLoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	collection := db.GetUsersCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		c.String(http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Set session or cookie
	c.SetCookie("user_id", user.ID.Hex(), 86400, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}

func Logout(c *gin.Context) {
	c.SetCookie("user_id", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}
