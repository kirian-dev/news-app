package post

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kir/news-app/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHTMX_Responses(t *testing.T) {
	handler, mockService := setupTestHandler()
	postID := primitive.NewObjectID()

	tests := []struct {
		name           string
		request        func() *http.Request
		mockService    func()
		expectedStatus int
		expectedHeader string
		expectedValue  string
	}{
		{
			name: "create post success",
			request: func() *http.Request {
				form := bytes.NewBufferString("title=Test&content=Content")
				req := httptest.NewRequest(http.MethodPost, "/posts", form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.Header.Set("HX-Request", "true")
				return req
			},
			mockService: func() {
				mockService.CreateFunc = func(ctx context.Context, title, content string) (*domain.Post, error) {
					return &domain.Post{
						ID:      postID,
						Title:   title,
						Content: content,
					}, nil
				}
			},
			expectedStatus: http.StatusNoContent,
			expectedHeader: HXTriggerHeader,
			expectedValue:  TriggerPostCreated,
		},
		{
			name: "update post success",
			request: func() *http.Request {
				form := bytes.NewBufferString("title=Updated&content=Content")
				req := httptest.NewRequest(http.MethodPut, "/posts/"+postID.Hex(), form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.Header.Set("HX-Request", "true")
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", postID.Hex())
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
				return req
			},
			mockService: func() {
				mockService.UpdateFunc = func(ctx context.Context, id, title, content string) error {
					return nil
				}
			},
			expectedStatus: http.StatusNoContent,
			expectedHeader: HXTriggerHeader,
			expectedValue:  TriggerPostUpdated,
		},
		{
			name: "delete post success",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/posts/"+postID.Hex(), nil)
				req.Header.Set("HX-Request", "true")
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", postID.Hex())
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
				return req
			},
			mockService: func() {
				mockService.DeleteFunc = func(ctx context.Context, id string) error {
					return nil
				}
			},
			expectedStatus: http.StatusNoContent,
			expectedHeader: HXTriggerHeader,
			expectedValue:  TriggerPostDeleted,
		},
		{
			name: "view post success",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/posts/"+postID.Hex(), nil)
				req.Header.Set("HX-Request", "true")
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", postID.Hex())
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
				return req
			},
			mockService: func() {
				mockService.GetByIDFunc = func(ctx context.Context, id string) (*domain.Post, error) {
					return &domain.Post{
						ID:      postID,
						Title:   "Test Post",
						Content: "Test content",
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedHeader: "",
			expectedValue:  "",
		},
		{
			name: "create post error",
			request: func() *http.Request {
				form := bytes.NewBufferString("title=&content=")
				req := httptest.NewRequest(http.MethodPost, "/posts", form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.Header.Set("HX-Request", "true")
				return req
			},
			mockService: func() {
				mockService.CreateFunc = func(ctx context.Context, title, content string) (*domain.Post, error) {
					return nil, nil
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectedHeader: HXErrorHeader,
			expectedValue:  ErrEmptyFields,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()
			w := httptest.NewRecorder()
			req := tt.request()

			switch req.Method {
			case http.MethodPost:
				handler.Create(w, req)
			case http.MethodPut:
				handler.Update(w, req)
			case http.MethodDelete:
				handler.Delete(w, req)
			case http.MethodGet:
				handler.View(w, req)
			}

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedHeader != "" {
				assert.Equal(t, tt.expectedValue, w.Header().Get(tt.expectedHeader))
			}
		})
	}
}

func TestHTMX_NonHTMXRequests(t *testing.T) {
	handler, mockService := setupTestHandler()
	postID := primitive.NewObjectID()

	tests := []struct {
		name           string
		request        func() *http.Request
		mockService    func()
		expectedStatus int
	}{
		{
			name: "create post without HTMX",
			request: func() *http.Request {
				form := bytes.NewBufferString("title=Test&content=Content")
				req := httptest.NewRequest(http.MethodPost, "/posts", form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			},
			mockService: func() {
				mockService.CreateFunc = func(ctx context.Context, title, content string) (*domain.Post, error) {
					return &domain.Post{
						ID:      postID,
						Title:   title,
						Content: content,
					}, nil
				}
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "update post without HTMX",
			request: func() *http.Request {
				form := bytes.NewBufferString("title=Updated&content=Content")
				req := httptest.NewRequest(http.MethodPut, "/posts/"+postID.Hex(), form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", postID.Hex())
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
				return req
			},
			mockService: func() {
				mockService.UpdateFunc = func(ctx context.Context, id, title, content string) error {
					return nil
				}
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()
			w := httptest.NewRecorder()
			req := tt.request()

			switch req.Method {
			case http.MethodPost:
				handler.Create(w, req)
			case http.MethodPut:
				handler.Update(w, req)
			}

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
