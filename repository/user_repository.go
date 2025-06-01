package repository

import (
	"context"
	"fmt"
	"time"

	"go-server/domain"
	"go-server/internal/tokenutil"
	"go-server/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database // Add this at the top of the file to define the DB variable.
var ContextTimeout time.Duration

func Create(c context.Context, user *domain.User) error {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()

	collection := (*DB).Collection(domain.CollectionUser)
	_, err := collection.InsertOne(c, user)

	return err
}

func CreateAccount(c context.Context, account *domain.Account) error {

	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()

	collection := (*DB).Collection(domain.CollectionAccount)
	_, err := collection.InsertOne(c, account)

	return err
}

func CreateUserMapping(c context.Context, userMapping *domain.UserMapping) error {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()

	collection := (*DB).Collection(domain.CollectionUserMapping)
	_, err := collection.InsertOne(c, userMapping)

	return err
}

func Fetch(c context.Context) ([]domain.User, error) {
	collection := (*DB).Collection(domain.CollectionUser)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(c, bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	var users []domain.User

	err = cursor.All(c, &users)
	if users == nil {
		return []domain.User{}, err
	}

	return users, err
}

func GetByEmail(c context.Context, email string) (domain.Account, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()

	collection := (*DB).Collection(domain.CollectionAccount)
	var account domain.Account
	err := collection.FindOne(c, bson.M{"email": email}).Decode(&account)
	return account, err
}

func GetByID(c context.Context, id primitive.ObjectID) (domain.User, error) {
	collection := (*DB).Collection(domain.CollectionUser)

	var user domain.User

	err := collection.FindOne(c, bson.M{"_id": id}).Decode(&user)
	return user, err
}

func GetUserMappingByPId(c context.Context, platformId string) (domain.UserMapping, error) {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()

	collection := (*DB).Collection(domain.CollectionUserMapping)
	var userMapping domain.UserMapping
	err := collection.FindOne(c, bson.M{"platformId": platformId}).Decode(&userMapping)
	return userMapping, err
}

func CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, secret, expiry)
}

func CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}

func GetUserByEmail(c context.Context, email string) (domain.User, domain.Account, error) {
	ctx, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	account, err := GetByEmail(ctx, email)
	fmt.Printf("GetUserByEmail GetByEmail%+v\n", account)
	if err != nil {
		return domain.User{}, account, err
	}
	userMapping, err := GetUserMappingByPId(ctx, account.ID.Hex())
	fmt.Printf("GetUserByEmail GetUserMappingByPId%+v\n", userMapping)
	if err != nil {
		return domain.User{}, account, err
	}
	user, err := GetByID(ctx, userMapping.UserId)
	fmt.Printf("GetUserByEmail GetByID%+v\n", user)
	return user, account, err
}

func GetUserByID(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	user, _, err := GetUserByEmail(ctx, email)
	return user, err
}

func UpdateUserDays(c context.Context, id primitive.ObjectID, days []string) error {
	_, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	collection := (*DB).Collection(domain.CollectionUser)

	// 构建过滤条件
	filter := bson.M{"_id": id}

	// 定义更新操作（使用 $set 精确更新字段）
	update := bson.M{
		"$set": bson.M{
			"days": days,
		},
	}

	// 执行更新
	_, err := collection.UpdateOne(context.Background(), filter, update)

	return err
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	return tokenutil.ExtractIDFromToken(requestToken, secret)
}

func GetProfileByID(c context.Context, userID string) (*domain.Profile, error) {
	// ctx, cancel := context.WithTimeout(c, ContextTimeout)
	// defer cancel()

	// user, err := GetByID(ctx, userID)
	// if err != nil {
	// 	return nil, err
	// }

	// return &domain.Profile{Name: user.Name, Email: user.ID.Hex()}, nil
	return nil, nil
}
