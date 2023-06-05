package recipes

import (
	"WhatCanIEat/WhatCanIEat/cmd/ingredients"
	"fmt"
)

type Recipe struct {
	Id                 int
	Title              string
	Ingredients        []ingredients.Ingredient
	MissingIngredients []ingredients.Ingredient
	Nutrition          Nutrition
}

type Nutrition struct {
	Carbs    float32 // percent
	Proteins float32 // percent
	Fat      float32 // percent
}

func DrawRecipes(recipes *[]Recipe) {
	out := ""
	for _, recipe := range *recipes {
		out +=
			"Recipe name: " + recipe.Title + "\n"
		if len(recipe.Ingredients) > 0 {
			out += "  Ingredients: \n"
			for _, ingredient := range recipe.Ingredients {
				out += "    " + ingredient.Name + "\n"
			}
		}
		if len(recipe.MissingIngredients) > 0 {
			out += "  Ingredients (missing): \n"
			for _, ingredient := range recipe.MissingIngredients {
				out += "    " + ingredient.Name + "\n"
			}
		}
		out += fmt.Sprintf(
			"  Nutrition: \n"+
				"    Carbs: %.1f%%\n"+
				"    Proteins: %.1f%%\n"+
				"    Fat: %.1f%%\n",
			recipe.Nutrition.Carbs, recipe.Nutrition.Proteins, recipe.Nutrition.Fat)
		out += "\n"
	}
	fmt.Println(out)
}
