package controllers

import (
	"context"
	"net/http"
	"sort"
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

	// Fetch public polls
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

	// Sort polls by total vote count (descending)
	sort.Slice(polls, func(i, j int) bool {
		sumVotes := func(votes []int) int {
			sum := 0
			for _, vote := range votes {
				sum += vote
			}
			return sum
		}
		return sumVotes(polls[i].Votes) > sumVotes(polls[j].Votes)
	})

	// Take top 5
	if len(polls) > 5 {
		polls = polls[:5]
	}

	// Check if user is logged in
	_, err = c.Cookie("user_id")
	isLoggedIn := (err == nil)
	isLoggedOut := !isLoggedIn

	// Render page
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":       "Home",
		"polls":       polls,
		"IsLoggedIn":  isLoggedIn,
		"IsLoggedOut": isLoggedOut,
	})
}
