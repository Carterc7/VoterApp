package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Poll struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserID   primitive.ObjectID `bson:"user_id"`
	Question string             `bson:"question"`
	Options  []string           `bson:"options"`
	Votes    []int              `bson:"votes"`
	Public   bool               `bson:"public"`
}
