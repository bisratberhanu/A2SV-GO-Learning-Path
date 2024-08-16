package usecases

import (
	"context"
	"task_manger_clean_architecture/domain"
	"task_manger_clean_architecture/domain/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestUserUseCase is the suite struct for testing UserUseCase.
type TestUserUseCase struct {
	suite.Suite
	UserUseCase domain.UserUseCase
	mockRepo    *mocks.UserRepository
}

// SetupTest sets up the necessary objects before each test.
func (c *TestUserUseCase) SetupTest() {
	c.mockRepo = new(mocks.UserRepository)
	c.UserUseCase = NewUserUseCase(c.mockRepo, time.Second*2)
}

// TestGetUser tests the GetUser method of the UserUseCase.
func (c *TestUserUseCase) TestGetUser() {
	userID := "1"
	user := domain.User{
		ID:        primitive.NewObjectID(),
		FirstName: stringPtr("John"),
		LastName:  stringPtr("Doe"),
		Email:     stringPtr("john.doe@example.com"),
	}

	c.mockRepo.On("GetUser", mock.Anything, userID).Return(user, nil)

	result, err := c.UserUseCase.GetUser(context.Background(), userID)
	assert.NoError(c.T(), err)
	assert.Equal(c.T(), user, result)

	c.mockRepo.AssertExpectations(c.T())
}

// TestGetUserByEmail tests the GetUserByEmail method of the UserUseCase.
func (c *TestUserUseCase) TestGetUserByEmail() {
	email := "john.doe@example.com"
	user := domain.User{
		ID:        primitive.NewObjectID(),
		FirstName: stringPtr("John"),
		LastName:  stringPtr("Doe"),
		Email:     stringPtr(email),
	}

	c.mockRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)

	result, err := c.UserUseCase.GetUserByEmail(context.Background(), email)
	assert.NoError(c.T(), err)
	assert.Equal(c.T(), user, result)

	c.mockRepo.AssertExpectations(c.T())
}

// TestGetUsers tests the GetUsers method of the UserUseCase.
func (c *TestUserUseCase) TestGetUsers() {
	users := []*domain.User{
		{
			ID:        primitive.NewObjectID(),
			FirstName: stringPtr("John"),
			LastName:  stringPtr("Doe"),
			Email:     stringPtr("john.doe@example.com"),
		},
		{
			ID:        primitive.NewObjectID(),
			FirstName: stringPtr("Jane"),
			LastName:  stringPtr("Doe"),
			Email:     stringPtr("jane.doe@example.com"),
		},
	}

	c.mockRepo.On("GetUsers", mock.Anything, int64(0), int64(10)).Return(users, nil)

	result, err := c.UserUseCase.GetUsers(context.Background(), 0, 10)
	assert.NoError(c.T(), err)
	assert.Equal(c.T(), users, result)

	c.mockRepo.AssertExpectations(c.T())
}

// TestLogin tests the Login method of the UserUseCase.
func (c *TestUserUseCase) TestLogin() {
	email := "john.doe@example.com"
	user := &domain.User{
		ID:        primitive.NewObjectID(),
		FirstName: stringPtr("John"),
		LastName:  stringPtr("Doe"),
		Email:     stringPtr(email),
	}

	c.mockRepo.On("Login", mock.Anything, email).Return(user, nil)

	result, err := c.UserUseCase.Login(context.Background(), email)
	assert.NoError(c.T(), err)
	assert.Equal(c.T(), user, result)

	c.mockRepo.AssertExpectations(c.T())
}

// TestPromote tests the Promote method of the UserUseCase.
func (c *TestUserUseCase) TestPromote() {
	userID := "1"
	userType := "ADMIN"

	c.mockRepo.On("Promote", mock.Anything, userID, userType).Return(nil, int64(1), int64(1))

	err, count1, count2 := c.UserUseCase.Promote(context.Background(), userID, userType)
	assert.NoError(c.T(), err)
	assert.Equal(c.T(), int64(1), count1)
	assert.Equal(c.T(), int64(1), count2)

	c.mockRepo.AssertExpectations(c.T())
}

// TestSignup tests the Signup method of the UserUseCase.
func (c *TestUserUseCase) TestSignup() {
	user := domain.User{
		FirstName: stringPtr("John"),
		LastName:  stringPtr("Doe"),
		Email:     stringPtr("john.doe@example.com"),
	}

	// Set up the mock expectation to match the modified user.
	c.mockRepo.On("Signup", mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(func(u domain.User) bool {
		// Validate the fields individually
		return u.FirstName != nil && *u.FirstName == "John" &&
			u.LastName != nil && *u.LastName == "Doe" &&
			u.Email != nil && *u.Email == "john.doe@example.com"
	})).Return("user_id", nil)

	// Call the Signup method.
	result, err := c.UserUseCase.Signup(context.Background(), user)

	// Check if no error occurred and the result matches the expected output.
	assert.NoError(c.T(), err)
	assert.Equal(c.T(), "user_id", result)

	// Verify that all expectations were met.
	c.mockRepo.AssertExpectations(c.T())
}


// TestUpdateAllTokens tests the UpdateAllTokens method of the UserUseCase.
func (c *TestUserUseCase) TestUpdateAllTokens() {
	token := "new_token"
	refreshToken := "new_refresh_token"
	userID := "1"

	c.mockRepo.On("UpdateAllTokens", token, refreshToken, userID).Return(nil)

	err := c.UserUseCase.UpdateAllTokens(token, refreshToken, userID)
	assert.NoError(c.T(), err)

	c.mockRepo.AssertExpectations(c.T())
}

// stringPtr is a helper function to create a string pointer.
func stringPtr(s string) *string {
	return &s
}

// TestUserUseCaseTestSuite runs the test suite.
func TestUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(TestUserUseCase))
}
