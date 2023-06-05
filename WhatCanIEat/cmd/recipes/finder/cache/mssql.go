package cache

import (
	recipesPackage "WhatCanIEat/WhatCanIEat/cmd/recipes"
	finderPackage "WhatCanIEat/WhatCanIEat/cmd/recipes/finder"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"log"
	"os"
	"sort"
	"strings"
)

type Mssql struct {
	db     *sql.DB
	finder *finderPackage.Finder
}

func (mssql *Mssql) Init(finder *finderPackage.Finder) error {
	mssql.finder = finder
	return nil
}

func (mssql *Mssql) Connect() error {
	// Create connection pool
	var err error

	server := os.Getenv("DB_SERVER")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		server, user, password, port, database)

	mssql.db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Println("Error creating connection pool: ", err)
		return err
	}

	ctx := context.Background()
	err = mssql.db.PingContext(ctx)
	if err != nil {
		return err
	}
	log.Printf("Connected to db.")

	err = mssql.ensureScheme()
	if err != nil {
		log.Println("Cannot ensure scheme, err: ", err)
		return err
	}
	log.Printf("DB scheme is valid.")

	return nil
}

// ensureScheme checks whether database contains required tables.
// When they don't - creates them.
// In case of error / conflict - returns error
// Assumes connected to the db and alive.
func (mssql *Mssql) ensureScheme() error {
	ctx := context.Background()

	// Recreate if not exists.
	tsql := `if not exists (select * from sysobjects where name='recipe' and xtype='U')
    create table recipe (
        ingredients VARCHAR(MAX) NOT NULL,
        recipes VARCHAR(MAX)
    )`

	_, err := mssql.db.ExecContext(ctx, tsql)
	if err != nil {
		return err
	}

	return nil
}

func (mssql *Mssql) Disconnect() error {
	if mssql.db != nil {
		err := mssql.db.Close()
		mssql.db = nil
		if err != nil {
			return err
		}
	}
	return nil
}

func (mssql *Mssql) FindByIngredientsNames(ingredients *[]string, numberOfRecipes int) ([]recipesPackage.Recipe, error) {
	var recipes []recipesPackage.Recipe = nil
	var err error

	if mssql.db != nil {
		recipes, err = mssql.searchCache(ingredients, numberOfRecipes)
		if err != nil {
			return nil, err
		}
	}

	// Not cached.
	if recipes == nil {
		recipes, err = (*mssql.finder).FindByIngredientsNames(ingredients, numberOfRecipes)
		if err != nil {
			return nil, err
		}
		if mssql.db != nil {
			err = mssql.saveIntoCache(ingredients, recipes)
			if err != nil {
				return nil, err
			}
		}
	}

	return recipes, nil
}

type recipesDbRow struct {
	ingredientsAsKey string
	recipesMarshaled []byte
}

func (mssql *Mssql) searchCache(ingredients *[]string, numberOfRecipes int) ([]recipesPackage.Recipe, error) {
	ingredientsAsKey := ingredeintsToKey(ingredients)
	ctx := context.Background()

	// Check if database is alive.
	err := mssql.db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	tsql := fmt.Sprintf("SELECT recipes FROM recipe WHERE ingredients='%s';", ingredientsAsKey)

	row := mssql.db.QueryRowContext(ctx, tsql)

	var recipesMarshaled []byte
	err = row.Scan(&recipesMarshaled)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var cachedRecipes []recipesPackage.Recipe
	err = json.Unmarshal(recipesMarshaled, &cachedRecipes)
	if err != nil {
		return nil, err
	}

	if len(cachedRecipes) < numberOfRecipes {
		return nil, nil
	}

	return cachedRecipes[:numberOfRecipes], nil
}

func (mssql *Mssql) saveIntoCache(ingredients *[]string, recipes []recipesPackage.Recipe) error {
	ctx := context.Background()

	// Check if database is alive.
	err := mssql.db.PingContext(ctx)
	if err != nil {
		return err
	}

	ingredientsAsKey := ingredeintsToKey(ingredients)

	// Delete if exists.
	tsql := fmt.Sprintf("DELETE FROM recipe WHERE ingredients='%s'", ingredientsAsKey)

	_, err = mssql.db.ExecContext(ctx,
		tsql,
		ingredientsAsKey)
	if err != nil {
		return err
	}

	// Insert.
	marshaled, err := json.Marshal(recipes)
	if err != nil {
		return err
	}
	tsql = fmt.Sprintf("INSERT INTO recipe (ingredients, recipes) VALUES ('%s', '%s');",
		ingredientsAsKey, marshaled)
	stmt, err := mssql.db.Prepare(tsql)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	row := stmt.QueryRowContext(ctx)

	var inserted []recipesDbRow
	err = row.Scan(&inserted)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	return nil
}

func ingredeintsToKey(ingredients *[]string) string {
	sorted := (*ingredients)[:]
	sort.Strings(sorted)
	return strings.Join(sorted, ",")
}
