package cache

import (
	"WhatCanIEat/WhatCanIEat/cmd/recipes"
	"WhatCanIEat/WhatCanIEat/cmd/recipes/finder"
)

type FinderCacher interface {
	Init(finder *finder.Finder) error
	Connect() error
	Disconnect() error
	FindByIngredientsNames(ingredients *[]string, numberOfRecipes int) ([]recipes.Recipe, error)
}
