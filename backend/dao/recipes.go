package dao

import (
	"encoding/json"
	"gotth/template/backend/auth"
	"gotth/template/backend/db"
	"gotth/template/backend/models"
	"gotth/template/backend/repository"
	"gotth/template/backend/store"
	"gotth/template/backend/utils"
	"net/http"

	"github.com/ostafen/clover/v2/query"
)

func ListRecipes(w http.ResponseWriter, r *http.Request, badges []string, searches []string) ([]models.RecipeCard, error) {
	var recipes []models.RecipeCard
	var q *query.Query
	var c query.Criteria

	recipeRepository := repository.NewRecipeRepository(db.GetCloverProvider())

	user, err := auth.GetUser(w, r)
	if err != nil {
		c = query.Field("private").IsFalse()
	} else {
		c = query.Field("private").IsFalse().Or(query.Field("user").Eq(user.Id))
	}

	if len(badges) > 0 {
		for _, badge := range badges {
			c = c.And(query.Field("keywords").Contains(badge))
		}
	}

	if len(searches) > 0 {
		for _, search := range searches {
			c = c.And(query.Field("name").Like(search).
				Or((query.Field("description").Like(search))))
		}
	}

	q = query.NewQuery(recipeRepository.Collection).Where(c)

	res, err := recipeRepository.FindDocuments(q, nil)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return recipes, err
	}

	for _, resRecipe := range res {
		var recipe models.RecipeCard
		data, _ := json.Marshal(resRecipe)
		json.Unmarshal(data, &recipe)
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

// Get Recipe from the ID
// According to the permission it can be possible that the Recipe will not be found.
func GetRecipe(w http.ResponseWriter, r *http.Request, id string, forBring bool) (models.Recipe, error) {
	var q *query.Query
	var recipe models.Recipe

	recipeRepository := repository.NewRecipeRepository(db.GetCloverProvider())
	user, err := auth.GetUser(w, r)
	if err != nil {
		q = query.NewQuery(recipeRepository.Collection).Where(query.Field("_id").Eq(id).And(query.Field("private").IsFalse()))
	} else if forBring {
		q = query.NewQuery(recipeRepository.Collection).Where(query.Field("_id").Eq(id))
	} else {
		q = query.NewQuery(recipeRepository.Collection).
			Where(query.Field("_id").Eq(id).
				And(query.Field("private").Eq(false).Or(query.Field("user").Eq(user.Id))))
	}

	res, err := recipeRepository.FindDocument(q, nil)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return recipe, err
	}

	data, _ := json.Marshal(res)
	json.Unmarshal(data, &recipe)

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
