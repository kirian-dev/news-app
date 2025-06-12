package server

import (
	"html/template"
	posthandler "news-app/internal/handlers/post"
	postrepo "news-app/internal/repository/post"
	postservice "news-app/internal/services/post"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) Handlers() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

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

	tmpl = template.Must(tmpl.ParseGlob("templates/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("templates/post/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("templates/modals/*.html"))

	repo := postrepo.NewMongoRepository(s.mongo.Client.Database("newsdb"))
	service := postservice.NewService(repo)
	handler := posthandler.New(service, tmpl, s.logger)

	posthandler.RegisterRoutes(r, handler, s.logger)

	s.http.Handler = r
}
