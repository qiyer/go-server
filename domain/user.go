package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUser = "users"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"createdAt"`         // 创建时间
	UpdatedAt time.Time          `bson:"updatedAt"`         // 最后更新时间
	UnionID   string             `bson:"unionId,omitempty"` // 系统唯一标识（用于多平台账号合并）
	Email     string             `bson:"email,omitempty"`   // 主邮箱（可空）
	Phone     string             `bson:"phone,omitempty"`   // 主手机号（可空）
	Verified  bool               `bson:"verified"`          // 是否已验证主联系方式
	Providers []AuthProvider     `bson:"providers"`         // 第三方登录信息数组
}

type AuthProvider struct {
	Platform    string    `bson:"platform"`    // 平台标识：wechat/facebook/google
	PlatformID  string    `bson:"platformId"`  // 第三方平台用户唯一ID
	AccessToken string    `bson:"accessToken"` // 访问令牌
	Metadata    bson.M    `bson:"metadata"`    // 平台返回的原始数据（如微信用户信息）
	ExpiresAt   time.Time `bson:"expiresAt"`   // 令牌过期时间
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	Fetch(c context.Context) ([]User, error)
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id string) (User, error)
}
