package finder

import (
	"WhatCanIEat/WhatCanIEat/cmd/errors"
	ingredientsPackage "WhatCanIEat/WhatCanIEat/cmd/ingredients"
	"WhatCanIEat/WhatCanIEat/cmd/recipes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type SpoonacularFinder struct{}

func (finder SpoonacularFinder) FindByIngredientsNames(ingredients *[]string,
	numberOfRecipes int) ([]recipes.Recipe, error) {
	/* Get Recipe Id, Title, Ingredients and MissingIngredients. */

	// TODO: Use (if exists) pathlib or sth similar, to auto format, validate and sanitize.
	recipesUrl := "https://api.spoonacular.com/recipes/findByIngredients" +
		"?ingredients=" + strings.Join(*ingredients, ",") +
		"&sort=min-missing-ingredients" +
		"&ranking=2" + // Minimize missing ingredients.
		"&number=" + strconv.Itoa(numberOfRecipes) +
		"&apiKey=" + os.Getenv("SPOONACULAR_API_KEY")

	recipesResponse, err := finder.getValidResponse(recipesUrl)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(recipesResponse.Body)

	parsedRecipes, err := finder.responseToRecipes(recipesResponse)
	if err != nil {
		return nil, err
	}

	/* Get Nutrition. */
	for idx, recipe := range parsedRecipes {
		// TODO: Use (if exists) pathlib or sth similar, to auto format, validate and sanitize.
		nutritionUrl := "https://api.spoonacular.com/recipes/" + strconv.Itoa(recipe.Id) + "/information" +
			"?includeNutrition=true" +
			"&apiKey=" + os.Getenv("SPOONACULAR_API_KEY")

		nutritionResponse, err := finder.getValidResponse(nutritionUrl)
		if err != nil {
			return nil, err
		}
		parsedRecipes[idx].Nutrition, err = finder.responseToNutrition(nutritionResponse)
		if err != nil {
			_ = nutritionResponse.Body.Close()
			return nil, err
		}
		_ = nutritionResponse.Body.Close()
	}

	return parsedRecipes, nil
}

func (finder SpoonacularFinder) getValidResponse(url string) (response *http.Response, err error) {
	//log.Println("Sending request to: \"" + url + "\"")
	// TODO: Move apiKey into secret.

	response, err = http.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, errors.ReceivedBadStatusCode(response.Status)
	}
	log.Println("Response: " + response.Status)

	return response, nil
}

func (finder SpoonacularFinder) responseToRecipes(response *http.Response) ([]recipes.Recipe, error) {
	type ingredientFromJson struct {
		Id   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	type recipeFromJson struct {
		Id                 int                  `json:"id,omitempty"`
		Title              string               `json:"title,omitempty"`
		Ingredients        []ingredientFromJson `json:"usedIngredients,omitempty"`
		MissingIngredients []ingredientFromJson `json:"missedIngredients,omitempty"`
	}

	var recipesFromJson []recipeFromJson
	err := json.NewDecoder(response.Body).Decode(&recipesFromJson)
	if err != nil {
		return nil, err
	}

	parsedRecipes := make([]recipes.Recipe, len(recipesFromJson))
	for idx, recipe := range recipesFromJson {

		ingredients := make([]ingredientsPackage.Ingredient, len(recipe.Ingredients))
		for idx_, ingredient := range recipe.Ingredients {
			ingredients[idx_] = ingredientsPackage.Ingredient{
				Id:   ingredient.Id,
				Name: ingredient.Name,
			}
		}

		missingIngredients := make([]ingredientsPackage.Ingredient, len(recipe.MissingIngredients))
		for idx_, ingredient := range recipe.MissingIngredients {
			missingIngredients[idx_] = ingredientsPackage.Ingredient{
				Id:   ingredient.Id,
				Name: ingredient.Name,
			}
		}

		parsedRecipes[idx] = recipes.Recipe{
			Id:                 recipe.Id,
			Title:              recipe.Title,
			Ingredients:        ingredients,
			MissingIngredients: missingIngredients,
		}
	}

	return parsedRecipes, nil
}

func (finder SpoonacularFinder) responseToNutrition(response *http.Response) (recipes.Nutrition, error) {
	type nutrition struct {
		Protein float32 `json:"percentProtein,omitempty"`
		Fat     float32 `json:"percentFat,omitempty"`
		Carbs   float32 `json:"percentCarbs,omitempty"`
	}
	type caloricBreakdownSection struct {
		Nutrition nutrition `json:"caloricBreakdown,omitempty"`
	}
	type nutritionSection struct {
		NutritionSection caloricBreakdownSection `json:"nutrition,omitempty"`
	}

	var nutritionFromJson nutritionSection
	err := json.NewDecoder(response.Body).Decode(&nutritionFromJson)
	if err != nil {
		return recipes.Nutrition{}, err
	}

	resultNutrition := recipes.Nutrition{
		Carbs:    nutritionFromJson.NutritionSection.Nutrition.Carbs,
		Proteins: nutritionFromJson.NutritionSection.Nutrition.Protein,
		Fat:      nutritionFromJson.NutritionSection.Nutrition.Fat,
	}

	return resultNutrition, nil
}
