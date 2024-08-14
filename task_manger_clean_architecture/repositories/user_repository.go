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

// GetUser implements domain.UserRepository.
func (u *userRepository) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	var user domain.User
	collection := u.database.Collection(u.collection)
	err := collection.FindOne(c, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, err
}

func (u *userRepository) GetUsers(ctx context.Context, startIndex int64, recordsPerPage int64) ([]*domain.User, error) {
    var allUsers []*domain.User

    collection := u.database.Collection(u.collection)

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

    cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
    if err != nil {
        fmt.Println("Aggregation error: ", err)
        return nil, fmt.Errorf("aggregation error: %v", err)
    }
    defer cursor.Close(ctx)

    var result []bson.M
    if err := cursor.All(ctx, &result); err != nil {
        fmt.Println("Cursor error: ", err)
        return nil, fmt.Errorf("cursor error: %v", err)
    }

    if len(result) == 0 {
        fmt.Println("No users found")
        return nil, fmt.Errorf("no users found")
    }

    userItems, ok := result[0]["user_items"].(primitive.A)
    if !ok {
        fmt.Println("Unexpected result format: 'user_items' is not of type primitive.A")
        fmt.Printf("Actual type of 'user_items': %T\n", result[0]["user_items"])  // Log actual type
        return nil, fmt.Errorf("unexpected result format")
    }

    for _, item := range userItems {
        userMap, ok := item.(primitive.M)
        if !ok {
            fmt.Println("Unexpected user item format")  // Log unexpected format error
            fmt.Printf("Item type: %T, Item: %+v\n", item, item)  // Inspect the item
            return nil, fmt.Errorf("unexpected user item format")
        }

        var user domain.User
        // Manually map fields from userMap to user
        if id, ok := userMap["_id"].(primitive.ObjectID); ok {
            user.ID = id
        }
        if firstname, ok := userMap["firstname"].(string); ok {
            user.FirstName = &firstname
        }
        if lastname, ok := userMap["lastname"].(string); ok {
            user.LastName = &lastname
        }
        if password, ok := userMap["password"].(string); ok {
            user.Password = &password
        }
        if email, ok := userMap["email"].(string); ok {
            user.Email = &email
        }
        if phone, ok := userMap["phone"].(string); ok {
            user.Phone = &phone
        }
        if token, ok := userMap["token"].(string); ok {
            user.Token = &token
        }
        if usertype, ok := userMap["usertype"].(string); ok {
            user.UserType = &usertype
        }
        if refreshtoken, ok := userMap["refreshtoken"].(string); ok {
            user.RefreshToken = &refreshtoken
        }
        if createdat, ok := userMap["createdat"].(int64); ok {
            user.CreatedAt = time.UnixMilli(createdat)
        }
        if updatedat, ok := userMap["updatedat"].(int64); ok {
            user.UpdatedAt = time.UnixMilli(updatedat)
        }
        if userid, ok := userMap["userid"].(string); ok {
            user.UserId = userid
        }

        allUsers = append(allUsers, &user)
    }

    return allUsers, nil
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
func (u *userRepository) Promote(ctx context.Context, user_id string, userType string)( error,int64, int64 ){
	collection := u.database.Collection(u.collection)

	// Update the user type to ADMIN
	filter := bson.M{"userid": user_id}
	update := bson.M{"$set": bson.M{"usertype": userType, "updatedat": time.Now()}}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating user type: %v", err),0,0
	}

	

	return nil, res.MatchedCount,res.ModifiedCount
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


// UpdateAllToken implements domain.UserRepository.
func (u *userRepository) UpdateAllTokens(token string, refreshToken string, user_id string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second * 100)
	var updateObject primitive.D
	updatedAt,_ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObject = append(updateObject, bson.E{Key: "updatedat",  Value: updatedAt })
	updateObject = append(updateObject, bson.E{Key: "token", Value: token})
	updateObject = append(updateObject, bson.E{Key: "refreshtoken", Value: refreshToken})
	upsert:= true
	filter:= bson.M{"userid": user_id}
	opt:= options.UpdateOptions{Upsert: &upsert,}
    collection := u.database.Collection(u.collection)
    _,err := collection.UpdateOne(ctx, filter,bson.D{{"$set", updateObject}},&opt)
	defer cancel()
	
    return err
}

