package finder

import (
	"WhatCanIEat/WhatCanIEat/cmd/recipes"
)

type Finder interface {
	FindByIngredientsNames(ingredients *[]string, numberOfRecipes int) ([]recipes.Recipe, error)
}
