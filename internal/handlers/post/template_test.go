package post

import (
	"bytes"
	"html/template"
	"news-app/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTemplates_Render(t *testing.T) {
	tmpl := template.New("")
	tmpl.Funcs(template.FuncMap{
		"add":      func(a, b int) int { return a + b },
		"subtract": func(a, b int) int { return a - b },
		"multiply": func(a, b int) int { return a * b },
		"sequence": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
		"objectIDToString": func(id primitive.ObjectID) string { return id.Hex() },
	})
	tmpl = template.Must(tmpl.ParseGlob("../../../templates/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("../../../templates/post/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("../../../templates/modals/*.html"))

	t.Run("index_template_with_posts", func(t *testing.T) {
		var buf bytes.Buffer
		posts := []*domain.Post{{Title: "Test Post", Content: "Test content"}}
		data := struct {
			Posts       []*domain.Post
			TotalCount  int64
			Page        int
			PageSize    int
			TotalPages  int
			Search      string
			RecentPosts []*domain.Post
		}{
			Posts:       posts,
			TotalCount:  1,
			Page:        1,
			PageSize:    10,
			TotalPages:  1,
			Search:      "",
			RecentPosts: posts,
		}
		err := tmpl.ExecuteTemplate(&buf, "index", data)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Test Post")
		assert.Contains(t, buf.String(), "Test content")
	})

	t.Run("create_form_template", func(t *testing.T) {
		var buf bytes.Buffer
		err := tmpl.ExecuteTemplate(&buf, "modals/create", nil)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Create Post")
	})

	t.Run("edit_form_template", func(t *testing.T) {
		var buf bytes.Buffer
		post := &domain.Post{Title: "Edit Post", Content: "Edit content"}
		err := tmpl.ExecuteTemplate(&buf, "modals/edit-content", post)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Edit Post")
		assert.Contains(t, buf.String(), "Edit content")
	})

	t.Run("view_content_template", func(t *testing.T) {
		var buf bytes.Buffer
		post := &domain.Post{Title: "View Post", Content: "View content"}
		err := tmpl.ExecuteTemplate(&buf, "modals/view-content", post)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "View Post")
		assert.Contains(t, buf.String(), "View content")
	})

	t.Run("delete_confirmation_template", func(t *testing.T) {
		var buf bytes.Buffer
		post := &domain.Post{Title: "Delete Post"}
		err := tmpl.ExecuteTemplate(&buf, "modals/delete-content", post)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Delete Post")
	})
}

func TestTemplates_ErrorHandling(t *testing.T) {
	tmpl := template.New("")
	tmpl = template.Must(tmpl.Parse(`
		{{define "test"}}
			{{.NonExistentField}}
		{{end}}
	`))

	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "test", struct{}{})
	assert.Error(t, err)
}
