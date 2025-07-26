package main

import (
	"gotth/template/backend/auth"
	"gotth/template/backend/configuration"
	"gotth/template/backend/db"
	"gotth/template/backend/handlers"
	"gotth/template/backend/store"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {

	cfg := configuration.SetConfiguration()

	db.NewMinioProvider(cfg)
	db.NewBadgerProvider(cfg)

	auth.InitCasdoor(cfg)
	store.InitStore()

	cloverProvider, err := db.NewCloverProvider(cfg)
	if err != nil {
		return
	}
	defer cloverProvider.Close()

	router := chi.NewMux()

	// Public Files
	router.Handle("/*", public())

	// Endpoints for Login / Logout
	router.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleLogin(w, r, cfg)
	})
	router.Get("/logout", handlers.HandleLogout)
	router.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleLoginCallback(w, r, cfg)
	})

	// Redirect to new paths
	router.Get("/redirect/recipe/{id}", handlers.RedirectToRecipe)
	router.Get("/redirect/recipe/add", handlers.RedirectToAddRecipe)
	router.Get("/redirect/recipe/edit/{id}", handlers.RedirectToEditRecipe)
	router.Get("/redirect/home", handlers.RedirectToHome)

	// Page for Recipe List
	router.Get("/", handlers.HandleListPage) // Recipe List Page
	router.Get("/recipes", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRecipes(w, r, cfg)
	})
	router.Put("/addlistbadges/{keyword}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleAddClosableBadge(w, r, cfg)
	})
	router.Put("/removelistbadges/{keyword}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRemoveClosableBadge(w, r, cfg)
	})
	router.Put("/removelistbadges", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRemoveAllClosableBadge(w, r, cfg)
	})

	// Page for showing a single Recipe
	router.Get("/recipe/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRecipePage(w, r, cfg)
	}) // Recipe Page
	router.Get("/recipe/{id}/servings/{count}", handlers.HandleServings)

	// Page for showing add a new Recipe
	router.Get("/recipe/add", handlers.HandleAddRecipePage) // Add Recipe Page
	router.Get("/recipe/add/ingredient", handlers.HandleAddIngredientInput)
	router.Delete("/recipe/add/ingredient", handlers.HandleRemoveIngredientInput)
	router.Get("/recipe/add/instruction", handlers.HandleAddInstructionInput)
	router.Delete("/recipe/add/instruction", handlers.HandleRemoveInstructionInput)
	router.Post("/recipe", handlers.HandleAddRecipe)
	router.Get("/recipe/add/keyword", handlers.HandleAddRecipeAddBadge)
	router.Delete("/recipe/add/keyword", handlers.HandleAddRecipeRemoveBadge)
	router.Get("/recipe/import/url", handlers.HandleRecipeImportUrl)

	router.Get("/images/{image}", handlers.HandleImageGet)

	// Page for editing a existing Recipe
	router.Get("/recipe/edit/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleEditRecipePage(w, r, cfg)
	}) // Edit Recipe Page
	router.Put("/recipe/{id}", handlers.HandleEditRecipe)
	router.Delete("/recipe/{id}", handlers.HandleDeleteRecipe)

	router.Get("/recipe/bring/{id}", handlers.HandleBringRequest)

	// listenAddr := os.Getenv("LISTEN_ADDR")
	// slog.Info("HTTP server started", "listenAddr", listenAddr)

	http.ListenAndServe(":3000", router)
}
