package postrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/kir/news-app/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository implements Repository interface using MongoDB
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository creates a new MongoDB repository
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("posts"),
	}
}

// Create implements Repository.Create
func (r *MongoRepository) Create(ctx context.Context, p *domain.Post) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("invalid post: %w", err)
	}

	res, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		return fmt.Errorf("failed to insert post: %w", err)
	}

	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("failed to convert InsertedID to ObjectID")
	}
	p.ID = objID
	return nil
}

// GetAll implements Repository.GetAll
func (r *MongoRepository) GetAll(ctx context.Context) ([]*domain.Post, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find posts: %w", err)
	}
	defer cursor.Close(ctx)

	var posts []*domain.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("failed to decode posts: %w", err)
	}
	return posts, nil
}

// GetByID implements Repository.GetByID
func (r *MongoRepository) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	var p domain.Post
	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&p); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to find post: %w", err)
	}
	return &p, nil
}

// Update implements Repository.Update
func (r *MongoRepository) Update(ctx context.Context, p *domain.Post) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("invalid post: %w", err)
	}

	p.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": p.ID},
		bson.M{
			"$set": bson.M{
				"title":      p.Title,
				"content":    p.Content,
				"updated_at": p.UpdatedAt,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("post not found")
	}
	return nil
}

// Delete implements Repository.Delete
func (r *MongoRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("post not found")
	}
	return nil
}

// GetPaginated implements Repository.GetPaginated
func (r *MongoRepository) GetPaginated(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 9
	}

	skip := (page - 1) * pageSize
	filter := bson.M{}

	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"title": bson.M{"$regex": search, "$options": "i"}},
				{"content": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %w", err)
	}

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find posts: %w", err)
	}
	defer cursor.Close(ctx)

	var posts []*domain.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("failed to decode posts: %w", err)
	}

	return &domain.PostList{
		Posts:      posts,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// GetRecent implements Repository.GetRecent
func (r *MongoRepository) GetRecent(ctx context.Context, limit int) ([]*domain.Post, error) {
	if limit < 1 {
		limit = 5
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find recent posts: %w", err)
	}
	defer cursor.Close(ctx)

	var posts []*domain.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("failed to decode recent posts: %w", err)
	}

	return posts, nil
}
