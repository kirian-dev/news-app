package post

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"news-app/internal/domain"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func setupTestHandler() (*Handler, *MockService) {
	mockService := &MockService{}
	tmpl := template.New("")
	tmpl.Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"multiply": func(a, b int) int {
			return a * b
		},
		"sequence": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
		"objectIDToString": func(id primitive.ObjectID) string {
			return id.Hex()
		},
	})
	tmpl = template.Must(tmpl.ParseGlob("../../../templates/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("../../../templates/post/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("../../../templates/modals/*.html"))
	logger, _ := zap.NewDevelopment()
	handler := New(mockService, tmpl, logger)
	return handler, mockService
}

func TestHandler_Create(t *testing.T) {
	handler, mockService := setupTestHandler()

	tests := []struct {
		name           string
		title          string
		content        string
		mockCreate     func(ctx context.Context, title, content string) (*domain.Post, error)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:    "successful creation",
			title:   "Test Post",
			content: "Test content",
			mockCreate: func(ctx context.Context, title, content string) (*domain.Post, error) {
				return &domain.Post{
					ID:      primitive.NewObjectID(),
					Title:   title,
					Content: content,
				}, nil
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  false,
		},
		{
			name:    "empty title",
			title:   "",
			content: "Test content",
			mockCreate: func(ctx context.Context, title, content string) (*domain.Post, error) {
				return nil, nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:    "empty content",
			title:   "Test Post",
			content: "",
			mockCreate: func(ctx context.Context, title, content string) (*domain.Post, error) {
				return nil, nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:    "service error",
			title:   "Test Post",
			content: "Test content",
			mockCreate: func(ctx context.Context, title, content string) (*domain.Post, error) {
				return nil, assert.AnError
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.CreateFunc = tt.mockCreate

			form := bytes.NewBufferString("")
			form.WriteString("title=" + tt.title)
			form.WriteString("&content=" + tt.content)

			req := httptest.NewRequest(http.MethodPost, "/posts", form)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			handler.Create(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				assert.Contains(t, w.Header().Get(HXErrorHeader), "")
			} else {
				assert.Equal(t, TriggerPostCreated, w.Header().Get(HXTriggerHeader))
			}
		})
	}
}

func TestHandler_Update(t *testing.T) {
	handler, mockService := setupTestHandler()
	postID := primitive.NewObjectID()

	tests := []struct {
		name           string
		id             string
		title          string
		content        string
		mockUpdate     func(ctx context.Context, id, title, content string) error
		expectedStatus int
		expectedError  bool
	}{
		{
			name:    "successful update",
			id:      postID.Hex(),
			title:   "Updated Title",
			content: "Updated content",
			mockUpdate: func(ctx context.Context, id, title, content string) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  false,
		},
		{
			name:    "empty title",
			id:      postID.Hex(),
			title:   "",
			content: "Updated content",
			mockUpdate: func(ctx context.Context, id, title, content string) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:    "empty content",
			id:      postID.Hex(),
			title:   "Updated Title",
			content: "",
			mockUpdate: func(ctx context.Context, id, title, content string) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:    "service error",
			id:      postID.Hex(),
			title:   "Updated Title",
			content: "Updated content",
			mockUpdate: func(ctx context.Context, id, title, content string) error {
				return assert.AnError
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.UpdateFunc = tt.mockUpdate

			form := bytes.NewBufferString("")
			form.WriteString("title=" + tt.title)
			form.WriteString("&content=" + tt.content)

			req := httptest.NewRequest(http.MethodPut, "/posts/"+tt.id, form)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			w := httptest.NewRecorder()

			handler.Update(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				assert.Contains(t, w.Header().Get(HXErrorHeader), "")
			} else {
				assert.Equal(t, TriggerPostUpdated, w.Header().Get(HXTriggerHeader))
			}
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	handler, mockService := setupTestHandler()
	postID := primitive.NewObjectID()

	tests := []struct {
		name           string
		id             string
		mockDelete     func(ctx context.Context, id string) error
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful delete",
			id:   postID.Hex(),
			mockDelete: func(ctx context.Context, id string) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  false,
		},
		{
			name: "service error",
			id:   postID.Hex(),
			mockDelete: func(ctx context.Context, id string) error {
				return assert.AnError
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.DeleteFunc = tt.mockDelete

			req := httptest.NewRequest(http.MethodDelete, "/posts/"+tt.id, nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			w := httptest.NewRecorder()

			handler.Delete(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				assert.Contains(t, w.Header().Get(HXErrorHeader), "")
			} else {
				assert.Equal(t, TriggerPostDeleted, w.Header().Get(HXTriggerHeader))
			}
		})
	}
}

func TestHandler_GetByID(t *testing.T) {
	handler, mockService := setupTestHandler()
	postID := primitive.NewObjectID()

	tests := []struct {
		name           string
		id             string
		mockGetByID    func(ctx context.Context, id string) (*domain.Post, error)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful get",
			id:   postID.Hex(),
			mockGetByID: func(ctx context.Context, id string) (*domain.Post, error) {
				return &domain.Post{
					ID:      postID,
					Title:   "Test Post",
					Content: "Test content",
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "post not found",
			id:   "nonexistent",
			mockGetByID: func(ctx context.Context, id string) (*domain.Post, error) {
				return nil, assert.AnError
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.GetByIDFunc = tt.mockGetByID

			req := httptest.NewRequest(http.MethodGet, "/posts/"+tt.id, nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			w := httptest.NewRecorder()

			handler.View(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				assert.Contains(t, w.Header().Get(HXErrorHeader), "")
			}
		})
	}
}

func TestHandler_GetPaginated(t *testing.T) {
	handler, mockService := setupTestHandler()

	tests := []struct {
		name             string
		page             int
		pageSize         int
		search           string
		mockGetPaginated func(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error)
		expectedStatus   int
		expectedError    bool
	}{
		{
			name:     "successful get paginated",
			page:     1,
			pageSize: 10,
			search:   "",
			mockGetPaginated: func(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error) {
				return &domain.PostList{
					Posts: []*domain.Post{
						{
							ID:      primitive.NewObjectID(),
							Title:   "Test Post",
							Content: "Test content",
						},
					},
					TotalCount: 1,
					Page:       1,
					PageSize:   10,
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:     "service error",
			page:     1,
			pageSize: 10,
			search:   "",
			mockGetPaginated: func(ctx context.Context, page, pageSize int, search string) (*domain.PostList, error) {
				return nil, assert.AnError
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.GetPaginatedFunc = tt.mockGetPaginated

			req := httptest.NewRequest(http.MethodGet, "/?page="+strconv.Itoa(tt.page)+"&page_size="+strconv.Itoa(tt.pageSize)+"&search="+tt.search, nil)
			w := httptest.NewRecorder()

			handler.Index(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				assert.Contains(t, w.Header().Get(HXErrorHeader), "")
			}
		})
	}
}
