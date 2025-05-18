package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUser        = "users"
	CollectionUserMapping = "user_mappings"
	CollectionAccount     = "accounts"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id" json:"userId"`
	Name           string             `bson:"name" json:"name"`
	Avatar         string             `bson:"avatar,omitempty" json:"avatar"`       // 头像URL（可空）
	Coins          uint64             `bson:"coins" json:"coins"`                   // 用户金币数量
	Level          int                `bson:"level" json:"level"`                   // 用户等级
	Experience     int64              `bson:"experience" json:"experience"`         // 用户经验值
	Chapter        int                `bson:"chapter" json:"chapter"`               // 当前章节
	Vip            int                `bson:"vip" json:"vip"`                       // VIP等级
	OnlineTime     int64              `bson:"onlineTime" json:"onlineTime"`         // 在线时长（秒）
	Days           []string           `bson:"days" json:"days"`                     // 连续登录天数
	GiftExp        int64              `bson:"giftExp" json:"giftExp"`               // 礼物奖励当前百分比
	ChallengeLevel int                `bson:"challengeLevel" json:"challengeLevel"` // 挑战等级
	Bosses         []Boss             `bson:"bosses" json:"bosses"`                 // 用户拥有的Boss列表
	Girls          string             `bson:"girls" json:"girls"`                   // 用户拥有的秘书列表
	GirlTrainLevel int                `bson:"grilTrainLevel" json:"grilTrainLevel"` // 秘书训练等级
	Capital        Capital            `bson:"capital" json:"capital"`               // 用户拥有的资产
	Build          Build              `bson:"build" json:"build"`                   // 用户拥有的小区
	Islands        []Island           `bson:"islands" json:"islands"`               // 用户拥有的岛屿列表
	Legacies       []Legacy           `bson:"legacies" json:"legacies"`             // 用户拥有的遗迹列表
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`           // 创建时间【用来计算以及多少天】
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`           // 最后更新时间
}

type UserMapping struct {
	PlatformID string             `bson:"platformId"` // 第三方平台用户唯一ID
	UserId     primitive.ObjectID `bson:"userId"`     // User在本系统的唯一ID
	Platform   string             `bson:"platform"`   // 平台标识：wechat/facebook/google
	Metadata   bson.M             `bson:"metadata"`   // 平台返回的原始数据（如微信用户信息）
	CreateAt   time.Time          `bson:"createAt"`   // 创建时间
}

type Account struct {
	ID        primitive.ObjectID `bson:"_id"`
	AccountId string             `bson:"accountId"` // 账号ID
	Password  string             `bson:"password"`
	Email     string             `bson:"email,omitempty"` // 主邮箱（可空）
	Phone     string             `bson:"phone,omitempty"` // 主手机号（可空）
	CreatedAt time.Time          `bson:"createdAt"`       // 创建时间
	UpdatedAt time.Time          `bson:"updatedAt"`       // 最后更新时间
}

type MGirl struct {
	GirlId uint   `bson:"girlId" json:"girlId"`
	Level  uint64 `bson:"level" json:"level"` // 宠物等级
}

type Boss struct {
	BossId string    `bson:"bossId" json:"bossId"`
	Time   time.Time `bson:"time" json:"time"` // 当前时间
}

type Capital struct {
	CapitalIds []string  `bson:"capitalIds" json:"capitalIds"` //资产ID列表
	Time       time.Time `bson:"time" json:"time"`             // 冷却时间
}

type Build struct {
	Id           string    `bson:"buildId" json:"buildId"`
	Level        uint      `bson:"level" json:"level"`               // 等级
	DisplayLevel uint      `bson:"displayLevel" json:"displayLevel"` // 显示等级
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`       // 创建时间
	UpdatedAt    time.Time `bson:"updatedAt" json:"updatedAt"`       // 最后更新时间
}

type Island struct {
	Id    string `bson:"id" json:"id"`
	Level uint   `bson:"level" json:"level"` // 等级
}

type Legacy struct {
	Id    string `bson:"id" json:"id"`
	Level uint   `bson:"level" json:"level"` // 等级
}
