package models

import (
	"time"
)

type RecipeCard struct {
	ID             string    `json:"_id,omitempty"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Image          string    `json:"image"`
	User           string    `json:"user"`
	UserName       string    `json:"userName"`
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

/*
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
*/

// Recipe represents a cooking recipe with all the relevant properties from Schema.org.
type Recipe struct {
	ID                 string          `json:"_id,omitempty"`
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	Image              string          `json:"image"`
	User               string          `json:"user"`
	UserName           string          `json:"userName"`
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
	Tip                string          `json:"tip"`
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

type RecipeUpdate struct {
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	Image              string          `json:"image"`
	User               string          `json:"user"`
	UserName           string          `json:"userName"`
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
	Tip                string          `json:"tip"`
}
