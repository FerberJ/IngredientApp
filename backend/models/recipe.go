package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type RecipeCard struct {
	ID             string    `bson:"_id,omitempty"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Image          string    `json:"image"`
	User           string    `json:"user"`
	Private        bool      `json:"private"`
	RecipeYield    string    `json:"recipeYield"`
	RecipeCategory string    `json:"recipeCategory"`
	RecipeCuisine  string    `json:"recipeCuisine"`
	PrepTime       string    `json:"prepTime"`
	CookTime       string    `json:"cookTime"`
	TotalTime      string    `json:"totalTime"`
	Cuisine        string    `json:"cuisine"`
	CreatedAt      time.Time `json:"createdAt"`
	Keywords       []string  `json:"keywords"`
}

func GetRecipeCardFilter() bson.M {
	return bson.M{
		"name":           1,
		"description":    1,
		"private":        1,
		"image":          1,
		"user":           1,
		"recipeYield":    1,
		"recipeCategory": 1,
		"recipeCuisine":  1,
		"prepTime":       1,
		"cookTime":       1,
		"totalTime":      1,
		"cuisine":        1,
		"createdAt":      1,
		"keywords":       1,
	}
}

// Recipe represents a cooking recipe with all the relevant properties from Schema.org.
type Recipe struct {
	ID                 string          `bson:"_id,omitempty"`
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	Image              string          `json:"image"`
	User               string          `json:"user"`
	Private            bool            `json:"private"`
	RecipeYield        string          `json:"recipeYield"`
	RecipeCategory     string          `json:"recipeCategory"`
	RecipeCuisine      string          `json:"recipeCuisine"`
	PrepTime           string          `json:"prepTime"`
	CookTime           string          `json:"cookTime"`
	TotalTime          string          `json:"totalTime"`
	Ingredients        []Ingredient    `json:"ingredients"`
	Instructions       []Instruction   `json:"instructions"`
	Nutrition          NutritionInfo   `json:"nutrition"`
	AggregateRating    AggregateRating `json:"aggregateRating"`
	Cuisine            string          `json:"cuisine"`
	Keywords           []string        `json:"keywords"`
	CreatedAt          time.Time       `json:"createdAt"`
	RecipeInstructions []string        `json:"recipeInstructions"`
}

type Name struct {
	Name string `json:"name"`
}

type Ingredient struct {
	Text   string  `json:"text"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type Instruction struct {
	Header string `json:"header"`
	Text   string `json:"text"`
	Image  string `json:"image"`
}

// NutritionInfo represents nutritional information as part of the recipe.
type NutritionInfo struct {
	Calories            string `json:"calories"`
	FatContent          string `json:"fatContent"`
	CarbohydrateContent string `json:"carbohydrateContent"`
	ProteinContent      string `json:"proteinContent"`
	ServingSize         int    `json:"servingSize"`
}

// AggregateRating holds information about the recipe's rating.
type AggregateRating struct {
	RatingValue float64 `json:"ratingValue"`
	BestRating  float64 `json:"bestRating"`
	WorstRating float64 `json:"worstRating"`
	RatingCount int     `json:"ratingCount"`
}

func GetRecipe() Recipe {
	return Recipe{
		Name:           "Spaghetti Bolognese",
		Description:    "A classic Italian pasta dish with a rich tomato and meat sauce.",
		Image:          "https://example.com/spaghetti.jpg",
		User:           "John Doe",
		RecipeYield:    "4 servings",
		RecipeCategory: "Main Course",
		RecipeCuisine:  "Italian",
		PrepTime:       "PT15M",
		CookTime:       "PT45M",
		TotalTime:      "PT1H",
		Ingredients: []Ingredient{
			{
				Text:   "Spaghetti",
				Amount: 300,
				Unit:   "g",
			},
			{
				Text:   "Ground beef",
				Amount: 500,
				Unit:   "g",
			},
			{
				Text:   "Tomato sauce",
				Amount: 500,
				Unit:   "g",
			},
			{
				Text:   "Garlic",
				Amount: 3,
			},
		},
		Instructions: []Instruction{
			{
				Text:  "Boil water for spaghetti.",
				Image: "",
			},
			{
				Text:  "Cook ground beef.",
				Image: "",
			},
			{
				Text:  "Add garlic and onion.",
				Image: "",
			},
			{
				Text:  "Add tomato sauce and simmer.",
				Image: "",
			},
		},
		Nutrition: NutritionInfo{
			Calories:            "600",
			FatContent:          "20g",
			CarbohydrateContent: "75g",
			ProteinContent:      "30g",
			ServingSize:         1,
		},
		AggregateRating: AggregateRating{
			RatingValue: 4.5,
			BestRating:  5,
			WorstRating: 1,
			RatingCount: 100,
		},
		Cuisine:            "Italian",
		Keywords:           []string{"Pasta", "Italian", "Bolognese"},
		CreatedAt:          time.Now(),
		RecipeInstructions: []string{"Step 1: Boil water.", "Step 2: Cook meat.", "Step 3: Mix sauce."},
	}
}

func GetRecipeCard() RecipeCard {
	return RecipeCard{
		Name:           "Spaghetti Bolognese",
		Description:    "A classic Italian pasta dish with a rich tomato and meat sauce.",
		Image:          "https://example.com/spaghetti.jpg",
		User:           "John Doe",
		RecipeYield:    "4 servings",
		RecipeCategory: "Main Course",
		RecipeCuisine:  "Italian",
		PrepTime:       "PT15M",
		CookTime:       "PT45M",
		TotalTime:      "PT1H",
		Cuisine:        "Italian",
		CreatedAt:      time.Now(),
		Keywords:       []string{"Pasta", "Italian", "Bolognese"},
	}
}

// GetRecipeCard2 function that returns a second RecipeCard
func GetRecipeCard2() RecipeCard {
	return RecipeCard{
		Name:           "Chicken Alfredo",
		Description:    "A creamy pasta dish made with grilled chicken and Alfredo sauce.",
		Image:          "https://example.com/chicken-alfredo.jpg",
		User:           "Jane Smith",
		RecipeYield:    "4 servings",
		RecipeCategory: "Main Course",
		RecipeCuisine:  "Italian",
		PrepTime:       "PT10M",
		CookTime:       "PT30M",
		TotalTime:      "PT40M",
		Cuisine:        "Italian",
		CreatedAt:      time.Now(),
		Keywords:       []string{"Chicken", "Italian", "Bolognese", "Main Course"},
	}
}
