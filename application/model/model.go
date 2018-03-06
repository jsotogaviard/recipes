package model

import (
	"database/sql"
	"github.com/rubenv/sql-migrate"
	"log"

	"jsotogaviard-api-test/application/constants"
	"jsotogaviard-api-test/application/serializer"
)

// User structure used for token retrieval
type User struct {
	Login          string `json:"login"`
	HashedPassword string `json:"hashedPassword"`
}

// Token retrieved
type Token struct {
	Token 	 string  `json:"token"`
}

// Rating used to rate a recipe
type Rating struct {
	Rating 	 int  `json:"rating"`
}

// Recipe structure
type Recipe struct {
	UniqueId        float64         `json:"uniqueId"`
	Name            string          `json:"name"`
	PreparationTime int32           `json:"preparationTime"`
	Difficulty      int32           `json:"difficulty"`
	Vegetarian      bool            `json:"vegetarian"`
	CreatedBy       float64         `json:"createdBy"`
	UpdatedBy       serializer.JsonNullFloat64 `json:"updatedBy"`
}

// Recipe slice for paged response
type RecipeSlice struct {
	Data []Recipe `json:"data"`
	NextUri string `json:"nextUri"`
}

// Migrate the database
func MigrateDatabase(db *sql.DB, action migrate.MigrationDirection){
	migrations := &migrate.FileMigrationSource{Dir: "./db"}

	_, err := migrate.Exec(db, constants.GetPostgres(), migrations, action)
	if err != nil {
		log.Fatal("Could not run migration ", err)
	}
}
