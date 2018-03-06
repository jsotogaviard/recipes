package handler

import (
	"net/http"
	"database/sql"
	"encoding/json"
	sq "github.com/Masterminds/squirrel"

	"jsotogaviard-api-test/application/model"
	"jsotogaviard-api-test/application/constants"
)

/**
Store the recipe rate in the database
 */
func RateRecipe(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// Get the recipe
	recipe := getRecipeOr404(db, w, r)
	if recipe == nil {
		return
	}

	// Get the rating
	rat := model.Rating{}
	decoder := json.NewDecoder(r.Body)
	if err1 := decoder.Decode(&rat); err1 != nil {
		respondError(w, http.StatusBadRequest, err1.Error())
		return
	}

	// Insert the rating
	_, err2 := sq.
		StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert(constants.GetTRecipeRating()).
		Columns(constants.GetRecipeUniqueId(), constants.GetRating()).
		Values(recipe.UniqueId, rat.Rating).
		RunWith(db).
		Query()

	if err2 != nil {
		respondError(w, http.StatusInternalServerError, err2.Error())
		return
	}

	// Send the response
	respondJSON(w, http.StatusNoContent, nil)
}