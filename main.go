package main

import (
	"fmt"
	"gotth/template/backend/auth"
	"gotth/template/backend/db"
	"gotth/template/backend/handlers"
	"gotth/template/backend/store"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db.NewRedisProvider("localhost:6379", "")

	auth.InitCasdoor()
	store.InitStore()

	mongoProvider, err := db.NewMongoProvider("mongodb://localhost:27017", "recipesDb")
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	defer mongoProvider.Close()

	router := chi.NewMux()

	// Public Files
	router.Handle("/*", public())

	// Endpoints for Login / Logout
	router.Get("/login", handlers.HandleLogin)
	router.Get("/logout", handlers.HandleLogout)
	router.Get("/callback", handlers.HandleLoginCallback)

	// Redirect to new paths
	router.Get("/redirect/recipe/{id}", handlers.RedirectToRecipe)
	router.Get("/redirect/home", handlers.RedirectToHome)

	// Page for Recipe List
	router.Get("/", handlers.HandleListPage) // Recipe List Page
	router.Get("/recipes", handlers.HandleRecipes)
	router.Put("/addlistbadges/{keyword}", handlers.HandleAddClosableBadge)
	router.Put("/removelistbadges/{keyword}", handlers.HandleRemoveClosableBadge)
	router.Put("/removelistbadges", handlers.HandleRemoveAllClosableBadge)

	// Page for showing a single Recipe
	router.Get("/recipe/{id}", handlers.HandleRecipePage) // Recipe Page
	router.Get("/recipe/{id}/servings/{count}", handlers.HandleServings)

	listenAddr := os.Getenv("LISTEN_ADDR")
	slog.Info("HTTP server started", "listenAddr", listenAddr)

	http.ListenAndServe(listenAddr, router)
}
