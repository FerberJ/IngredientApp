package handlers

import (
	"gotth/template/backend/configuration"
	"gotth/template/backend/dao"
	"gotth/template/backend/store"
	"gotth/template/view/components"
	"net/http"
)

func HandleRecipes(w http.ResponseWriter, r *http.Request, cfg configuration.Configutration) {
	var selectedBadges []string
	var selectedSearches []string

	s := store.GetStore()
	valB, err := s.GetValue("badgeList", w, r)
	if err == nil && valB != nil {
		selectedBadges = valB.([]string)
	}

	valS, err := s.GetValue("searchList", w, r)
	if err == nil && valS != nil {
		selectedSearches = valS.([]string)
	}

	recipes, err := dao.ListRecipes(w, r, selectedBadges, selectedSearches, false)
	if err != nil {
		return
	}

	filteredRecipes, selectedFilters := dao.FilterRecipeCards(w, r, recipes)

	components.RecipesList(filteredRecipes, selectedFilters, selectedSearches, cfg).Render(r.Context(), w)
}
