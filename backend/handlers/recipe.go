package handlers

import (
	"fmt"
	"gotth/template/backend/auth"
	"gotth/template/backend/dao"
	"gotth/template/backend/db"
	"gotth/template/backend/models"
	"gotth/template/backend/repository"
	"gotth/template/backend/utils"
	recipe_components "gotth/template/view/components/recipe"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func HandleServings(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")
	servingCount := chi.URLParam(r, "count")

	recipe, err := dao.GetRecipe(w, r, recipeID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	servingSice := recipe.Nutrition.ServingSize
	if servingCount != "" {
		num, err := strconv.ParseInt(servingCount, 10, 64)
		if err == nil {
			if num < 1 {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Cant go below 1 Portion"))
				return
			} else {
				servingSice = int(num)
			}
		}
	}

	dao.GetServing(w, r, servingSice, &recipe)

	w.Header().Set("HX-Replace-Url", "/recipe/"+recipeID+"?serving="+strconv.Itoa(servingSice))
	recipe_components.Servings(recipe.Ingredients, uint(servingSice), recipeID).Render(r.Context(), w)
}

func HandleAddRecipe(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(w, r)

	mongoRepository := repository.NewRecipeRepository(db.GetMongoProvider())
	err = r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	privateStr := r.FormValue("private")
	description := r.FormValue("description")
	prepTime := r.FormValue("preptime")
	cookTime := r.FormValue("cooktime")
	totalTime := utils.GetTotalTime(prepTime, cookTime)
	servings, _ := strconv.Atoi(r.FormValue("servings"))

	// Ingredients
	amounts := r.Form["amount"]
	units := r.Form["unit"]
	ingredient := r.Form["ingredient"]

	// Instructions
	instruction := r.Form["instruction"]
	instructionDescription := r.Form["instruction_description"]
	keywords := r.Form["keyword"]

	// Improved file handling
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Cant handle image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageName, err := dao.AddImage(file, handler)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ingredients []models.Ingredient
	for i := 0; i < len(amounts); i++ {
		fAmount, err := strconv.ParseFloat(amounts[i], 64)
		if err != nil {
			continue
		}

		ingredient := models.Ingredient{
			Amount: fAmount,
			Text:   ingredient[i],
			Unit:   units[i],
		}
		ingredients = append(ingredients, ingredient)
	}

	var instructions []models.Instruction
	for i := 0; i < len(instruction); i++ {
		instruction := models.Instruction{
			Header: instruction[i],
			Text:   instructionDescription[i],
		}
		instructions = append(instructions, instruction)
	}

	recipe := models.Recipe{
		Name:        title,
		Description: description,
		Private:     privateStr == "on",
		PrepTime:    prepTime,
		CookTime:    cookTime,
		TotalTime:   totalTime,
		Nutrition: models.NutritionInfo{
			ServingSize: servings,
		},
		Ingredients:  ingredients,
		Instructions: instructions,
		Image:        imageName,
		User:         user.Id,
		UserName:     user.Name,
		Keywords:     keywords,
	}

	id, err := mongoRepository.InsertDocument(recipe, nil)

	fmt.Println(id, id.Hex())
	w.Header().Set("HX-Replace-Url", "/recipe/"+id.Hex()+"?serving="+strconv.Itoa(servings))
}
