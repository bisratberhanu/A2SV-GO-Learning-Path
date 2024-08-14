package repositories

import (
	"context"
	"fmt"
	"task_manger_clean_architecture/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	database   *mongo.Database
	collection string
}

// UpdateAllToken implements domain.UserRepository.
func (u *userRepository) UpdateAllTokens(token string, refreshToken string, user_id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
	var updateObject primitive.D
	updatedAt,_ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObject = append(updateObject, bson.E{Key: "updatedat",  Value: updatedAt })
	upsert:= true
	filter:= bson.M{"userid": user_id}
	opt:= options.UpdateOptions{Upsert: &upsert,}
    collection := u.database.Collection(u.collection)
    _,err := collection.UpdateOne(ctx, filter,bson.D{{"$set", updateObject}},&opt)
	defer cancel()
	
    return err
}

func NewUserRepository(db *mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

// GetUser implements domain.UserRepository.
func (u *userRepository) GetUser(c context.Context, user_id string) (domain.User, error) {
	var user domain.User
	collection := u.database.Collection(u.collection)
	err := collection.FindOne(c, bson.M{"userid": user_id}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, err
}

// GetUsers implements domain.UserRepository.
func (u *userRepository) GetUsers(ctx context.Context, startIndex int64, recordsPerPage int64) ([]*domain.User, error) {
	var allUsers []*domain.User

	matchStage := bson.D{{"$match", bson.D{}}}
	groupStage := bson.D{{"$group", bson.D{
		{"_id", "null"},
		{"total_count", bson.D{{"$sum", 1}}},
		{"data", bson.D{{"$push", "$$ROOT"}}},
	}}}
	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordsPerPage}}}},
		}},
	}

	cursor, err := u.database.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &allUsers); err != nil {
		return nil, err
	}

	if len(allUsers) > 0 {
		return allUsers, nil
	}

	return nil, fmt.Errorf("no users found")
}

// Login implements domain.UserRepository.
func (u *userRepository) Login(ctx context.Context, email string) (*domain.User, error) {
	collection := u.database.Collection(u.collection)
	var user domain.User

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Promote implements domain.UserRepository.
func (u *userRepository) Promote(ctx context.Context, user_id string, userType string) error {
	collection := u.database.Collection(u.collection)

	// Update the user type to ADMIN
	filter := bson.M{"userid": user_id}
	update := bson.M{"$set": bson.M{"usertype": userType, "updatedat": time.Now()}}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating user type: %v", err)
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	if res.ModifiedCount == 0 {
		return fmt.Errorf("user type was already %s", userType)
	}

	return nil
}

// Signup implements domain.UserRepository.
func (u *userRepository) Signup(ctx context.Context, user domain.User) (interface{}, error) {
	// Set timestamps and IDs
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.ID = primitive.NewObjectID()
	user.UserId = user.ID.Hex()

	// Insert the user into the database
	result, err := u.database.Collection(u.collection).InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("user item was not created: %v", err)
	}

	return result.InsertedID, nil
}
