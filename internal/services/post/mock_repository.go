package post

import (
	"context"

	"github.com/kir/news-app/internal/domain"
)

// MockRepository is a mock implementation of domain.Repository
type MockRepository struct {
	CreateFunc       func(ctx context.Context, post *domain.Post) error
	GetAllFunc       func(ctx context.Context) ([]*domain.Post, error)
	GetByIDFunc      func(ctx context.Context, id string) (*domain.Post, error)
	UpdateFunc       func(ctx context.Context, post *domain.Post) error
	DeleteFunc       func(ctx context.Context, id string) error
	GetPaginatedFunc func(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error)
	GetRecentFunc    func(ctx context.Context, limit int) ([]*domain.Post, error)
}

func (m *MockRepository) Create(ctx context.Context, post *domain.Post) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, post)
	}
	return nil
}

func (m *MockRepository) GetAll(ctx context.Context) ([]*domain.Post, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockRepository) Update(ctx context.Context, post *domain.Post) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, post)
	}
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockRepository) GetPaginated(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error) {
	if m.GetPaginatedFunc != nil {
		return m.GetPaginatedFunc(ctx, page, pageSize, search)
	}
	return nil, nil
}

func (m *MockRepository) GetRecent(ctx context.Context, limit int) ([]*domain.Post, error) {
	if m.GetRecentFunc != nil {
		return m.GetRecentFunc(ctx, limit)
	}
	return nil, nil
}
