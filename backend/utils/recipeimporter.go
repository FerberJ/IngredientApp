package utils

import (
	"encoding/json"
	"gotth/template/backend/models"
	"log"
	"net/http"
	neturl "net/url" //
	"strconv"
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
			recipe.Image = s.Find("[itemprop='thumbnail']").First().Text()

			// IMAGE-------

			// And in Method 2, replace the current image extraction with:
			// Find image in microdata
			imgSrc, imgExists := s.Find("[itemprop='image']").First().Attr("src")
			if imgExists && imgSrc != "" {
				recipe.Image = imgSrc
			} else {
				// Try content attribute
				imgContent, contentExists := s.Find("[itemprop='image']").First().Attr("content")
				if contentExists && imgContent != "" {
					recipe.Image = imgContent
				} else {
					// Try nested img tag
					nestedImg, nestedImgExists := s.Find("[itemprop='image'] img").First().Attr("src")
					if nestedImgExists && nestedImg != "" {
						recipe.Image = nestedImg
					}
				}
			}

			// Also add a fallback method at the end of ImportRecipe if no image was found:
			// Fallback to looking for the largest/most prominent image if no schema-based image was found
			if recipe.Image == "" {
				var largestImg string
				var largestArea int

				doc.Find("img").Each(func(i int, s *goquery.Selection) {
					src, exists := s.Attr("src")
					if !exists || src == "" {
						return
					}

					width, _ := s.Attr("width")
					height, _ := s.Attr("height")

					// Convert attributes to integers (with defaults if missing)
					w := 0
					h := 0
					if width != "" {
						if wInt, err := strconv.Atoi(width); err == nil {
							w = wInt
						}
					}
					if height != "" {
						if hInt, err := strconv.Atoi(height); err == nil {
							h = hInt
						}
					}

					// Compute area - if dimensions exist
					area := w * h

					// If no dimensions in attributes, prioritize by position in document
					if area == 0 {
						area = 1000 - i // Earlier images get higher priority
					}

					if area > largestArea {
						largestArea = area
						largestImg = src
					}
				})

				if largestImg != "" {
					// Ensure the URL is absolute
					if !strings.HasPrefix(largestImg, "http") {
						base, err := neturl.Parse(url)
						if err == nil {
							imgURL, err := base.Parse(largestImg)
							if err == nil {
								largestImg = imgURL.String()
							}
						}
					}
					recipe.Image = largestImg
				}
			}

			//---------

			// Ingredients
			s.Find("[itemprop='recipeIngredient']").Each(func(i int, s *goquery.Selection) {
				ingredientStr := strings.TrimSpace(s.Text())
				ingredientArr := strings.Split(ingredientStr, " ")
				text := ingredientStr
				var amount float64 = 0
				amountIndex := -1

				for i, s := range ingredientArr {
					a, err := strconv.ParseFloat(s, 10)
					if err == nil {
						amount = a
						amountIndex = i
						break
					}
				}

				if amountIndex > -1 {
					text = strings.Replace(ingredientStr, ingredientArr[amountIndex]+" ", "", 1)
				}

				ingredient := models.Ingredient{
					Text:   text,
					Amount: amount,
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

	// Image

	recipe.CookTime, _ = FormatDuration(recipe.CookTime)
	recipe.PrepTime, _ = FormatDuration(recipe.PrepTime)

	yields := strings.Split(recipe.RecipeYield, " ")
	for _, yield := range yields {
		servingI, err := strconv.ParseInt(yield, 10, 0)
		if err == nil {
			recipe.Nutrition = models.NutritionInfo{
				ServingSize: int(servingI),
			}
			break
		}
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

	// In the parseRecipeFromJSON function, add this image extraction logic
	if image, ok := data["image"]; ok {
		// Handle different image formats in JSON-LD
		switch v := image.(type) {
		case string:
			// Direct URL
			recipe.Image = v
		case map[string]interface{}:
			// ImageObject format
			if url, ok := v["url"].(string); ok {
				recipe.Image = url
			}
		case []interface{}:
			// Array of images, use the first one
			if len(v) > 0 {
				if imgStr, ok := v[0].(string); ok {
					recipe.Image = imgStr
				} else if imgObj, ok := v[0].(map[string]interface{}); ok {
					if url, ok := imgObj["url"].(string); ok {
						recipe.Image = url
					}
				}
			}
		}
	}

	// Extract ingredients (could be an array of strings)
	if ingredients, ok := data["recipeIngredient"].([]interface{}); ok {
		for _, ing := range ingredients {
			if ingStr, ok := ing.(string); ok {
				ingredientArr := strings.Split(ingStr, " ")
				text := ingStr
				var amount float64 = 0
				amountIndex := -1

				for i, s := range ingredientArr {
					a, err := strconv.ParseFloat(s, 10)
					if err == nil {
						amount = a
						amountIndex = i
						break
					}
				}

				if amountIndex > -1 {
					text = strings.Replace(ingStr, ingredientArr[amountIndex]+" ", "", 1)
				}

				ingredient := models.Ingredient{
					Text:   text,
					Amount: amount,
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
