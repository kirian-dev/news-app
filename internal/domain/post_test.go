package domain

import (
	"testing"
	"time"
)

func TestNewPost(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		content string
		wantErr bool
	}{
		{
			name:    "valid post",
			title:   "Valid Title",
			content: "This is a valid content with more than 10 characters",
			wantErr: false,
		},
		{
			name:    "short title",
			title:   "A",
			content: "Valid content",
			wantErr: true,
		},
		{
			name:    "long title",
			title:   string(make([]byte, 201)),
			content: "Valid content",
			wantErr: true,
		},
		{
			name:    "short content",
			title:   "Valid Title",
			content: "Short",
			wantErr: true,
		},
		{
			name:    "empty title",
			title:   "",
			content: "Valid content",
			wantErr: true,
		},
		{
			name:    "empty content",
			title:   "Valid Title",
			content: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := NewPost(tt.title, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if post.Title != tt.title {
					t.Errorf("NewPost() title = %v, want %v", post.Title, tt.title)
				}
				if post.Content != tt.content {
					t.Errorf("NewPost() content = %v, want %v", post.Content, tt.content)
				}
				if post.CreatedAt.IsZero() {
					t.Error("NewPost() CreatedAt is zero")
				}
				if post.UpdatedAt.IsZero() {
					t.Error("NewPost() UpdatedAt is zero")
				}
			}
		})
	}
}

func TestPost_Update(t *testing.T) {
	post, err := NewPost("Original Title", "Original content with more than 10 characters")
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name      string
		title     string
		content   string
		wantErr   bool
		checkTime bool
	}{
		{
			name:      "valid update",
			title:     "Updated Title",
			content:   "Updated content with more than 10 characters",
			wantErr:   false,
			checkTime: true,
		},
		{
			name:    "invalid title",
			title:   "A",
			content: "Valid content",
			wantErr: true,
		},
		{
			name:    "invalid content",
			title:   "Valid Title",
			content: "Short",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := post.UpdatedAt
			err := post.Update(tt.title, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if post.Title != tt.title {
					t.Errorf("Post.Update() title = %v, want %v", post.Title, tt.title)
				}
				if post.Content != tt.content {
					t.Errorf("Post.Update() content = %v, want %v", post.Content, tt.content)
				}
				if tt.checkTime && !post.UpdatedAt.After(oldUpdatedAt) {
					t.Error("Post.Update() UpdatedAt was not updated")
				}
			}
		})
	}
}

func TestPost_Validate(t *testing.T) {
	post := &Post{
		Title:     "Test Title",
		Content:   "Test content with more than 10 characters",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name    string
		post    *Post
		wantErr bool
	}{
		{
			name:    "valid post",
			post:    post,
			wantErr: false,
		},
		{
			name: "invalid title",
			post: &Post{
				Title:     "A",
				Content:   post.Content,
				CreatedAt: post.CreatedAt,
				UpdatedAt: post.UpdatedAt,
			},
			wantErr: true,
		},
		{
			name: "invalid content",
			post: &Post{
				Title:     post.Title,
				Content:   "Short",
				CreatedAt: post.CreatedAt,
				UpdatedAt: post.UpdatedAt,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.post.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Post.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
