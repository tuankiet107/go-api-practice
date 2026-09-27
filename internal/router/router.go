package router

import (
	"database/sql"
	"net/http"

	"go-api-practice/internal/handler"
	"go-api-practice/internal/repository"
	"go-api-practice/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(db *sql.DB) http.Handler {
	userRepository := repository.NewPostgresUserRepository(db)
	return newRouter(userRepository)
}

func newRouter(userRepository repository.UserRepository) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.Health)

	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", userHandler.GetAll)
		r.Get("/{id}", userHandler.GetByID)
		r.Post("/", userHandler.Create)
		r.Put("/{id}", userHandler.Update)
		r.Delete("/{id}", userHandler.Delete)
	})

	return r
}
