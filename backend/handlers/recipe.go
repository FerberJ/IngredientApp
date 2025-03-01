package handlers

import (
	"gotth/template/backend/auth"
	"gotth/template/backend/db"
	"gotth/template/backend/models"
	"gotth/template/backend/repository"
	recipe_components "gotth/template/view/components/recipe"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleServings(w http.ResponseWriter, r *http.Request) {
	var filter bson.M
	recipeID := chi.URLParam(r, "id")
	servingCount := chi.URLParam(r, "count")

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
	if servingCount != "" {
		num, err := strconv.ParseInt(servingCount, 10, 64)
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

	if servingSice < 1 {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Bad Gateway"))
		return
	}

	w.Header().Set("HX-Replace-Url", "/recipe/"+recipeID+"?serving="+strconv.Itoa(servingSice))
	recipe_components.Servings(recipe.Ingredients, uint(servingSice), recipeID).Render(r.Context(), w)
}
