package controllers

import (
	"context"
	"net/http"
	"time"
	"voterapp/db"
	"voterapp/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func ShowHomePage(c *gin.Context) {
	collection := db.GetPollsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"public": true})
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to fetch polls")
		return
	}
	defer cursor.Close(ctx)

	var polls []models.Poll
	if err := cursor.All(ctx, &polls); err != nil {
		c.String(http.StatusInternalServerError, "Error decoding polls")
		return
	}

	// Check if user is logged in by looking for the cookie
	_, err = c.Cookie("user_id")
	isLoggedIn := (err == nil)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":      "Home",
		"polls":      polls,
		"IsLoggedIn": isLoggedIn,
	})
}
