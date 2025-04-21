package dao

import (
	"gotth/template/backend/auth"
	"gotth/template/backend/db"
	"gotth/template/backend/models"
	"gotth/template/backend/repository"
	"gotth/template/backend/store"
	"gotth/template/backend/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListRecipes(w http.ResponseWriter, r *http.Request) ([]models.RecipeCard, error) {
	var recipes []models.RecipeCard
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
		return recipes, err
	}

	for _, resRecipe := range res {
		var recipe models.RecipeCard
		data, _ := bson.Marshal(resRecipe)
		bson.Unmarshal(data, &recipe)
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

// Get Recipe from the ID
// According to the permission it can be possible that the Recipe will not be found.
func GetRecipe(w http.ResponseWriter, r *http.Request, id string, forBring bool) (models.Recipe, error) {
	var filter bson.M
	var recipe models.Recipe
	recipeIDObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid recipe ID format"))
		return recipe, err
	}
	recipeRepository := repository.NewRecipeRepository(db.GetMongoProvider())
	user, err := auth.GetUser(w, r)
	if err != nil {
		filter = bson.M{"_id": recipeIDObjectID, "private": false}
	} else if forBring {
		filter = bson.M{"_id": recipeIDObjectID}
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
		return recipe, err
	}

	data, _ := bson.Marshal(res)
	bson.Unmarshal(data, &recipe)

	return recipe, nil
}

// Ajust the Ingredients according to servingSiceStr
// If the servingSiceStr is empty, it will return the default servingSice that is defined in the recipe
func GetServing(servingSice int, recipe *models.Recipe) {
	ingredients := make([]models.Ingredient, 0, len(recipe.Ingredients))
	devider := float64(recipe.Nutrition.ServingSize) / float64(servingSice)
	for _, ingredient := range recipe.Ingredients {
		ingre := ingredient
		ingre.Amount = ingre.Amount / devider
		ingredients = append(ingredients, ingre)
	}

	recipe.Ingredients = ingredients
	recipe.Nutrition.ServingSize = servingSice
}

// Filter the Recipe Cards. Filte it by the badgelist on the redis store
func FilterRecipeCards(w http.ResponseWriter, r *http.Request, recipes []models.RecipeCard) ([]models.RecipeCard, []string) {
	var selectedFilters []string
	var filteredRecipes []models.RecipeCard

	s := store.GetStore()
	val, err := s.GetValue("badgeList", w, r)
	if err == nil && val != nil {
		selectedFilters = val.([]string)
	}

	if len(selectedFilters) == 0 {
		filteredRecipes = recipes
	} else {
		for _, rec := range recipes {
			if utils.ContainsAllKeywords(rec.Keywords, selectedFilters) {
				filteredRecipes = append(filteredRecipes, rec)
			}
		}
	}

	return filteredRecipes, selectedFilters
}
