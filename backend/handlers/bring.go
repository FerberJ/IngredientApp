package handlers

import (
	"gotth/template/backend/dao"
	"gotth/template/view/components"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func HandleBringRequest(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")
	servingSiceStr := r.URL.Query().Get("serving")

	recipe, err := dao.GetRecipe(w, r, recipeID, true)
	if err != nil {
		return
	}

	// Get the servingSice. The Serving sice can not be below 1
	servingSice := recipe.Nutrition.ServingSize
	if servingSiceStr != "" {
		num, err := strconv.ParseInt(servingSiceStr, 10, 64)
		if err == nil {
			if num >= 1 {
				servingSice = int(num)
			}
		}
	}

	dao.GetServing(servingSice, &recipe)

	components.BringRecipe(recipe.Name, recipe.Ingredients).Render(r.Context(), w)
}
