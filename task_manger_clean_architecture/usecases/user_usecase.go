package usecases

import (
	"context"
	"fmt"
	"task_manger_clean_architecture/domain"
	"time"
)

type UserUseCase struct {
	UserRepository domain.UserRepository
	contextTimeout time.Duration
}



func NewUserUseCase(userRepository domain.UserRepository, timeout time.Duration) domain.UserUseCase {
	return &UserUseCase{
		UserRepository: userRepository,
		contextTimeout: timeout,
	}
}

// GetUser implements domain.UserUseCase.
func (u *UserUseCase) GetUser(c context.Context, user_id string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.UserRepository.GetUser(ctx, user_id)
}
func (u *UserUseCase) GetUserByEmail(c context.Context,email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.UserRepository.GetUserByEmail(ctx, email)
}

// GetUsers implements domain.UserUseCase.
func (u *UserUseCase) GetUsers(c context.Context, startIndex int64, recordsPerPage int64) ([]*domain.User, error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()
    
    users, err := u.UserRepository.GetUsers(ctx, startIndex, recordsPerPage)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve users: %v", err)
    }

    return users, nil
}


// Login implements domain.UserUseCase.
func (u *UserUseCase) Login(c context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.UserRepository.Login(ctx, email)
}

// Promote implements domain.UserUseCase.
func (u *UserUseCase) Promote(c context.Context, user_id string, userType string) (error,int64, int64) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.UserRepository.Promote(ctx, user_id, userType)
}

// Signup implements domain.UserUseCase.
func (u *UserUseCase) Signup(c context.Context, user domain.User) (interface{}, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.UserRepository.Signup(ctx, user)
}
// UpdateToken implements domain.UserUseCase.
func (u *UserUseCase) UpdateAllTokens(token string, refreshToken string, user_id string) error {
	return u.UserRepository.UpdateAllTokens(token, refreshToken, user_id)
}

