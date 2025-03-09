package utils

import (
	"encoding/json"
	"gotth/template/backend/models"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ImportRecipe(url string) models.Recipe {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Add a user agent to appear more like a browser
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// Make HTTP GET request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch the website: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is OK (200)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code: %d", resp.StatusCode)
	}

	// Load HTML document from the response body
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	// Method 1: Extract recipe data from JSON-LD script tags
	var recipe models.Recipe
	doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
		scriptContent := s.Text()

		// Check if the script contains Recipe schema
		if strings.Contains(scriptContent, "schema.org") && strings.Contains(scriptContent, "Recipe") {
			// Try to parse the JSON
			var jsonData map[string]interface{}
			if err := json.Unmarshal([]byte(scriptContent), &jsonData); err == nil {
				// Handle @graph structure where Recipe might be in an array
				if graph, hasGraph := jsonData["@graph"].([]interface{}); hasGraph {
					for _, item := range graph {
						if itemMap, isMap := item.(map[string]interface{}); isMap {
							if itemType, hasType := itemMap["@type"].(string); hasType && itemType == "Recipe" {
								parseRecipeFromJSON(itemMap, &recipe)
								break
							}
						}
					}
				} else if recipeType, hasType := jsonData["@type"].(string); hasType && recipeType == "Recipe" {
					// Direct Recipe object
					parseRecipeFromJSON(jsonData, &recipe)
				}
			}
		}
	})

	// Method 2: Fallback to microdata/RDFa if JSON-LD not found
	if recipe.Name == "" {
		// Find elements with itemtype="http://schema.org/Recipe" or similar
		doc.Find("[itemtype$='Recipe'], [typeof$='Recipe']").Each(func(i int, s *goquery.Selection) {
			// Basic recipe info
			recipe.Name = s.Find("[itemprop='name']").First().Text()
			recipe.Description = s.Find("[itemprop='description']").First().Text()
			//recipe.Author = s.Find("[itemprop='author']").First().Text()
			recipe.RecipeYield = s.Find("[itemprop='recipeYield']").First().Text()
			recipe.RecipeCategory = s.Find("[itemprop='recipeCategory']").First().Text()
			recipe.RecipeCuisine = s.Find("[itemprop='recipeCuisine']").First().Text()
			recipe.PrepTime, _ = s.Find("[itemprop='prepTime']").First().Attr("content")
			recipe.CookTime, _ = s.Find("[itemprop='cookTime']").First().Attr("content")
			recipe.TotalTime, _ = s.Find("[itemprop='totalTime']").First().Attr("content")

			// Ingredients
			s.Find("[itemprop='recipeIngredient']").Each(func(i int, s *goquery.Selection) {
				ingredientStr := strings.TrimSpace(s.Text())
				ingredient := models.Ingredient{
					Text: ingredientStr,
				}
				if ingredientStr != "" {
					recipe.Ingredients = append(recipe.Ingredients, ingredient)
				}
			})

			// Instructions
			s.Find("[itemprop='recipeInstructions']").Each(func(i int, s *goquery.Selection) {
				instructionStr := strings.TrimSpace(s.Text())
				instruction := models.Instruction{
					Text: instructionStr,
				}
				if instructionStr != "" {
					recipe.Instructions = append(recipe.Instructions, instruction)
				}
			})

			// Nutrition
			recipe.Nutrition.Calories = s.Find("[itemprop='calories']").First().Text()
		})
	}

	return recipe
}

func parseRecipeFromJSON(data map[string]interface{}, recipe *models.Recipe) {
	// Extract basic info
	if name, ok := data["name"].(string); ok {
		recipe.Name = name
	}
	if desc, ok := data["description"].(string); ok {
		recipe.Description = desc
	}

	if prepTime, ok := data["prepTime"].(string); ok {
		recipe.PrepTime = prepTime
	}
	if cookTime, ok := data["cookTime"].(string); ok {
		recipe.CookTime = cookTime
	}
	if totalTime, ok := data["totalTime"].(string); ok {
		recipe.TotalTime = totalTime
	}
	if keywords, ok := data["keywords"].(string); ok {
		recipe.Keywords = append(recipe.Keywords, keywords)
	}
	if yield, ok := data["recipeYield"].(string); ok {
		recipe.RecipeYield = yield
	}
	if category, ok := data["recipeCategory"].(string); ok {
		recipe.RecipeCategory = category
	}
	if cuisine, ok := data["recipeCuisine"].(string); ok {
		recipe.RecipeCuisine = cuisine
	}

	// Extract ingredients (could be an array of strings)
	if ingredients, ok := data["recipeIngredient"].([]interface{}); ok {
		for _, ing := range ingredients {
			if ingStr, ok := ing.(string); ok {
				ingredient := models.Ingredient{
					Text: ingStr,
				}
				recipe.Ingredients = append(recipe.Ingredients, ingredient)
			}
		}
	}

	// Extract instructions (could be an array of strings or objects)
	if instructions, ok := data["recipeInstructions"].([]interface{}); ok {
		for _, inst := range instructions {
			if instStr, ok := inst.(string); ok {
				instruction := models.Instruction{
					Text: instStr,
				}
				recipe.Instructions = append(recipe.Instructions, instruction)
			} else if instObj, ok := inst.(map[string]interface{}); ok {
				if text, ok := instObj["text"].(string); ok {
					instruction := models.Instruction{
						Text: text,
					}
					recipe.Instructions = append(recipe.Instructions, instruction)
				}
			}
		}
	}

	// Extract nutrition
	if nutrition, ok := data["nutrition"].(map[string]interface{}); ok {
		if calories, ok := nutrition["calories"].(string); ok {
			recipe.Nutrition.Calories = calories
		}
	}
}
