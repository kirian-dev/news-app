package postrepo

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"news-app/internal/domain"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testDB   *mongo.Database
	testRepo domain.Repository
	resource *dockertest.Resource
)

func TestMain(m *testing.M) {
	// Run MongoDB in Docker
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "6",
		Env: []string{
			"MONGO_INITDB_DATABASE=test",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		panic(err)
	}

	// Wait for MongoDB to be ready
	if err := pool.Retry(func() error {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:"+resource.GetPort("27017/tcp")))
		if err != nil {
			return err
		}
		return client.Ping(context.Background(), nil)
	}); err != nil {
		panic(err)
	}

	// Connect to test database
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:"+resource.GetPort("27017/tcp")))
	if err != nil {
		panic(err)
	}
	testDB = client.Database("test")
	testRepo = NewMongoRepository(testDB)

	// Run tests
	code := m.Run()

	// Clean up
	if err := pool.Purge(resource); err != nil {
		panic(err)
	}
	os.Exit(code)
}

func TestMongoRepository_Create(t *testing.T) {
	ctx := context.Background()
	post, err := domain.NewPost("Test Title", "Test content with more than 10 characters")
	require.NoError(t, err)

	// Test successful creation
	err = testRepo.Create(ctx, post)
	assert.NoError(t, err)
	assert.NotEmpty(t, post.ID)

	// Test creation of post with existing ID
	err = testRepo.Create(ctx, post)
	assert.Error(t, err)
}

