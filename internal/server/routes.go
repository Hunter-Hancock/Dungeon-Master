package server

import (
	"Hunter-Hancock/dungeon-master/internal/db"
	"Hunter-Hancock/dungeon-master/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	authStore := db.NewAuthStore(s.db)
	authHandler := handler.NewAuthHandler(authStore)

	r.Get("/test", handler.HelloWorldHandler)
	r.Get("/auth", authHandler.Handle)

	return r
}
