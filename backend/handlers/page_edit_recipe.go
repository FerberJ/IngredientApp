package handlers

import (
	"gotth/template/backend/auth"
	"gotth/template/backend/dao"
	"gotth/template/view/home"
	"net/http"

	"github.com/go-chi/chi"
)

func HandleEditRecipePage(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")
	authenticated := false
	avatar := ""

	// If the User exist the avatar will be set
	user, err := auth.GetUser(w, r)
	if err == nil {
		authenticated = true
		avatar = user.Avatar
	}

	recipe, err := dao.GetRecipe(w, r, recipeID)
	if err != nil {
		return
	}

	recipe.Image = "http://localhost:3000/images/" + recipe.Image

	home.EditRecipeIndex(avatar, authenticated, recipe).Render(r.Context(), w)
}
