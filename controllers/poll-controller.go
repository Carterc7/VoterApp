package controllers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"
	"voterapp/db"
	"voterapp/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ShowUserPolls(c *gin.Context) {
	userIDHex, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	_, err = c.Cookie("user_id")
	isLoggedIn := (err == nil)

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid user ID")
		return
	}

	collection := db.GetPollsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to load your polls")
		return
	}
	defer cursor.Close(ctx)

	var polls []models.Poll
	if err = cursor.All(ctx, &polls); err != nil {
		c.String(http.StatusInternalServerError, "Failed to decode polls")
		return
	}

	c.HTML(http.StatusOK, "user-polls.html", gin.H{
		"polls":      polls,
		"IsLoggedIn": isLoggedIn,
	})
}

func DeletePoll(c *gin.Context) {
	userIDHex, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid user ID")
		return
	}

	pollIDStr := c.Param("id")
	pollID, err := primitive.ObjectIDFromHex(pollIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid poll ID")
		return
	}

	collection := db.GetPollsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{
		"_id":     pollID,
		"user_id": userID,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete poll")
		return
	}
	if res.DeletedCount == 0 {
		c.String(http.StatusForbidden, "Poll not found or not authorized to delete")
		return
	}

	c.Redirect(http.StatusSeeOther, "/mypolls")
}

func ShowPoll(c *gin.Context) {
	idStr := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid poll ID")
		return
	}

	collection := db.GetPollsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var poll models.Poll
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&poll)
	if err != nil {
		c.String(http.StatusNotFound, "Poll not found")
		return
	}

	_, err = c.Cookie("user_id")
	isLoggedIn := (err == nil)

	c.HTML(http.StatusOK, "poll.html", gin.H{
		"poll":       poll,
		"IsLoggedIn": isLoggedIn,
	})
}

func Vote(c *gin.Context) {
	idStr := c.PostForm("id")
	optionStr := c.PostForm("option")

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid poll ID")
		return
	}

	option, err := strconv.Atoi(optionStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid option")
		return
	}

	collection := db.GetPollsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var poll models.Poll
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&poll)
	if err != nil {
		c.String(http.StatusNotFound, "Poll not found")
		return
	}

	if option < 0 || option >= len(poll.Options) {
		c.String(http.StatusBadRequest, "Invalid vote option")
		return
	}

	poll.Votes[option]++

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"votes": poll.Votes}})
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to record vote")
		return
	}

	c.Redirect(http.StatusSeeOther, "/poll/"+idStr+"/results")
}

func SearchPoll(c *gin.Context) {
	idStr := c.Query("id")
	if len(idStr) != 24 {
		c.String(http.StatusBadRequest, "Poll ID must be 24 characters")
		return
	}

	c.Redirect(http.StatusSeeOther, "/poll/"+idStr)
}

func ShowResults(c *gin.Context) {
	idStr := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid poll ID")
		return
	}

	collection := db.GetPollsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var poll models.Poll
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&poll)
	if err != nil {
		c.String(http.StatusNotFound, "Poll not found")
		return
	}

	_, err = c.Cookie("user_id")
	isLoggedIn := (err == nil)

	// Compute total votes
	totalVotes := 0
	for _, v := range poll.Votes {
		totalVotes += v
	}

	// Compute bar percentages as floats
	barPercents := make([]float64, len(poll.Votes))
	for i, votes := range poll.Votes {
		if totalVotes > 0 {
			barPercents[i] = (float64(votes) / float64(totalVotes)) * 100
		} else {
			barPercents[i] = 0
		}
	}

	c.HTML(http.StatusOK, "results.html", gin.H{
		"poll":        poll,
		"IsLoggedIn":  isLoggedIn,
		"BarPercents": barPercents,
	})
}

func ShowCreatePollForm(c *gin.Context) {
	_, err := c.Cookie("user_id")
	isLoggedIn := (err == nil)

	c.HTML(http.StatusOK, "create_poll.html", gin.H{
		"IsLoggedIn": isLoggedIn,
	})
}

func CreatePoll(c *gin.Context) {
	userIDHex, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid user ID")
		return
	}

	question := c.PostForm("question")
	optionsRaw := c.PostForm("options")
	publicStr := c.PostForm("public")

	options := strings.Split(optionsRaw, "\n")

	cleanOptions := []string{}
	for _, opt := range options {
		opt = strings.TrimSpace(opt)
		if opt != "" {
			cleanOptions = append(cleanOptions, opt)
		}
	}

	if question == "" || len(cleanOptions) < 2 {
		c.String(http.StatusBadRequest, "Question and at least 2 options are required.")
		return
	}

	isPublic := publicStr == "on"

	poll := models.Poll{
		ID:       primitive.NewObjectID(),
		UserID:   userID,
		Question: question,
		Options:  cleanOptions,
		Votes:    make([]int, len(cleanOptions)),
		Public:   isPublic,
	}

	collection := db.GetPollsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, poll)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create poll")
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}
