package handlers

import (
	"gotth/template/backend/auth"
	"gotth/template/backend/db"
	"gotth/template/backend/models"
	"gotth/template/backend/repository"
	"gotth/template/backend/store"
	"gotth/template/backend/utils"
	"gotth/template/view/components"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func HandleRecipes(w http.ResponseWriter, r *http.Request) {
	var filter bson.M
	recipeRepository := repository.NewRecipeRepository(db.GetMongoProvider())
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

	var valSlice []string
	s := store.GetStore()
	val, err := s.GetValue("badgeList", w, r)
	if err == nil && val != nil {
		valSlice = val.([]string)
	}

	var filteredRecipes []models.RecipeCard

	if len(valSlice) == 0 {
		filteredRecipes = recipes
	} else {
		for _, rec := range recipes {
			if utils.ContainsAllKeywords(rec.Keywords, valSlice) {
				filteredRecipes = append(filteredRecipes, rec)
			}
		}
	}

	components.RecipesList(filteredRecipes, valSlice).Render(r.Context(), w)
}
