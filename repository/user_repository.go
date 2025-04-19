package repository

import (
	"context"
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

func GetByID(c context.Context, id string) (domain.User, error) {
	collection := (*DB).Collection(domain.CollectionUser)

	var user domain.User

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&user)
	return user, err
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
	if err != nil {
		user, err := GetByID(ctx, account.AccountId)
		return user, account, err
	}
	return domain.User{}, account, err
}

func GetUserByID(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()
	return GetByID(ctx, email)
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	return tokenutil.ExtractIDFromToken(requestToken, secret)
}

func GetProfileByID(c context.Context, userID string) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(c, ContextTimeout)
	defer cancel()

	user, err := GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.Profile{Name: user.Name, Email: user.ID.Hex()}, nil
}
