package domain

import (
	"context"
)

// Repository defines the interface for post storage operations
type Repository interface {
	Create(ctx context.Context, post *Post) error
	GetAll(ctx context.Context) ([]*Post, error)
	GetByID(ctx context.Context, id string) (*Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id string) error
	GetPaginated(ctx context.Context, page, pageSize int, search string) (*PostList, error)
	GetRecent(ctx context.Context, limit int) ([]*Post, error)
}
