package cmd

import (
	"WhatCanIEat/WhatCanIEat/cmd/errors"
	"WhatCanIEat/WhatCanIEat/cmd/recipes"
	finderPackage "WhatCanIEat/WhatCanIEat/cmd/recipes/finder"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var (
	ingredients     *[]string
	numberOfRecipes int

	rootCmd = &cobra.Command{
		Use:   "root --ingredients=tomatoes,eggs,pasta --numberOfRecipes=5",
		Short: "Generate numberOfRecipes of possible recipes using given ingredients.",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Connect to db.
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside rootCmd Run with args:\n")
			fmt.Printf("\tingredients: %v\n", *ingredients)
			fmt.Printf("\tnumberOfRecipes: %v\n", numberOfRecipes)

			var finder finderPackage.Finder = finderPackage.SpoonacularFinder{}

			possibleRecipes, err := finder.FindByIngredientsNames(ingredients, numberOfRecipes)
			if err != nil {
				log.Fatalln(err)
			}

			recipes.DrawRecipes(&possibleRecipes)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			// Disconnect from db.
		},
		Args: func(cmd *cobra.Command, args []string) error {
			// TODO: Check are given words ingredients.
			if numberOfRecipes < 1 {
				return errors.InvalidNumberOfIRecipes(numberOfRecipes)
			}
			if len(*ingredients) < 1 {
				return errors.InvalidIngredients(*ingredients)
			}
			return nil
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	ingredients = rootCmd.Flags().StringSlice("ingredients",
		nil,
		"List of ingredients to use. Example value: \"tomatoes,eggs,pasta\"")

	err := rootCmd.MarkFlagRequired("ingredients")
	if err != nil {
		return
	}

	rootCmd.Flags().IntVar(&numberOfRecipes, "numberOfRecipes",
		0,
		"Number of recipes to generate. Example value: \"5\"")

	err = rootCmd.MarkFlagRequired("numberOfRecipes")
	if err != nil {
		return
	}
}
