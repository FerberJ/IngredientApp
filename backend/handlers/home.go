package handlers

import (
	"fmt"
	"gotth/template/backend/auth"
	"gotth/template/backend/db"
	"gotth/template/backend/models"
	"gotth/template/backend/repository"
	"gotth/template/view/home"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleListPage(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(w, r)
	if err != nil {
		home.ListIndex("", false).Render(r.Context(), w)
		return
	}

	fmt.Println(user)

	home.ListIndex(user.Avatar, true).Render(r.Context(), w)
}

func HandleRecipePage(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	// Get Recipe
	var filter bson.M
	recipeID := chi.URLParam(r, "id")
	servingSiceStr := r.URL.Query().Get("serving")

	recipeIDObjectID, err := primitive.ObjectIDFromHex(recipeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid recipe ID format"))
		return
	}
	recipeRepository := repository.NewRecipeRepository(db.GetMongoProvider())
	user, err := auth.GetUser(w, r)
	if err != nil {
		filter = bson.M{"_id": recipeIDObjectID, "private": false}
	} else {
		authenticated = true
		filter = bson.M{
			"$and": []bson.M{
				{"_id": recipeIDObjectID},
				{"$or": []bson.M{
					{"private": false},
					{"user": user.Id},
				}},
			},
		}

	}

	res, err := recipeRepository.FindDocument(filter, nil)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	var recipe models.Recipe
	data, _ := bson.Marshal(res)
	bson.Unmarshal(data, &recipe)

	servingSice := recipe.Nutrition.ServingSize
	if servingSiceStr != "" {
		num, err := strconv.ParseInt(servingSiceStr, 10, 64)
		if err == nil {
			servingSice = int(num)
		}
	}

	ingredients := make([]models.Ingredient, 0, len(recipe.Ingredients))
	devider := float64(recipe.Nutrition.ServingSize) / float64(servingSice)
	for _, ingredient := range recipe.Ingredients {
		ingre := ingredient
		ingre.Amount = ingre.Amount / devider
		ingredients = append(ingredients, ingre)
	}

	recipe.Ingredients = ingredients

	avatar := ""
	if user != nil {
		avatar = user.Avatar
	}

	home.RecipeIndex(avatar, authenticated, recipe, uint(servingSice)).Render(r.Context(), w)
}
