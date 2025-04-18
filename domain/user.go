package domain

import (
	"context"
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
	ID             primitive.ObjectID `bson:"_id"`
	Name           string             `bson:"name"`
	Avatar         string             `bson:"avatar,omitempty"` // 头像URL（可空）
	Coins          int64              `bson:"coins"`            // 用户金币数量
	Level          int                `bson:"level"`            // 用户等级
	Experience     int64              `bson:"experience"`       // 用户经验值
	Chapter        int                `bson:"chapter"`          // 当前章节
	Vip            int                `bson:"vip"`              // VIP等级
	OnlineTime     int64              `bson:"onlineTime"`       // 在线时长（秒）
	Days           int                `bson:"days"`             // 连续登录天数
	GiftExp        int64              `bson:"giftExp"`          // 礼物奖励当前百分比
	ChallengeLevel int                `bson:"challengeLevel"`   // 挑战等级
	Bosses         []Boss             `bson:"bosses"`           // 用户拥有的Boss列表
	Grils          []Gril             `bson:"grils"`            // 用户拥有的秘书列表
	GrilTrainLevel int                `bson:"grilTrainLevel"`   // 秘书训练等级
	Capital        Capital            `bson:"capital"`          // 用户拥有的资产
	Builds         []Build            `bson:"builds"`           // 用户拥有的建筑列表
	CreatedAt      time.Time          `bson:"createdAt"`        // 创建时间【用来计算以及多少天】
	UpdatedAt      time.Time          `bson:"updatedAt"`        // 最后更新时间
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

type Gril struct {
	PetId     string    `bson:"petId"`
	Level     string    `bson:"level"`     // 宠物等级
	CreatedAt time.Time `bson:"createdAt"` // 创建时间
}

type Boss struct {
	BossId string    `bson:"bossId"`
	Time   time.Time `bson:"time"` // 当前时间
}

type Capital struct {
	CapitalIds []string  `bson:"capitalIds"` //资产ID列表
	Time       time.Time `bson:"time"`       // 冷却时间
}

type Build struct {
	BuildId   string    `bson:"buildId"`
	Level     string    `bson:"level"`     // 等级
	CreatedAt time.Time `bson:"createdAt"` // 创建时间
	UpdatedAt time.Time `bson:"updatedAt"` // 最后更新时间
}

type UserRepository interface {
	CreateAccount(c context.Context, account *Account) error
	CreateUserMapping(c context.Context, userMapping *UserMapping) error
	Create(c context.Context, user *User) error
	Fetch(c context.Context) ([]User, error)
	GetByEmail(c context.Context, email string) (Account, error)
	GetByID(c context.Context, id string) (User, error)
}
