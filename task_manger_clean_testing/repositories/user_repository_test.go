package repositories

import (
	"context"
	"task_manger_clean_architecture/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	mockRepo       domain.UserRepository
	mongoClient    *mongo.Client
	testDatabase   *mongo.Database
	testCollection string
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	suite.Require().NoError(err)

	// Ping the MongoDB server to ensure connection
	err = client.Ping(context.Background(), nil)
	suite.Require().NoError(err)

	// Set up test database and collection
	suite.mongoClient = client
	suite.testDatabase = client.Database("test_database")
	suite.testCollection = "test_users"

	// Initialize the repository with the test database and collection
	suite.mockRepo = NewUserRepository(suite.testDatabase, suite.testCollection)
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	// Drop the test database
	err := suite.testDatabase.Drop(context.Background())
	suite.Require().NoError(err)

	// Disconnect the MongoDB client
	err = suite.mongoClient.Disconnect(context.Background())
	suite.Require().NoError(err)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	// Clear the collection after each test
	err := suite.testDatabase.Collection(suite.testCollection).Drop(context.Background())
	suite.Require().NoError(err)
}
func (suite *UserRepositoryTestSuite) TestSignup() {
	firstName := "Test"
	lastName := "User"
	password := "password123"
	email := "testuser@example.com"
	phone := "1234567890"
	userType := "USER"
	
	ID:= primitive.NewObjectID()
	UserId:= ID.Hex()
	newUser := domain.User{
		FirstName: &firstName,
		LastName:  &lastName,
		Password:  &password,
		Email:     &email,
		Phone:     &phone,
		UserType:  &userType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserId: UserId,
		ID: ID,
	}

	insertedID, err := suite.mockRepo.Signup(context.Background(), newUser)
	suite.NoError(err)

	insertedIDStr, ok := insertedID.(primitive.ObjectID)
	suite.Require().True(ok, "Inserted ID is not of type primitive.ObjectID")

	// Log the inserted ID for debugging
	suite.T().Logf("Inserted ID: %s", insertedIDStr.Hex())

	// Verify user was inserted
	fetchedUser, err := suite.mockRepo.GetUser(context.Background(), UserId)
	suite.NoError(err)

	// Log the fetched user ID for debugging
	suite.T().Logf("Fetched User ID: %s", fetchedUser.UserId)

	suite.Equal(UserId, fetchedUser.UserId)
}


