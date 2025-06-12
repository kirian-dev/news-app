package post

import (
	"context"
	"errors"
	"testing"
	"time"

	"news-app/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestService_Create(t *testing.T) {
	tests := []struct {
		name          string
		title         string
		content       string
		mockCreate    func(ctx context.Context, p *domain.Post) error
		expectedError bool
		expectedPost  *domain.Post
	}{
		{
			name:    "successful creation",
			title:   "Test Post",
			content: "Test content with more than 10 characters",
			mockCreate: func(ctx context.Context, p *domain.Post) error {
				return nil
			},
			expectedError: false,
			expectedPost: &domain.Post{
				Title:   "Test Post",
				Content: "Test content with more than 10 characters",
			},
		},
		{
			name:    "empty title",
			title:   "",
			content: "Test content with more than 10 characters",
			mockCreate: func(ctx context.Context, p *domain.Post) error {
				return nil
			},
			expectedError: true,
		},
		{
			name:    "empty content",
			title:   "Test Post",
			content: "",
			mockCreate: func(ctx context.Context, p *domain.Post) error {
				return nil
			},
			expectedError: true,
		},
		{
			name:    "content too short",
			title:   "Test Post",
			content: "Short",
			mockCreate: func(ctx context.Context, p *domain.Post) error {
				return nil
			},
			expectedError: true,
		},
		{
			name:    "repository error",
			title:   "Test Post",
			content: "Test content with more than 10 characters",
			mockCreate: func(ctx context.Context, p *domain.Post) error {
				return errors.New("repository error")
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{
				CreateFunc: tt.mockCreate,
			}
			service := NewService(repo)

			post, err := service.Create(context.Background(), tt.title, tt.content)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, post)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, post)
				assert.Equal(t, tt.expectedPost.Title, post.Title)
				assert.Equal(t, tt.expectedPost.Content, post.Content)
				assert.NotEmpty(t, post.ID)
				assert.NotZero(t, post.CreatedAt)
				assert.NotZero(t, post.UpdatedAt)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	existingID := primitive.NewObjectID()
	tests := []struct {
		name          string
		id            string
		title         string
		content       string
		mockGetByID   func(ctx context.Context, id string) (*domain.Post, error)
		mockUpdate    func(ctx context.Context, p *domain.Post) error
		expectedError bool
	}{
		{
			name:    "successful update",
			id:      existingID.Hex(),
			title:   "Updated Title",
			content: "Updated content with more than 10 characters",
			mockGetByID: func(ctx context.Context, id string) (*domain.Post, error) {
				return &domain.Post{
					ID:        existingID,
					Title:     "Original Title",
					Content:   "Original content",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			mockUpdate: func(ctx context.Context, p *domain.Post) error {
				return nil
			},
			expectedError: false,
		},
		{
			name:    "post not found",
			id:      "nonexistent",
			title:   "Updated Title",
			content: "Updated content",
			mockGetByID: func(ctx context.Context, id string) (*domain.Post, error) {
				return nil, errors.New("post not found")
			},
			mockUpdate:    nil,
			expectedError: true,
		},
		{
			name:    "invalid update data",
			id:      existingID.Hex(),
			title:   "",
			content: "Updated content",
			mockGetByID: func(ctx context.Context, id string) (*domain.Post, error) {
				return &domain.Post{
					ID:        existingID,
					Title:     "Original Title",
					Content:   "Original content",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			mockUpdate:    nil,
			expectedError: true,
		},
		{
			name:    "update error",
			id:      existingID.Hex(),
			title:   "Updated Title",
			content: "Updated content with more than 10 characters",
			mockGetByID: func(ctx context.Context, id string) (*domain.Post, error) {
				return &domain.Post{
					ID:        existingID,
					Title:     "Original Title",
					Content:   "Original content",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			mockUpdate: func(ctx context.Context, p *domain.Post) error {
				return errors.New("update error")
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{
				GetByIDFunc: tt.mockGetByID,
				UpdateFunc:  tt.mockUpdate,
			}
			service := NewService(repo)

			err := service.Update(context.Background(), tt.id, tt.title, tt.content)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetPaginated(t *testing.T) {
	postID := primitive.NewObjectID()
	tests := []struct {
		name             string
		page             int
		pageSize         int
		search           string
		mockGetPaginated func(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error)
		expectedError    bool
		expectedList     *domain.PostList
	}{
		{
			name:     "successful pagination",
			page:     1,
			pageSize: 10,
			search:   "test",
			mockGetPaginated: func(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error) {
				return &domain.PostList{
					Posts: []*domain.Post{
						{
							ID:        postID,
							Title:     "Test Post",
							Content:   "Test content",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
					},
					TotalCount: 1,
					Page:       page,
					PageSize:   pageSize,
				}, nil
			},
			expectedError: false,
			expectedList: &domain.PostList{
				Posts: []*domain.Post{
					{
						ID:      postID,
						Title:   "Test Post",
						Content: "Test content",
					},
				},
				TotalCount: 1,
				Page:       1,
				PageSize:   10,
			},
		},
		{
			name:     "invalid page",
			page:     0,
			pageSize: 10,
			mockGetPaginated: func(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error) {
				return &domain.PostList{
					Page:     1,
					PageSize: pageSize,
				}, nil
			},
			expectedError: false,
			expectedList: &domain.PostList{
				Page:     1,
				PageSize: 10,
			},
		},
		{
			name:     "invalid page size",
			page:     1,
			pageSize: 0,
			mockGetPaginated: func(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error) {
				return &domain.PostList{
					Page:     page,
					PageSize: 9,
				}, nil
			},
			expectedError: false,
			expectedList: &domain.PostList{
				Page:     1,
				PageSize: 9,
			},
		},
		{
			name:     "repository error",
			page:     1,
			pageSize: 10,
			mockGetPaginated: func(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error) {
				return nil, errors.New("repository error")
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{
				GetPaginatedFunc: tt.mockGetPaginated,
			}
			service := NewService(repo)

			list, err := service.GetPaginated(context.Background(), tt.page, tt.pageSize, tt.search)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, list)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, list)
				assert.Equal(t, tt.expectedList.Page, list.Page)
				assert.Equal(t, tt.expectedList.PageSize, list.PageSize)
				if tt.expectedList.Posts != nil {
					assert.Equal(t, len(tt.expectedList.Posts), len(list.Posts))
				}
			}
		})
	}
}

func TestService_GetAll(t *testing.T) {
	postID := primitive.NewObjectID()
	tests := []struct {
		name          string
		mockGetAll    func(ctx context.Context) ([]*domain.Post, error)
		expectedError bool
		expectedPosts []*domain.Post
	}{
		{
			name: "successful get all",
			mockGetAll: func(ctx context.Context) ([]*domain.Post, error) {
				return []*domain.Post{
					{
						ID:        postID,
						Title:     "Test Post",
						Content:   "Test content",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil
			},
			expectedError: false,
			expectedPosts: []*domain.Post{
				{
					ID:      postID,
					Title:   "Test Post",
					Content: "Test content",
				},
			},
		},
		{
			name: "empty list",
			mockGetAll: func(ctx context.Context) ([]*domain.Post, error) {
				return []*domain.Post{}, nil
			},
			expectedError: false,
			expectedPosts: []*domain.Post{},
		},
		{
			name: "repository error",
			mockGetAll: func(ctx context.Context) ([]*domain.Post, error) {
				return nil, errors.New("repository error")
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{
				GetAllFunc: tt.mockGetAll,
			}
			service := NewService(repo)

			posts, err := service.GetAll(context.Background())

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, posts)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, posts)
				assert.Equal(t, len(tt.expectedPosts), len(posts))
				if len(posts) > 0 {
					assert.Equal(t, tt.expectedPosts[0].Title, posts[0].Title)
					assert.Equal(t, tt.expectedPosts[0].Content, posts[0].Content)
				}
			}
		})
	}
}

func TestService_GetByID(t *testing.T) {
	postID := primitive.NewObjectID()
	tests := []struct {
		name          string
		id            string
		mockGetByID   func(ctx context.Context, id string) (*domain.Post, error)
		expectedError bool
		expectedPost  *domain.Post
	}{
		{
			name: "successful get by id",
			id:   postID.Hex(),
			mockGetByID: func(ctx context.Context, id string) (*domain.Post, error) {
				return &domain.Post{
					ID:        postID,
					Title:     "Test Post",
					Content:   "Test content",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			expectedError: false,
			expectedPost: &domain.Post{
				ID:      postID,
				Title:   "Test Post",
				Content: "Test content",
			},
		},
		{
			name: "post not found",
			id:   "nonexistent",
			mockGetByID: func(ctx context.Context, id string) (*domain.Post, error) {
				return nil, errors.New("post not found")
			},
			expectedError: true,
		},
		{
			name: "invalid id format",
			id:   "invalid",
			mockGetByID: func(ctx context.Context, id string) (*domain.Post, error) {
				return nil, errors.New("invalid id format")
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{
				GetByIDFunc: tt.mockGetByID,
			}
			service := NewService(repo)

			post, err := service.GetByID(context.Background(), tt.id)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, post)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, post)
				assert.Equal(t, tt.expectedPost.Title, post.Title)
				assert.Equal(t, tt.expectedPost.Content, post.Content)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockDelete    func(ctx context.Context, id string) error
		expectedError bool
	}{
		{
			name: "successful delete",
			id:   primitive.NewObjectID().Hex(),
			mockDelete: func(ctx context.Context, id string) error {
				return nil
			},
			expectedError: false,
		},
		{
			name: "post not found",
			id:   "nonexistent",
			mockDelete: func(ctx context.Context, id string) error {
				return errors.New("post not found")
			},
			expectedError: true,
		},
		{
			name: "invalid id format",
			id:   "invalid",
			mockDelete: func(ctx context.Context, id string) error {
				return errors.New("invalid id format")
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{
				DeleteFunc: tt.mockDelete,
			}
			service := NewService(repo)

			err := service.Delete(context.Background(), tt.id)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetRecent(t *testing.T) {
	postID := primitive.NewObjectID()
	tests := []struct {
		name          string
		limit         int
		mockGetRecent func(ctx context.Context, limit int) ([]*domain.Post, error)
		expectedError bool
		expectedPosts []*domain.Post
	}{
		{
			name:  "successful get recent",
			limit: 5,
			mockGetRecent: func(ctx context.Context, limit int) ([]*domain.Post, error) {
				return []*domain.Post{
					{
						ID:        postID,
						Title:     "Test Post",
						Content:   "Test content",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil
			},
			expectedError: false,
			expectedPosts: []*domain.Post{
				{
					ID:      postID,
					Title:   "Test Post",
					Content: "Test content",
				},
			},
		},
		{
			name:  "invalid limit",
			limit: 0,
			mockGetRecent: func(ctx context.Context, limit int) ([]*domain.Post, error) {
				return []*domain.Post{}, nil
			},
			expectedError: false,
			expectedPosts: []*domain.Post{},
		},
		{
			name:  "repository error",
			limit: 5,
			mockGetRecent: func(ctx context.Context, limit int) ([]*domain.Post, error) {
				return nil, errors.New("repository error")
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{
				GetRecentFunc: tt.mockGetRecent,
			}
			service := NewService(repo)

			posts, err := service.GetRecent(context.Background(), tt.limit)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, posts)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, posts)
				assert.Equal(t, len(tt.expectedPosts), len(posts))
				if len(posts) > 0 {
					assert.Equal(t, tt.expectedPosts[0].Title, posts[0].Title)
					assert.Equal(t, tt.expectedPosts[0].Content, posts[0].Content)
				}
			}
		})
	}
}
