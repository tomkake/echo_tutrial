package mocks

import (
	"context"
	"apiserver/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserInteractor struct {
	mock.Mock
}

func (m *MockUserInteractor) CreateNewUser(ctx context.Context, name, email, plainPassword string) (*domain.User, error) {
	args := m.Called(ctx, name, email, plainPassword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserInteractor) FindUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserInteractor) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserInteractor) UpdateExistingUser(ctx context.Context, id string, name, email *string, plainPassword *string) (*domain.User, error) {
	args := m.Called(ctx, id, name, email, plainPassword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserInteractor) RemoveUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