func (suite *UserRepositoryTestSuite) TestGetUser() {
	// Step 1: Add a new user to the database
	firstName := "Test"
	lastName := "User"
	password := "password123"
	email := "testuser@example.com"
	phone := "1234567890"
	userType := "USER"
	newUser := domain.User{
		ID:        primitive.NewObjectID(),
		FirstName: &firstName,
		LastName:  &lastName,
		Password:  &password,
		Email:     &email,
		Phone:     &phone,
		UserType:  &userType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserId:    primitive.NewObjectID().Hex(),
	}

	_, err := suite.mockRepo.Signup(context.Background(), newUser)
	suite.NoError(err)

	// Step 2: Retrieve the user by ID
	fetchedUser, err := suite.mockRepo.GetUser(context.Background(), newUser.UserId)
	suite.NoError(err)
	suite.NotNil(fetchedUser)
	suite.Equal(newUser.UserId, fetchedUser.UserId)
	suite.Equal(*newUser.FirstName, *fetchedUser.FirstName)
	suite.Equal(*newUser.LastName, *fetchedUser.LastName)
	suite.Equal(*newUser.Email, *fetchedUser.Email)
	suite.Equal(*newUser.Phone, *fetchedUser.Phone)
	suite.Equal(*newUser.UserType, *fetchedUser.UserType)

	// Step 3: Test case where no user is found
	fetchedUser, err = suite.mockRepo.GetUser(context.Background(), "nonexistent_id")
	suite.Error(err)
	suite.Equal(mongo.ErrNoDocuments, err)
	suite.Equal(domain.User{}, fetchedUser)
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail() {
	// Step 1: Add a new user to the database
	firstName := "Test"
	lastName := "User"
	password := "password123"
	email := "testuser@example.com"
	phone := "1234567890"
	userType := "USER"
	newUser := domain.User{
		ID:        primitive.NewObjectID(),
		FirstName: &firstName,
		LastName:  &lastName,
		Password:  &password,
		Email:     &email,
		Phone:     &phone,
		UserType:  &userType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserId:    primitive.NewObjectID().Hex(),
	}

	_, err := suite.mockRepo.Signup(context.Background(), newUser)
	suite.NoError(err)

	// Step 2: Retrieve the user by email
	fetchedUser, err := suite.mockRepo.GetUserByEmail(context.Background(), *newUser.Email)
	suite.NoError(err)
	suite.NotNil(fetchedUser)
	suite.Equal(newUser.UserId, fetchedUser.UserId)
	suite.Equal(*newUser.FirstName, *fetchedUser.FirstName)
	suite.Equal(*newUser.LastName, *fetchedUser.LastName)
	suite.Equal(*newUser.Email, *fetchedUser.Email)
	suite.Equal(*newUser.Phone, *fetchedUser.Phone)
	suite.Equal(*newUser.UserType, *fetchedUser.UserType)

	// Step 3: Test case where no user is found
	fetchedUser, err = suite.mockRepo.GetUserByEmail(context.Background(), "nonexistent@example.com")
	suite.Error(err)
	suite.Equal(mongo.ErrNoDocuments, err)
	suite.Equal(domain.User{}, fetchedUser)
}


func (suite *UserRepositoryTestSuite) TestGetUsers() {
    // Step 1: Add new users to the database
    firstName1 := "Test1"
    lastName1 := "User1"
    password1 := "password123"
    email1 := "testuser1@example.com"
    phone1 := "1234567890"
    userType1 := "USER"
    newUser1 := domain.User{
        ID:        primitive.NewObjectID(),
        FirstName: &firstName1,
        LastName:  &lastName1,
        Password:  &password1,
        Email:     &email1,
        Phone:     &phone1,
        UserType:  &userType1,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserId:    primitive.NewObjectID().Hex(),
    }

    firstName2 := "Test2"
    lastName2 := "User2"
    password2 := "password123"
    email2 := "testuser2@example.com"
    phone2 := "0987654321"
    userType2 := "USER"
    newUser2 := domain.User{
        ID:        primitive.NewObjectID(),
        FirstName: &firstName2,
        LastName:  &lastName2,
        Password:  &password2,
        Email:     &email2,
        Phone:     &phone2,
        UserType:  &userType2,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserId:    primitive.NewObjectID().Hex(),
    }

    _, err := suite.mockRepo.Signup(context.Background(), newUser1)
    suite.NoError(err)
    _, err = suite.mockRepo.Signup(context.Background(), newUser2)
    suite.NoError(err)

    // Step 2: Retrieve users with pagination
    users, err := suite.mockRepo.GetUsers(context.Background(), 0, 2)
    suite.NoError(err)
    suite.Len(users, 2)

    user1 := users[0]
    suite.Equal(newUser1.UserId, user1.UserId)
    suite.Equal(*newUser1.FirstName, *user1.FirstName)
    suite.Equal(*newUser1.LastName, *user1.LastName)
    suite.Equal(*newUser1.Email, *user1.Email)
    suite.Equal(*newUser1.Phone, *user1.Phone)
    suite.Equal(*newUser1.UserType, *user1.UserType)

    user2 := users[1]
    suite.Equal(newUser2.UserId, user2.UserId)
    suite.Equal(*newUser2.FirstName, *user2.FirstName)
    suite.Equal(*newUser2.LastName, *user2.LastName)
    suite.Equal(*newUser2.Email, *user2.Email)
    suite.Equal(*newUser2.Phone, *user2.Phone)
    suite.Equal(*newUser2.UserType, *user2.UserType)

    // Step 3: Test case where no users are found
    users, err = suite.mockRepo.GetUsers(context.Background(), 10, 2)
    suite.Error(err)
    suite.Nil(users)
}



func (suite *UserRepositoryTestSuite) TestLogin() {
    // Step 1: Add a new user to the database
    firstName := "Test"
    lastName := "User"
    password := "password123"
    email := "testuser@example.com"
    phone := "1234567890"
    userType := "USER"
    newUser := domain.User{
        ID:        primitive.NewObjectID(),
        FirstName: &firstName,
        LastName:  &lastName,
        Password:  &password,
        Email:     &email,
        Phone:     &phone,
        UserType:  &userType,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserId:    primitive.NewObjectID().Hex(),
    }

    _, err := suite.mockRepo.Signup(context.Background(), newUser)
    suite.NoError(err)

    // Step 2: Test successful login
    user, err := suite.mockRepo.Login(context.Background(), email)
    suite.NoError(err)
    suite.NotNil(user)
    suite.Equal(newUser.UserId, user.UserId)
    suite.Equal(*newUser.FirstName, *user.FirstName)
    suite.Equal(*newUser.LastName, *user.LastName)
    suite.Equal(*newUser.Email, *user.Email)
    suite.Equal(*newUser.Phone, *user.Phone)
    suite.Equal(*newUser.UserType, *user.UserType)

    // Step 3: Test login with non-existent email
    user, err = suite.mockRepo.Login(context.Background(), "nonexistent@example.com")
    suite.Error(err)
    suite.Nil(user)
    suite.EqualError(err, "user not found")
}


func (suite *UserRepositoryTestSuite) TestPromote() {
    // Step 1: Add a new user to the database
    firstName := "Test"
    lastName := "User"
    password := "password123"
    email := "testuser@example.com"
    phone := "1234567890"
    userType := "USER"
    newUser := domain.User{
        ID:        primitive.NewObjectID(),
        FirstName: &firstName,
        LastName:  &lastName,
        Password:  &password,
        Email:     &email,
        Phone:     &phone,
        UserType:  &userType,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserId:    primitive.NewObjectID().Hex(),
    }

    _, err := suite.mockRepo.Signup(context.Background(), newUser)
    suite.NoError(err)

    // Step 2: Promote the user to ADMIN
    err, matchedCount, modifiedCount := suite.mockRepo.Promote(context.Background(), newUser.UserId, "ADMIN")
    suite.NoError(err)
    suite.Equal(int64(1), matchedCount)
    suite.Equal(int64(1), modifiedCount)

    // Step 3: Verify the user type has been updated
    user, err := suite.mockRepo.Login(context.Background(), email)
    suite.NoError(err)
    suite.NotNil(user)
    suite.Equal("ADMIN", *user.UserType)

    // Step 4: Test promoting a non-existent user
    err, matchedCount, modifiedCount = suite.mockRepo.Promote(context.Background(), "nonexistent_user_id", "ADMIN")
    suite.Error(err)
    suite.Equal(int64(0), matchedCount)
    suite.Equal(int64(0), modifiedCount)
    suite.EqualError(err, "user not found")
}


func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
