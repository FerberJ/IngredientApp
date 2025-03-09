package handlers

import (
	"gotth/template/backend/auth"
	add_recipe_components "gotth/template/view/components/addRecipe"
	"gotth/template/view/home"
	"net/http"
)

func HandleAddRecipePage(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	avatar := ""

	// If the User exist the avatar will be set
	user, err := auth.GetUser(w, r)
	if err == nil {
		authenticated = true
		avatar = user.Avatar
	}

	home.CreateRecipeIndex(avatar, authenticated).Render(r.Context(), w)
}

func HandleAddRecipeAddBadge(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	add_recipe_components.ClosableBadge(keyword).Render(r.Context(), w)
}

func HandleAddRecipeRemoveBadge(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
