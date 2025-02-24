package handlers

import (
	"gotth/template/backend/auth"
	"gotth/template/backend/db"
	"gotth/template/backend/models"
	"gotth/template/backend/repository"
	"gotth/template/view/components"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func HandeleRecipes(w http.ResponseWriter, r *http.Request) {
	var filter bson.M
	recipeRepository := repository.NewRecipeRepository(db.GetProvider())
	user, err := auth.GetUser(w, r)
	if err != nil {
		filter = bson.M{"private": false}
	} else {
		filter = bson.M{"$or": []bson.M{{"private": false}, {"user": user.Id}}}

	}

	res, err := recipeRepository.FindDocumentsFields(filter, models.GetRecipeCardFilter(), nil)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	var recipes []models.RecipeCard

	for _, resRecipe := range res {
		var recipe models.RecipeCard
		data, _ := bson.Marshal(resRecipe)
		bson.Unmarshal(data, &recipe)
		recipes = append(recipes, recipe)
	}

	components.RecipesList(recipes).Render(r.Context(), w)
}
