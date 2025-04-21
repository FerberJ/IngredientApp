package handlers

import (
	"gotth/template/backend/auth"
	"gotth/template/backend/configuration"
	"gotth/template/backend/dao"
	"gotth/template/view/home"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func HandleRecipePage(w http.ResponseWriter, r *http.Request, cfg configuration.Configutration) {
	authenticated := false
	avatar := ""

	// If the User exist the avatar will be set
	user, err := auth.GetUser(w, r)
	if err == nil {
		authenticated = true
		avatar = user.Name
	}

	recipeID := chi.URLParam(r, "id")
	servingSiceStr := r.URL.Query().Get("serving")

	recipe, err := dao.GetRecipe(w, r, recipeID, false)
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

	isUsers := false
	if user != nil {
		isUsers = recipe.User == user.Id
	}

	home.RecipeIndex(avatar, authenticated, recipe, uint(recipe.Nutrition.ServingSize), isUsers, cfg).Render(r.Context(), w)
}
