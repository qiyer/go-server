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
	Coin    uint64 `form:"coin"`    // 多少金币，默认 0
}

type CoinAutoQueueRequest struct {
	CoinAutoQueue []CoinAutoRequest `form:"queue"`  //金币增长序列
	UserID        string            `form:"userID"` // 用户 id
}

type UserInfoRequest struct {
	Name   string `form:"name"`   //金币增长序列
	UserID string `form:"userID"` // 用户 id
}

type CheckInRequest struct {
	UserID string `form:"userID"` // 用户 id
}

type OnlineRewardsRequest struct {
	UserID string `form:"userID"` // 用户 id
}

type LevelUpRequest struct {
	UserID string `form:"userID"` // 用户 id
	RoleID string `form:"roleID"` // 角色 id
	Level  int    `form:"level"`  // 升了几个等级，默认 1
}

type PassChapterRequest struct {
	UserID  string `form:"userID"`  // 用户 id
	Chapter int    `form:"chapter"` // 章节
}
