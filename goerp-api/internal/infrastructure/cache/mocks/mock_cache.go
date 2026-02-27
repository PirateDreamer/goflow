package mocks

import (
	"context"
	"time"
)

type MockCache struct {
	SetFunc    func(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetFunc    func(ctx context.Context, key string) (string, error)
	DeleteFunc func(ctx context.Context, key string) error
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.SetFunc(ctx, key, value, expiration)
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	return m.GetFunc(ctx, key)
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	return m.DeleteFunc(ctx, key)
}
