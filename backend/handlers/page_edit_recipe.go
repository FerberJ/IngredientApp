package handlers

import (
	"gotth/template/backend/auth"
	"gotth/template/backend/configuration"
	"gotth/template/backend/dao"
	"gotth/template/view/home"
	"net/http"

	"github.com/go-chi/chi"
)

func HandleEditRecipePage(w http.ResponseWriter, r *http.Request, cfg configuration.Configutration) {
	recipeID := chi.URLParam(r, "id")
	authenticated := false
	avatar := ""

	// If the User exist the avatar will be set
	user, err := auth.GetUser(w, r)
	if err == nil {
		authenticated = true
		avatar = user.Name
	}

	recipe, err := dao.GetRecipe(w, r, recipeID, false)
	if err != nil {
		return
	}

	if recipe.User != user.Id {
		return
	}

	recipe.Image = cfg.AppAddress + "/images/" + recipe.Image

	home.EditRecipeIndex(avatar, authenticated, recipe).Render(r.Context(), w)
}
