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

type ClickEarnRequest struct {
	Clicker int `form:"clicker"` // 连点击，默认 1
}

type UserInfoRequest struct {
	Name   string `form:"name"`   //金币增长序列
	UserID string `form:"userID"` // 用户 id
}

type LevelUpRequest struct {
	RoleID uint `form:"roleID"` // 角色 id
	Level  int  `form:"level"`  // 升了几个等级，默认 1
}

type PassChapterRequest struct {
	Chapter int `form:"chapter"` // 章节
}

type UnLockRoleRequest struct {
	RoleID uint `form:"roleID"` // 角色 id
}

type UnLockVehicleRequest struct {
	VehicleID uint `form:"vehicleID"` // 坐骑 id
}

type UnLockCapitalRequest struct {
	CapitalID uint `form:"capitalID"` //资产 id
}

type CheckInRequest struct {
	Id int `form:"id"`
}

type OnlineRequest struct {
	Id int `form:"id"`
}

type VehicleDisplayRequest struct {
	DisplayLevel int `form:"displayLevel"` // 显示等级
}

type BossRequest struct {
	ID int `form:"id"` // boss id
}
