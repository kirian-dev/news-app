package post

import (
	"context"
	"news-app/internal/domain"
)

// MockService implements PostService interface for testing
type MockService struct {
	CreateFunc       func(ctx context.Context, title, content string) (*domain.Post, error)
	GetAllFunc       func(ctx context.Context) ([]*domain.Post, error)
	GetByIDFunc      func(ctx context.Context, id string) (*domain.Post, error)
	UpdateFunc       func(ctx context.Context, id, title, content string) error
	DeleteFunc       func(ctx context.Context, id string) error
	GetPaginatedFunc func(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error)
	GetRecentFunc    func(ctx context.Context, limit int) ([]*domain.Post, error)
}

func (m *MockService) Create(ctx context.Context, title, content string) (*domain.Post, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, title, content)
	}
	return nil, nil
}

func (m *MockService) GetAll(ctx context.Context) ([]*domain.Post, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockService) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockService) Update(ctx context.Context, id, title, content string) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, title, content)
	}
	return nil
}

func (m *MockService) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockService) GetPaginated(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error) {
	if m.GetPaginatedFunc != nil {
		return m.GetPaginatedFunc(ctx, page, pageSize, search)
	}
	return nil, nil
}

func (m *MockService) GetRecent(ctx context.Context, limit int) ([]*domain.Post, error) {
	if m.GetRecentFunc != nil {
		return m.GetRecentFunc(ctx, limit)
	}
	return nil, nil
}
