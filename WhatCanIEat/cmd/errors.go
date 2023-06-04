package cmd

import "fmt"

type InvalidNumberOfIRecipes int

func (e InvalidNumberOfIRecipes) Error() string {
	return fmt.Sprintf(" \"%d\" is not a valid number of ingredients.", int(e))
}

type InvalidIngredients []string

func (e InvalidIngredients) Error() string {
	return fmt.Sprintf(" \"%s\" are not valid ingredients.", []string(e))
}
