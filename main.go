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

	router.Handle("/*", public())
	router.Get("/", handlers.HandleHome)
	router.Get("/login", handlers.HandleLogin)
	router.Get("/logout", handlers.HandleLogout)
	router.Get("/recipies", handlers.HandleRecipes)
	router.Get("/callback", handlers.HandleLoginCallback)
	router.Put("/addlistbadges/{keyword}", handlers.HandleAddClosableBadge)
	router.Put("/removelistbadges/{keyword}", handlers.HandleRemoveClosableBadge)
	router.Put("/removelistbadges", handlers.HandleRemoveAllClosableBadge)

	listenAddr := os.Getenv("LISTEN_ADDR")
	slog.Info("HTTP server started", "listenAddr", listenAddr)

	http.ListenAndServe(listenAddr, router)
}
