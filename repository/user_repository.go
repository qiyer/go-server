package repository

import (
	"context"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	database              mongo.Database
	userCollection        string
	userMappingCollection string
	accountCollection     string
}

func NewUserRepository(db mongo.Database, userCollection string, userMappingCollection string, accountCollection string) domain.UserRepository {
	return &userRepository{
		database:              db,
		userCollection:        userCollection,
		userMappingCollection: userMappingCollection,
		accountCollection:     accountCollection,
	}
}

func (ur *userRepository) Create(c context.Context, user *domain.User) error {
	collection := ur.database.Collection(ur.userCollection)

	_, err := collection.InsertOne(c, user)

	return err
}

func (ur *userRepository) CreateAccount(c context.Context, account *domain.Account) error {
	collection := ur.database.Collection(ur.accountCollection)

	_, err := collection.InsertOne(c, account)

	return err
}

func (ur *userRepository) CreateUserMapping(c context.Context, userMapping *domain.UserMapping) error {
	collection := ur.database.Collection(ur.userMappingCollection)

	_, err := collection.InsertOne(c, userMapping)

	return err
}

func (ur *userRepository) Fetch(c context.Context) ([]domain.User, error) {
	collection := ur.database.Collection(ur.userCollection)

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

func (ur *userRepository) GetByEmail(c context.Context, email string) (domain.Account, error) {
	collection := ur.database.Collection(ur.accountCollection)
	var account domain.Account
	err := collection.FindOne(c, bson.M{"email": email}).Decode(&account)
	return account, err
}

func (ur *userRepository) GetByID(c context.Context, id string) (domain.User, error) {
	collection := ur.database.Collection(ur.userCollection)

	var user domain.User

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&user)
	return user, err
}