func TestMongoRepository_GetByID(t *testing.T) {
	ctx := context.Background()
	post, err := domain.NewPost("Test Title", "Test content with more than 10 characters")
	require.NoError(t, err)

	// Create post for test
	err = testRepo.Create(ctx, post)
	require.NoError(t, err)

	// Test successful retrieval
	found, err := testRepo.GetByID(ctx, post.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, post.ID, found.ID)
	assert.Equal(t, post.Title, found.Title)
	assert.Equal(t, post.Content, found.Content)

	// Test retrieval of non-existent post
	_, err = testRepo.GetByID(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestMongoRepository_Update(t *testing.T) {
	ctx := context.Background()
	post, err := domain.NewPost("Original Title", "Original content with more than 10 characters")
	require.NoError(t, err)

	// Create post for test
	err = testRepo.Create(ctx, post)
	require.NoError(t, err)

	// Test successful update
	newTitle := "Updated Title"
	newContent := "Updated content with more than 10 characters"
	oldUpdatedAt := post.UpdatedAt
	post.Title = newTitle
	post.Content = newContent
	err = testRepo.Update(ctx, post)
	assert.NoError(t, err)

	// Check update
	updated, err := testRepo.GetByID(ctx, post.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, newTitle, updated.Title)
	assert.Equal(t, newContent, updated.Content)
	assert.True(t, updated.UpdatedAt.After(oldUpdatedAt))

	// Test update of non-existent post
	nonExistentPost := &domain.Post{
		ID:        primitive.NewObjectID(),
		Title:     "Non-existent Title",
		Content:   "Non-existent content with more than 10 characters",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = testRepo.Update(ctx, nonExistentPost)
	assert.Error(t, err)
}

func TestMongoRepository_Delete(t *testing.T) {
	ctx := context.Background()
	post, err := domain.NewPost("Test Title", "Test content with more than 10 characters")
	require.NoError(t, err)

	// Create post for test
	err = testRepo.Create(ctx, post)
	require.NoError(t, err)

	// Test successful deletion
	err = testRepo.Delete(ctx, post.ID.Hex())
	assert.NoError(t, err)

	// Check that post is deleted
	_, err = testRepo.GetByID(ctx, post.ID.Hex())
	assert.Error(t, err)

	// Test deletion of non-existent post
	err = testRepo.Delete(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestMongoRepository_GetPaginated(t *testing.T) {
	ctx := context.Background()

	// Clean up collection before test
	err := testDB.Collection("posts").Drop(ctx)
	require.NoError(t, err)

	// Create test posts
	for i := 0; i < 15; i++ {
		post, err := domain.NewPost(
			fmt.Sprintf("Test Title %d", i),
			fmt.Sprintf("Test content %d with more than 10 characters", i),
		)
		require.NoError(t, err)
		err = testRepo.Create(ctx, post)
		require.NoError(t, err)
	}

	// Test pagination
	tests := []struct {
		name     string
		page     int
		pageSize int
		search   string
		want     int
	}{
		{
			name:     "first page",
			page:     1,
			pageSize: 10,
			want:     10,
		},
		{
			name:     "second page",
			page:     2,
			pageSize: 10,
			want:     5,
		},
		{
			name:     "search by title",
			page:     1,
			pageSize: 10,
			search:   "Title 1",
			want:     6,
		},
		{
			name:     "search by content",
			page:     1,
			pageSize: 10,
			search:   "content 1",
			want:     6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := testRepo.GetPaginated(ctx, tt.page, tt.pageSize, tt.search)
			assert.NoError(t, err)
			assert.Len(t, result.Posts, tt.want)
		})
	}
}

func TestMongoRepository_GetRecent(t *testing.T) {
	ctx := context.Background()

	// Create test posts
	for i := 0; i < 10; i++ {
		post, err := domain.NewPost(
			fmt.Sprintf("Test Title %d", i),
			fmt.Sprintf("Test content %d with more than 10 characters", i),
		)
		require.NoError(t, err)
		err = testRepo.Create(ctx, post)
		require.NoError(t, err)
		// For different timestamps
		time.Sleep(time.Millisecond * 100)
	}

	// Test getting recent posts
	tests := []struct {
		name  string
		limit int
		want  int
	}{
		{
			name:  "get 5 recent posts",
			limit: 5,
			want:  5,
		},
		{
			name:  "get all recent posts",
			limit: 10,
			want:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts, err := testRepo.GetRecent(ctx, tt.limit)
			assert.NoError(t, err)
			assert.Len(t, posts, tt.want)

			// Check order (from newest to oldest)
			for i := 1; i < len(posts); i++ {
				assert.True(t, posts[i-1].CreatedAt.After(posts[i].CreatedAt))
			}
		})
	}
}

func TestMongoRepository_GetAll(t *testing.T) {
	ctx := context.Background()

	// Clean up collection before test
	err := testDB.Collection("posts").Drop(ctx)
	require.NoError(t, err)

	// Create test posts
	expectedPosts := make([]*domain.Post, 5)
	for i := 0; i < 5; i++ {
		post, err := domain.NewPost(
			fmt.Sprintf("Test Title %d", i),
			fmt.Sprintf("Test content %d with more than 10 characters", i),
		)
		require.NoError(t, err)
		err = testRepo.Create(ctx, post)
		require.NoError(t, err)
		expectedPosts[i] = post
	}

	// Test getting all posts
	posts, err := testRepo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, posts, len(expectedPosts))

	// Check that all posts are returned
	postMap := make(map[string]*domain.Post)
	for _, post := range posts {
		postMap[post.ID.Hex()] = post
	}

	for _, expectedPost := range expectedPosts {
		actualPost, exists := postMap[expectedPost.ID.Hex()]
		assert.True(t, exists)
		assert.Equal(t, expectedPost.Title, actualPost.Title)
		assert.Equal(t, expectedPost.Content, actualPost.Content)
	}
}

func TestMongoRepository_ContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	post, err := domain.NewPost("Test Title", "Test content with more than 10 characters")
	require.NoError(t, err)

	// Test all operations with timeout
	_, err = testRepo.GetAll(ctx)
	assert.Error(t, err)

	_, err = testRepo.GetByID(ctx, "test")
	assert.Error(t, err)

	err = testRepo.Create(ctx, post)
	assert.Error(t, err)

	err = testRepo.Update(ctx, post)
	assert.Error(t, err)

	err = testRepo.Delete(ctx, "test")
	assert.Error(t, err)

	_, err = testRepo.GetPaginated(ctx, 1, 10, "")
	assert.Error(t, err)

	_, err = testRepo.GetRecent(ctx, 5)
	assert.Error(t, err)
}
