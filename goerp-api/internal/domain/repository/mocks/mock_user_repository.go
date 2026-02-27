package mocks

import (
	"context"
	"goerp-api/internal/domain/entity"
)

type MockUserRepository struct {
	CreateFunc         func(ctx context.Context, user *entity.User) error
	FindByIDFunc       func(ctx context.Context, id uint) (*entity.User, error)
	FindByUsernameFunc func(ctx context.Context, username string) (*entity.User, error)
	FindByEmailFunc    func(ctx context.Context, email string) (*entity.User, error)
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	return m.CreateFunc(ctx, user)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	return m.FindByIDFunc(ctx, id)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	return m.FindByUsernameFunc(ctx, username)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return m.FindByEmailFunc(ctx, email)
}
