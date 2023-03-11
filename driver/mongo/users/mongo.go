package users

import (
	"context"
	"marketplace-backend/business/users"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) users.Repository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

/*
Create
*/

func (ur *userRepository) Create(domain *users.Domain) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := ur.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return users.Domain{}, err
	}

	result, err := ur.GetByID(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return users.Domain{}, err
	}

	return result, err
}

/*
Read
*/

func (ur *userRepository) GetByID(id primitive.ObjectID) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := ur.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)
	if err != nil {
		return users.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (ur *userRepository) GetByEmail(email string) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := ur.collection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&result)

	return result.ToDomain(), err
}

/*
Update
*/

func (ur *userRepository) Update(domain *users.Domain) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := ur.collection.UpdateOne(ctx, bson.M{
		"_id": domain.ID,
	}, bson.M{
		"$set": FromDomain(domain),
	})
	if err != nil {
		return users.Domain{}, err
	}

	result, err := ur.GetByID(domain.ID)
	if err != nil {
		return users.Domain{}, err
	}

	return result, nil
}

/*
Delete
*/
