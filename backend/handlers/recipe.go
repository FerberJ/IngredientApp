package handlers

import (
	"gotth/template/backend/auth"
	"gotth/template/backend/dao"
	"gotth/template/backend/db"
	"gotth/template/backend/models"
	"gotth/template/backend/repository"
	"gotth/template/backend/utils"
	recipe_components "gotth/template/view/components/recipe"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleServings(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")
	servingCount := chi.URLParam(r, "count")

	recipe, err := dao.GetRecipe(w, r, recipeID, false)
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

	dao.GetServing(servingSice, &recipe)

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

	// Setting Time to 00:00 if empty
	if prepTime == "" {
		prepTime = "00:00"
	}

	if cookTime == "" {
		cookTime = "00:00"
	}

	// Ingredients
	amounts := r.Form["amount"]
	units := r.Form["unit"]
	ingredient := r.Form["ingredient"]

	selectedRadio := r.FormValue("radio-image")

	// Instructions
	instruction := r.Form["instruction"]
	instructionDescription := r.Form["instruction_description"]
	keywords := r.Form["keyword"]

	var imageName string
	switch selectedRadio {
	case "", "upload":
		// Improved file handling
		file, handler, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Cant handle image", http.StatusBadRequest)
			return
		}
		defer file.Close()

		imageName, err = dao.AddImage(file, handler)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "import":
		imageUrl := r.FormValue("imageUrl")
		imageName, err = dao.AddImageFromURL(imageUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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

	w.Header().Set("HX-Redirect", "/recipe/"+id.Hex()+"?serving="+strconv.Itoa(servings))
}

func HandleDeleteRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")
	user, err := auth.GetUser(w, r)

	recipe, err := dao.GetRecipe(w, r, recipeID, false)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	if recipe.User != user.Id {
		http.Error(w, "No Authorized", http.StatusBadRequest)
		return
	}

	dao.DeleteImage(recipe.Image)

	mongoRepository := repository.NewRecipeRepository(db.GetMongoProvider())

	objectID, err := primitive.ObjectIDFromHex(recipeID)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": objectID}
	err = mongoRepository.DeleteDocument(filter, nil)
	w.Header().Set("HX-Redirect", "/")
}

func HandleEditRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")
	user, err := auth.GetUser(w, r)

	recipe, err := dao.GetRecipe(w, r, recipeID, false)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}
	recipe.User = user.Id
	recipe.UserName = user.Name

	mongoRepository := repository.NewRecipeRepository(db.GetMongoProvider())
	err = r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	recipe.Name = r.FormValue("title")
	recipe.Private = r.FormValue("private") == "on"
	recipe.Description = r.FormValue("description")
	prepTime := r.FormValue("preptime")
	cookTime := r.FormValue("cooktime")
	totalTime := utils.GetTotalTime(prepTime, cookTime)
	servings, _ := strconv.Atoi(r.FormValue("servings"))
	recipe.Nutrition.ServingSize = servings

	// Setting Time to 00:00 if empty
	if prepTime == "" {
		prepTime = "00:00"
	}

	if cookTime == "" {
		cookTime = "00:00"
	}

	recipe.PrepTime = prepTime
	recipe.CookTime = cookTime
	recipe.TotalTime = totalTime

	// Ingredients
	amounts := r.Form["amount"]
	units := r.Form["unit"]
	ingredient := r.Form["ingredient"]

	selectedRadio := r.FormValue("radio-image")

	// Instructions
	instruction := r.Form["instruction"]
	instructionDescription := r.Form["instruction_description"]
	keywords := r.Form["keyword"]
	recipe.Keywords = append(recipe.Keywords, keywords...)

	var imageName string
	switch selectedRadio {
	case "", "upload":
		// Improved file handling
		file, handler, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Cant handle image", http.StatusBadRequest)
			return
		}
		defer file.Close()

		imageName, err = dao.AddImage(file, handler)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		recipe.Image = imageName
	case "import":
		break
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
	recipe.Ingredients = ingredients

	var instructions []models.Instruction
	for i := 0; i < len(instruction); i++ {
		instruction := models.Instruction{
			Header: instruction[i],
			Text:   instructionDescription[i],
		}
		instructions = append(instructions, instruction)
	}
	recipe.Instructions = instructions

	/*
		recipe = models.Recipe{
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
	*/

	objectID, err := primitive.ObjectIDFromHex(recipeID)
	if err != nil {
		log.Fatal(err)
	}

	recipeUpdate := models.RecipeUpdate{
		Name:               recipe.Name,
		Description:        recipe.Description,
		Image:              recipe.Image,
		User:               recipe.User,
		UserName:           recipe.UserName,
		Private:            recipe.Private,
		RecipeYield:        recipe.RecipeYield,
		RecipeCategory:     recipe.RecipeCategory,
		RecipeCuisine:      recipe.Cuisine,
		PrepTime:           recipe.PrepTime,
		CookTime:           recipe.CookTime,
		TotalTime:          recipe.TotalTime,
		Ingredients:        recipe.Ingredients,
		Instructions:       recipe.Instructions,
		Nutrition:          recipe.Nutrition,
		AggregateRating:    recipe.AggregateRating,
		Cuisine:            recipe.Cuisine,
		Keywords:           recipe.Keywords,
		CreatedAt:          recipe.CreatedAt,
		RecipeInstructions: recipe.RecipeInstructions,
		Tip:                recipe.Tip,
	}

	filter := bson.M{"_id": objectID}
	recipeBson, err := bson.Marshal(recipeUpdate)
	updateMap := bson.M{}
	err = bson.Unmarshal(recipeBson, &updateMap)
	err = mongoRepository.UpdateDocument(filter, updateMap, nil)

	w.Header().Set("HX-Redirect", "/recipe/"+recipeID+"?serving="+strconv.Itoa(servings))
}
