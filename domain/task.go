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

type CoinAutoRequest struct {
	Clicker int    `form:"clicker"` // 连点击，默认 1
	IsAd    bool   `form:"isAd"`    // 是否开启广告，默认关闭
	Times   int    `form:"times"`   // 多倍收益，默认倍数 1
	UserID  string `form:"userID"`  // 用户 id
}
