package cmd

import (
	"WhatCanIEat/WhatCanIEat/cmd/errors"
	"WhatCanIEat/WhatCanIEat/cmd/recipes"
	finderPackage "WhatCanIEat/WhatCanIEat/cmd/recipes/finder"
	finderCache "WhatCanIEat/WhatCanIEat/cmd/recipes/finder/cache"
	"github.com/spf13/cobra"
	"log"
)

var (
	ingredients     []string
	numberOfRecipes int
	finder          finderPackage.Finder
	cachedFinder    finderCache.FinderCacher

	rootCmd = &cobra.Command{
		Use:   "root --ingredients=tomatoes,eggs,pasta --numberOfRecipes=5",
		Short: "Generate numberOfRecipes of possible recipes using given ingredients.",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Connect to db.
			finder = finderPackage.SpoonacularFinder{}
			cachedFinder = &finderCache.Mssql{}
			err := cachedFinder.Init(&finder)
			if err != nil {
				log.Print("Cache disabled because cannot initialize cache, error: ", err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Received:\n")
			log.Printf("\tingredients: %v\n", ingredients)
			log.Printf("\tnumberOfRecipes: %v\n", numberOfRecipes)

			err := cachedFinder.Connect()
			if err != nil {
				log.Print("Cache disabled because cannot connect to cache, error: ", err)
				_ = cachedFinder.Disconnect()
			}

			possibleRecipes, err := cachedFinder.FindByIngredientsNames(&ingredients, numberOfRecipes)
			if err != nil {
				log.Fatalln(err)
			}

			recipes.DrawRecipes(&possibleRecipes)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			err := cachedFinder.Disconnect()
			if err != nil {
				log.Print("Failed to disconnect cache, error: ", err)
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			// TODO: Sanitize and check are given words ingredients.
			if numberOfRecipes < 1 {
				return errors.InvalidNumberOfIRecipes(numberOfRecipes)
			}
			if len(ingredients) < 1 {
				return errors.InvalidIngredients(ingredients)
			}
			return nil
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringSliceVar(&ingredients, "ingredients",
		nil,
		"List of ingredients to use. Example value: \"tomatoes,eggs,pasta\"")

	err := rootCmd.MarkFlagRequired("ingredients")
	if err != nil {
		log.Fatalln("Cannot set flag \"ingredients\" required, error: ", err)
	}

	rootCmd.Flags().IntVar(&numberOfRecipes, "numberOfRecipes",
		0,
		"Number of recipes to generate. Example value: \"5\"")

	err = rootCmd.MarkFlagRequired("numberOfRecipes")
	if err != nil {
		log.Fatalln("Cannot set flag \"numberOfRecipes\" required, error: ", err)
	}
}
