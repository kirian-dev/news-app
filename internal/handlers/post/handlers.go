package post

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/kir/news-app/internal/domain"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for posts
type Handler struct {
	service   PostService
	templates *template.Template
	logger    *zap.Logger
}

// New creates a new post handler
func New(service PostService, templates *template.Template, logger *zap.Logger) *Handler {
	return &Handler{
		service:   service,
		templates: templates,
		logger:    logger,
	}
}

// handleError is a helper function to handle errors consistently
func (h *Handler) handleError(w http.ResponseWriter, err error, message string, status int) {
	h.logger.Error(message, zap.Error(err))
	w.Header().Set(HXErrorHeader, message)
	http.Error(w, message, status)
}

// handleHTMXSuccess is a helper function to handle successful HTMX responses
func (h *Handler) handleHTMXSuccess(w http.ResponseWriter, trigger string) {
	w.Header().Set(HXTriggerHeader, trigger)
	w.WriteHeader(http.StatusNoContent)
}

// Web handlers

// Index handles the main page request
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 9
	if sizeStr := r.URL.Query().Get("page_size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	search := r.URL.Query().Get("search")

	response, err := h.service.GetPaginated(ctx, page, pageSize, search)
	if err != nil {
		h.handleError(w, err, ErrFailedToLoadPosts, http.StatusInternalServerError)
		return
	}

	recentPosts, err := h.service.GetRecent(ctx, 5)
	if err != nil {
		h.logger.Error("failed to get recent posts", zap.Error(err))
	}

	totalPages := int(response.TotalCount) / pageSize
	if int(response.TotalCount)%pageSize > 0 {
		totalPages++
	}

	data := struct {
		Posts       []*domain.Post
		TotalCount  int64
		Page        int
		PageSize    int
		TotalPages  int
		Search      string
		RecentPosts []*domain.Post
	}{
		Posts:       response.Posts,
		TotalCount:  response.TotalCount,
		Page:        response.Page,
		PageSize:    response.PageSize,
		TotalPages:  totalPages,
		Search:      search,
		RecentPosts: recentPosts,
	}

	if r.Header.Get("HX-Request") == "true" {
		if err := h.templates.ExecuteTemplate(w, "post/posts-list", data); err != nil {
			h.handleError(w, err, ErrInternalServer, http.StatusInternalServerError)
		}
		return
	}

	if err := h.templates.ExecuteTemplate(w, "index", data); err != nil {
		h.handleError(w, err, ErrInternalServer, http.StatusInternalServerError)
	}
}

// CreateForm handles the post creation form request
func (h *Handler) CreateForm(w http.ResponseWriter, r *http.Request) {
	if err := h.templates.ExecuteTemplate(w, "modals/create", nil); err != nil {
		h.handleError(w, err, ErrInternalServer, http.StatusInternalServerError)
	}
}

// HTMX handlers

// View handles the post view request
func (h *Handler) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.handleError(w, nil, "Post ID is required", http.StatusBadRequest)
		return
	}

	post, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.handleError(w, err, ErrPostNotFound, http.StatusNotFound)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "modals/view-content", post); err != nil {
		h.handleError(w, err, "Error displaying the post", http.StatusInternalServerError)
	}
}

// EditForm handles the post edit form request
func (h *Handler) EditForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	h.logger.Info("editing post", zap.String("id", id))

	post, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.handleError(w, err, ErrPostNotFound, http.StatusNotFound)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "modals/edit-content", post); err != nil {
		h.handleError(w, err, "Error displaying edit form", http.StatusInternalServerError)
	}
}

// DeleteForm handles the post deletion form request
func (h *Handler) DeleteForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	h.logger.Info("deleting post", zap.String("id", id))

	post, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.handleError(w, err, ErrPostNotFound, http.StatusNotFound)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "modals/delete-content", post); err != nil {
		h.handleError(w, err, "Error displaying delete form", http.StatusInternalServerError)
	}
}

// Create handles the post creation request
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		h.handleError(w, err, ErrInvalidFormData, http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		h.handleError(w, nil, ErrEmptyFields, http.StatusBadRequest)
		return
	}

	h.logger.Info("creating post", zap.String("title", title))

	_, err := h.service.Create(ctx, title, content)
	if err != nil {
		h.logger.Error("failed to create post", zap.Error(err))
		h.handleError(w, err, err.Error(), http.StatusBadRequest)
		return
	}
	h.handleHTMXSuccess(w, TriggerPostCreated)
}

// Update handles the post update request
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	h.logger.Info("updating post", zap.String("id", id))

	if err := r.ParseForm(); err != nil {
		h.handleError(w, err, ErrInvalidFormData, http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		h.handleError(w, nil, ErrEmptyFields, http.StatusBadRequest)
		return
	}

	if err := h.service.Update(ctx, id, title, content); err != nil {
		h.handleError(w, err, err.Error(), http.StatusBadRequest)
		return
	}
	h.handleHTMXSuccess(w, TriggerPostUpdated)
}

// Delete handles the post deletion request
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	h.logger.Info("deleting post", zap.String("id", id))

	if err := h.service.Delete(ctx, id); err != nil {
		h.handleError(w, err, ErrFailedToDeletePost, http.StatusInternalServerError)
		return
	}
	h.handleHTMXSuccess(w, TriggerPostDeleted)
}
