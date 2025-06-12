package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidTitle   = errors.New("title must be between 3 and 200 characters")
	ErrInvalidContent = errors.New("content must be at least 10 characters")
)

// Post represents a blog post with a title, content, and timestamps.
type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title" validate:"required,min=3,max=200"`
	Content   string             `bson:"content" json:"content" validate:"required,min=10"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// PostList is a paginated list of posts.
type PostList struct {
	Posts      []*Post `json:"posts"`
	TotalCount int64   `json:"total_count"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
}

// NewPost creates a new post with the given title and content.
// It returns an error if the title or content is invalid.
func NewPost(title, content string) (*Post, error) {
	if err := validatePostData(title, content); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Post{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Validate checks if the post's title and content meet the required constraints.
func (p *Post) Validate() error {
	return validatePostData(p.Title, p.Content)
}

// Update changes the post's title and content.
// It returns an error if the new data is invalid.
func (p *Post) Update(title, content string) error {
	if err := validatePostData(title, content); err != nil {
		return err
	}

	p.Title = title
	p.Content = content
	p.UpdatedAt = time.Now()
	return nil
}

// validatePostData validates the title and content according to length rules.
func validatePostData(title, content string) error {
	if len(title) < 3 || len(title) > 200 {
		return ErrInvalidTitle
	}
	if len(content) < 10 {
		return ErrInvalidContent
	}
	return nil
}
