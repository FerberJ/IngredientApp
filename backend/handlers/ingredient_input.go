package handlers

import (
	add_recipe_components "gotth/template/view/components/addRecipe"
	"net/http"
)

func HandleAddIngredientInput(w http.ResponseWriter, r *http.Request) {
	add_recipe_components.IngredientInput().Render(r.Context(), w)
}

func HandleRemoveIngredientInput(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
