package service

import (
	"context"
	"fmt"
	"goerp-api/internal/domain/derrors"
	"goerp-api/internal/domain/entity"
	"goerp-api/internal/domain/repository"
	"goerp-api/internal/infrastructure/cache"
	"goerp-api/internal/infrastructure/email"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo     repository.UserRepository
	cache    cache.Cache
	emailSvc email.EmailService
}

func NewUserService(repo repository.UserRepository, cache cache.Cache, emailSvc email.EmailService) *UserService {
	return &UserService{
		repo:     repo,
		cache:    cache,
		emailSvc: emailSvc,
	}
}

func (s *UserService) Register(ctx context.Context, username, email, password string) (*entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, derrors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, derrors.ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id uint) (*entity.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) SendEmailVerificationCode(ctx context.Context, emailAddr string) error {
	// 简单生成 6 位验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 存入缓存，5 分钟过期
	key := fmt.Sprintf("email_code:%s", emailAddr)
	if err := s.cache.Set(ctx, key, code, 5*time.Minute); err != nil {
		return err
	}

	// 发送邮件
	return s.emailSvc.SendCode(emailAddr, code)
}

func (s *UserService) LoginByEmailCode(ctx context.Context, emailAddr, code string) (*entity.User, error) {
	// 从缓存中获取验证码
	key := fmt.Sprintf("email_code:%s", emailAddr)
	val, err := s.cache.Get(ctx, key)
	if err != nil {
		return nil, derrors.ErrVerificationExpired
	}

	if val != code {
		return nil, derrors.ErrInvalidVerification
	}

	// 验证通过，删除验证码
	_ = s.cache.Delete(ctx, key)

	// 根据邮箱查找用户，若不存在则自动注册（无感注册）
	user, err := s.repo.FindByEmail(ctx, emailAddr)
	if err != nil {
		// 用户不存在，自动创建账号
		// 默认用户名取邮箱 @ 前的部分
		username := strings.SplitN(emailAddr, "@", 2)[0]
		user = &entity.User{
			Username: username,
			Email:    emailAddr,
			Password: "", // 邮箱验证码登录，无需密码
		}
		if createErr := s.repo.Create(ctx, user); createErr != nil {
			return nil, createErr
		}
	}

	return user, nil
}
