# WhatCanIEat
Being hungry while the fridge is full? </br>
Provide me a list of ingredients you have and cook something together!

## Description
This app takes ingredients and number of recipes to generate.</br>
Then fetches them from [Spoonacular](https://spoonacular.com/food-api) api.


## Setup
Env variable storing api key: **SPOONACULAR_API_KEY** 

## Usage
CLI flags:  </br>
- **--ingredients** 
- **--numberOfRecipes**

Example `--ingredients=tomatoes,eggs,pasta --numberOfRecipes=5`

## Example result
![](res/example_result.jpg)

## Dependencies
[cobra](https://github.com/spf13/cobra)