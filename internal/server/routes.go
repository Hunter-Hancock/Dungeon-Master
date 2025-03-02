package server

import (
	"Hunter-Hancock/dungeon-master/internal/db"
	"Hunter-Hancock/dungeon-master/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/joho/godotenv/autoload"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/google"
	"os"
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

	goth.UseProviders(
		discord.New(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"), "http://localhost:6969/auth/discord/callback", "identify", "email"),
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://localhost:6969/auth/google/callback"),
	)

	authStore := db.NewAuthStore(s.db)
	authHandler := handler.NewAuthHandler(authStore)

	r.Get("/auth/{provider}/callback", authHandler.Login)
	r.Get("/auth/{provider}", authHandler.Begin)
	r.Get("/auth/logout", authHandler.Logout)
	r.Get("/auth/me", authHandler.Me)

	return r
}
