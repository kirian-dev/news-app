package post

import (
	"context"

	"github.com/kir/news-app/internal/domain"
)

type PostService interface {
	Create(ctx context.Context, title, content string) (*domain.Post, error)
	GetAll(ctx context.Context) ([]*domain.Post, error)
	GetByID(ctx context.Context, id string) (*domain.Post, error)
	Update(ctx context.Context, id, title, content string) error
	Delete(ctx context.Context, id string) error
	GetPaginated(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error)
	GetRecent(ctx context.Context, limit int) ([]*domain.Post, error)
}
