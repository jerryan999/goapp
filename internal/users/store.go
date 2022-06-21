package users

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Database       = "goapp"
	UserCollection = "user"
)

type store interface {
	Create(ctx context.Context, u *User) error
	ReadByEmail(ctx context.Context, email string) (*User, error)
}

type userStore struct {
	mongoClient    *mongo.Client
	userCollection *mongo.Collection
}

func (us *userStore) Create(ctx context.Context, u *User) error {
	_, err := us.userCollection.InsertOne(ctx, u)
	if err != nil {
		return fmt.Errorf("userstore create: %w", err)
	}
	return nil
}

func (us *userStore) ReadByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := us.userCollection.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&u)
	if err != nil {
		return nil, fmt.Errorf("userstore readbyEmail: %w", err)
	}
	return &u, nil
}

func newStore(mongoClient *mongo.Client) (*userStore, error) {
	return &userStore{
		mongoClient:    mongoClient,
		userCollection: mongoClient.Database(Database).Collection(UserCollection),
	}, nil
}
