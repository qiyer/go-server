package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionTask = "tasks"
)

type Task struct {
	ID     primitive.ObjectID `bson:"_id" json:"-"`
	Title  string             `bson:"title" form:"title" binding:"required" json:"title"`
	UserID primitive.ObjectID `bson:"userID" json:"-"`
}

type CoinAuto struct {
	Clicker int                `bson:"clicker" json:"clicker"` // Clicker is the number of clicks
	IsAd    bool               `bson:"isAd" json:"isAd"`
	Times   int                `bson:"times" json:"times"`
	UserID  primitive.ObjectID `bson:"userID" json:"userID"`
}
