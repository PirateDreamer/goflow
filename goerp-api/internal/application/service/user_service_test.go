package service_test

import (
	"context"
	"errors"
	"goerp-api/internal/application/service"
	"goerp-api/internal/domain/derrors"
	"goerp-api/internal/domain/entity"
	repoMocks "goerp-api/internal/domain/repository/mocks"
	cacheMocks "goerp-api/internal/infrastructure/cache/mocks"
	emailMocks "goerp-api/internal/infrastructure/email/mocks"
	"testing"
	"time"
)

func TestUserService_LoginByEmailCode(t *testing.T) {
	mockRepo := &repoMocks.MockUserRepository{}
	mockCache := &cacheMocks.MockCache{}
	mockEmail := &emailMocks.MockEmailService{}
	svc := service.NewUserService(mockRepo, mockCache, mockEmail)

	ctx := context.Background()
	emailAddr := "test@example.com"
	code := "123456"

	t.Run("success login", func(t *testing.T) {
		mockCache.GetFunc = func(ctx context.Context, key string) (string, error) {
			return code, nil
		}
		mockCache.DeleteFunc = func(ctx context.Context, key string) error {
			return nil
		}
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{Email: emailAddr}, nil
		}

		user, err := svc.LoginByEmailCode(ctx, emailAddr, code)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user.Email != emailAddr {
			t.Errorf("expected email %s, got %s", emailAddr, user.Email)
		}
	})

	t.Run("invalid code", func(t *testing.T) {
		mockCache.GetFunc = func(ctx context.Context, key string) (string, error) {
			return "wrongcode", nil
		}

		_, err := svc.LoginByEmailCode(ctx, emailAddr, code)
		if err == nil || err.Error() != derrors.ErrInvalidVerification.Error() {
			t.Errorf("expected %v, got %v", derrors.ErrInvalidVerification, err)
		}
	})

	t.Run("expired code", func(t *testing.T) {
		mockCache.GetFunc = func(ctx context.Context, key string) (string, error) {
			return "", errors.New("not found")
		}

		_, err := svc.LoginByEmailCode(ctx, emailAddr, code)
		if err == nil || err.Error() != derrors.ErrVerificationExpired.Error() {
			t.Errorf("expected %v, got %v", derrors.ErrVerificationExpired, err)
		}
	})

	t.Run("auto register new user", func(t *testing.T) {
		mockCache.GetFunc = func(ctx context.Context, key string) (string, error) {
			return code, nil
		}
		mockCache.DeleteFunc = func(ctx context.Context, key string) error {
			return nil
		}
		// FindByEmail 返回错误，模拟用户不存在
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*entity.User, error) {
			return nil, errors.New("not found")
		}
		// Create 自动注册新用户
		mockRepo.CreateFunc = func(ctx context.Context, user *entity.User) error {
			// 验证自动生成的用户名是邮箱前缀
			if user.Username != "test" {
				t.Errorf("expected username 'test', got '%s'", user.Username)
			}
			if user.Email != emailAddr {
				t.Errorf("expected email %s, got %s", emailAddr, user.Email)
			}
			return nil
		}

		user, err := svc.LoginByEmailCode(ctx, emailAddr, code)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user == nil || user.Email != emailAddr {
			t.Errorf("expected user with email %s", emailAddr)
		}
	})

	t.Run("auto register fails", func(t *testing.T) {
		mockCache.GetFunc = func(ctx context.Context, key string) (string, error) {
			return code, nil
		}
		mockCache.DeleteFunc = func(ctx context.Context, key string) error {
			return nil
		}
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*entity.User, error) {
			return nil, errors.New("not found")
		}
		mockRepo.CreateFunc = func(ctx context.Context, user *entity.User) error {
			return errors.New("db error")
		}

		_, err := svc.LoginByEmailCode(ctx, emailAddr, code)
		if err == nil {
			t.Error("expected error when create fails, got nil")
		}
	})
}

func TestUserService_SendEmailVerificationCode(t *testing.T) {
	mockRepo := &repoMocks.MockUserRepository{}
	mockCache := &cacheMocks.MockCache{}
	mockEmail := &emailMocks.MockEmailService{}
	svc := service.NewUserService(mockRepo, mockCache, mockEmail)

	ctx := context.Background()
	emailAddr := "test@example.com"

	t.Run("success send", func(t *testing.T) {
		mockCache.SetFunc = func(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
			return nil
		}
		mockEmail.SendCodeFunc = func(to, code string) error {
			return nil
		}

		err := svc.SendEmailVerificationCode(ctx, emailAddr)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}
