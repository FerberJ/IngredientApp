package handlers

import (
	"gotth/template/backend/configuration"
	"gotth/template/backend/dao"
	"gotth/template/view/components"
	"net/http"
)

func HandleRecipes(w http.ResponseWriter, r *http.Request, cfg configuration.Configutration) {
	recipes, err := dao.ListRecipes(w, r)
	if err != nil {
		return
	}

	filteredRecipes, selectedFilters := dao.FilterRecipeCards(w, r, recipes)

	components.RecipesList(filteredRecipes, selectedFilters, cfg).Render(r.Context(), w)
}
