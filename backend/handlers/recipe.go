package handlers

import (
	"gotth/template/backend/dao"
	recipe_components "gotth/template/view/components/recipe"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func HandleServings(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")
	servingCount := chi.URLParam(r, "count")

	recipe, err := dao.GetRecipe(w, r, recipeID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	servingSice := recipe.Nutrition.ServingSize
	if servingCount != "" {
		num, err := strconv.ParseInt(servingCount, 10, 64)
		if err == nil {
			if num < 1 {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Cant go below 1 Portion"))
				return
			} else {
				servingSice = int(num)
			}
		}
	}

	dao.GetServing(w, r, servingSice, &recipe)

	w.Header().Set("HX-Replace-Url", "/recipe/"+recipeID+"?serving="+strconv.Itoa(servingSice))
	recipe_components.Servings(recipe.Ingredients, uint(servingSice), recipeID).Render(r.Context(), w)
}
