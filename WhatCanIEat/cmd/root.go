package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	ingredients     *[]string
	numberOfRecipes *int

	rootCmd = &cobra.Command{
		Use:   "root --ingredients=tomatoes,eggs,pasta --numberOfRecipes=5",
		Short: "Generate numberOfRecipes of possible recipes using given ingredients.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			//fmt.Printf("Inside rootCmd PersistentPreRun with args: %v\n", args)
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			//fmt.Printf("Inside rootCmd PreRun with args: %v\n", args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside rootCmd Run with args:\n")
			fmt.Printf("\tingredients: %v\n", *ingredients)
			fmt.Printf("\tnumberOfRecipes: %v\n", *numberOfRecipes)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			//fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			//fmt.Printf("Inside rootCmd PersistentPostRun with args: %v\n", args)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			// TODO: Check are given words ingredients.
			if *numberOfRecipes < 1 {
				return InvalidNumberOfIRecipes(*numberOfRecipes)
			}
			if len(*ingredients) < 1 {
				return InvalidIngredients(*ingredients)
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

	numberOfRecipes = rootCmd.Flags().Int("numberOfRecipes",
		0,
		"Number of recipes to generate. Example value: \"5\"")

	err = rootCmd.MarkFlagRequired("numberOfRecipes")
	if err != nil {
		return
	}
}
