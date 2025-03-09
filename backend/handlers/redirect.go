package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

func RedirectToRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")
	w.Header().Set("HX-Redirect", "/recipe/"+recipeID)
}

func RedirectToHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Redirect", "/")
}

func RedirectToAddRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Redirect", "/recipe/add")
}
