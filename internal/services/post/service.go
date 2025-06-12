package post

import (
	"context"
	"fmt"
	"news-app/internal/domain"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, title, content string) (*domain.Post, error) {
	post, err := domain.NewPost(title, content)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	if err := s.repo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to save post: %w", err)
	}

	return post, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*domain.Post, error) {
	posts, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	return posts, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return post, nil
}

func (s *Service) Update(ctx context.Context, id, title, content string) error {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get post for update: %w", err)
	}

	if err := post.Update(title, content); err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	if err := s.repo.Update(ctx, post); err != nil {
		return fmt.Errorf("failed to save updated post: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

func (s *Service) GetPaginated(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 9
	}

	posts, err := s.repo.GetPaginated(ctx, page, pageSize, search)
	if err != nil {
		return nil, fmt.Errorf("failed to get paginated posts: %w", err)
	}
	return posts, nil
}

func (s *Service) GetRecent(ctx context.Context, limit int) ([]*domain.Post, error) {
	if limit < 1 {
		limit = 5
	}

	posts, err := s.repo.GetRecent(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent posts: %w", err)
	}
	return posts, nil
}
