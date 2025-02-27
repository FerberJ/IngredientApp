package handlers

import (
	"gotth/template/backend/auth"
	"gotth/template/backend/db"
	"gotth/template/backend/models"
	"gotth/template/backend/repository"
	"gotth/template/view/components"
	"net/http"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleRecipe(w http.ResponseWriter, r *http.Request) {
	var filter bson.M
	recipeID := chi.URLParam(r, "id")
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

	components.Recipe(recipe).Render(r.Context(), w)
}
