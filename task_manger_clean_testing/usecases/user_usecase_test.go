package usecases

import (
	"context"
	"task_manger_clean_architecture/domain"
	"task_manger_clean_architecture/domain/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUser(t *testing.T) {
    mockRepo := new(mocks.UserRepository)
    userID := "1"
    user := domain.User{
        ID:        primitive.NewObjectID(),
        FirstName: stringPtr("John"),
        LastName:  stringPtr("Doe"),
        Email:     stringPtr("john.doe@example.com"),
    }

    mockRepo.On("GetUser", mock.Anything, userID).Return(user, nil)

    result, err := mockRepo.GetUser(context.Background(), userID)
    assert.NoError(t, err)
    assert.Equal(t, user, result)

    mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail(t *testing.T) {
    mockRepo := new(mocks.UserRepository)
    email := "john.doe@example.com"
    user := domain.User{
        ID:        primitive.NewObjectID(),
        FirstName: stringPtr("John"),
        LastName:  stringPtr("Doe"),
        Email:     stringPtr(email),
    }

    mockRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)

    result, err := mockRepo.GetUserByEmail(context.Background(), email)
    assert.NoError(t, err)
    assert.Equal(t, user, result)

    mockRepo.AssertExpectations(t)
}

func TestGetUsers(t *testing.T) {
    mockRepo := new(mocks.UserRepository)
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

    mockRepo.On("GetUsers", mock.Anything, int64(0), int64(10)).Return(users, nil)

    result, err := mockRepo.GetUsers(context.Background(), 0, 10)
    assert.NoError(t, err)
    assert.Equal(t, users, result)

    mockRepo.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
    mockRepo := new(mocks.UserRepository)
    email := "john.doe@example.com"
    user := &domain.User{
        ID:        primitive.NewObjectID(),
        FirstName: stringPtr("John"),
        LastName:  stringPtr("Doe"),
        Email:     stringPtr(email),
    }

    mockRepo.On("Login", mock.Anything, email).Return(user, nil)

    result, err := mockRepo.Login(context.Background(), email)
    assert.NoError(t, err)
    assert.Equal(t, user, result)

    mockRepo.AssertExpectations(t)
}

func TestPromote(t *testing.T) {
    mockRepo := new(mocks.UserRepository)
    userID := "1"
    userType := "ADMIN"

    mockRepo.On("Promote", mock.Anything, userID, userType).Return(nil, int64(1), int64(1))

    err, count1, count2 := mockRepo.Promote(context.Background(), userID, userType)
    assert.NoError(t, err)
    assert.Equal(t, int64(1), count1)
    assert.Equal(t, int64(1), count2)

    mockRepo.AssertExpectations(t)
}

func TestSignup(t *testing.T) {
    mockRepo := new(mocks.UserRepository)
    user := domain.User{
        ID:        primitive.NewObjectID(),
        FirstName: stringPtr("John"),
        LastName:  stringPtr("Doe"),
        Email:     stringPtr("john.doe@example.com"),
    }

    mockRepo.On("Signup", mock.Anything, user).Return("user_id", nil)

    result, err := mockRepo.Signup(context.Background(), user)
    assert.NoError(t, err)
    assert.Equal(t, "user_id", result)

    mockRepo.AssertExpectations(t)
}

func TestUpdateAllTokens(t *testing.T) {
    mockRepo := new(mocks.UserRepository)
	token := "new_token"
    refreshToken := "new_refresh_token"
    userID := "1"

    mockRepo.On("UpdateAllTokens", token, refreshToken, userID).Return(nil)

    err := mockRepo.UpdateAllTokens(token, refreshToken, userID)
    assert.NoError(t, err)

    mockRepo.AssertExpectations(t)
}

func stringPtr(s string) *string {
    return &s
}