package handler

import (
	"net/http"
	sq "github.com/Masterminds/squirrel"
	"database/sql"
	"jsotogaviard-api-test/application/model"
	"strconv"
	"encoding/json"
	"jsotogaviard-api-test/application/constants"
	"jsotogaviard-api-test/application/serializer"
)

/**
Get recipes with paging
 */
func GetRecipes(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// Get limit or set default value
	params := r.URL.Query()
	var limit *uint64
	_, okLimit := params[constants.GetLimit()]
	if !okLimit {
		limit = constants.GetLimitDefaultValue()
	} else {
		limit = GetParamAsNumberOr400(constants.GetLimit(), w, r)
	}

	// Get offset or set default value
	var offset *uint64
	_, okOffset := params[constants.GetOffset()]
	if !okOffset {
		offset = constants.GetOffsetDefaultValue()
	} else {
		offset = GetParamAsNumberOr400(constants.GetOffset(), w, r)
	}

	// Get recipes from database
	recipes := make([]model.Recipe, 0)
	query := sq.
		StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select(constants.GetStar()).
		From(constants.GetTRecipe()).
		Limit(*limit).
		Offset(*offset)

	// Get parameter name
	name, okName := GetParam(constants.GetName(), r)
	if okName {
		query = query.Where(constants.GetName() + " LIKE ?", "%" + name + "%")
	}

	// Get parameter vegetarian
	vegetarianBool, okVegetarian := GetParam(constants.GetVegetarian(), r)
	if okVegetarian {
		vegetarian, err4 := strconv.ParseBool(vegetarianBool)
		if err4 != nil {
			respondError(w, http.StatusInternalServerError, err4.Error())
			return
		}
		query = query.Where(sq.Eq{constants.GetVegetarian(): vegetarian})
	}

	// Get parameter difficulty
	difficultyString, okDifficulty := GetParam(constants.GetDifficulty(), r)
	if okDifficulty {
		difficulty, err5 := strconv.Atoi(difficultyString)
		if err5 != nil {
			respondError(w, http.StatusInternalServerError, err5.Error())
			return
		}
		query = query.Where(sq.Eq{constants.GetDifficulty(): difficulty})
	}

	// Run the query in the database
	rows, err := query.
		RunWith(db).
		Query()
	defer rows.Close()

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Iterate over the result
	for rows.Next() {
		recipe :=  model.Recipe{}
		err1 := rows.Scan(
			&recipe.UniqueId,
			&recipe.Name,
			&recipe.PreparationTime,
			&recipe.Difficulty,
			&recipe.Vegetarian,
			&recipe.CreatedBy,
			&recipe.UpdatedBy)
		if err1 != nil {
			respondError(w, http.StatusInternalServerError, err1.Error())
			return
		}
		recipes = append(recipes, recipe)
	}

	err2 := rows.Err()
	if err2 != nil {
		respondError(w, http.StatusInternalServerError, err2.Error())
		return
	}

	// Compute next uri
	var nextUri string
	if uint64(len(recipes)) == *limit {
		nextLimit := strconv.FormatUint(*limit, 10)
		nextOffset := strconv.FormatUint(*offset + * limit, 10)
		nextUri = constants.GetSlash() + constants.GetRecipes() + constants.GetQuestionMark() +
			constants.GetLimit() + constants.GetEquals() + nextLimit + constants.GetAnd() +
			constants.GetOffset() + constants.GetEquals() + nextOffset

	} else {
		nextUri = ""
	}

	recipeSlice := model.RecipeSlice{
		Data : recipes,
		NextUri:  nextUri,
	}
	respondJSON(w, http.StatusOK, recipeSlice)
}

/**
Create the recipe
 */
func CreateRecipe(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// Authentication
	userId := GetUserOr401(w, r)
	if userId == nil {
		return
	}

	// Get recipe from body
	recipe := model.Recipe{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&recipe); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	// Insert the recipe in the database
	err := sq.
		StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert(constants.GetTRecipe()).
		Columns(constants.GetName(), constants.GetPreparationTime(), constants.GetDifficulty(),
		constants.GetVegetarian(), constants.GetCreatedBy()).
		Values(
		&recipe.Name,
		&recipe.PreparationTime,
		&recipe.Difficulty,
		&recipe.Vegetarian,
		userId).
		Suffix("RETURNING \"" + constants.GetUniqueIdDbColumn() + "\"").
		RunWith(db).
		QueryRow().
		Scan(&recipe.UniqueId)

	recipe.CreatedBy = *userId

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, recipe)
}

/**
Get a particular recipe
 */
func GetRecipe(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// Get the recipe from path
	recipe := getRecipeOr404(db, w, r)
	if recipe == nil {
		return
	}
	respondJSON(w, http.StatusOK, recipe)
}

/**
Update the recipe
 */
func UpdateRecipe(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// Authenticate
	userId := GetUserOr401(w, r)
	if userId == nil {
		return
	}

	// Get the recipe
	recipe := getRecipeOr404(db, w, r)
	if recipe == nil {
		return
	}

	// Decode the updated recipe
	newRecipe := model.Recipe{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newRecipe); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	// Update 'in place' in the database
	_, err := sq.
		StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update(constants.GetTRecipe()).
		Set(constants.GetName(), &newRecipe.Name).
		Set(constants.GetPreparationTime(), &newRecipe.PreparationTime).
		Set(constants.GetDifficulty(), &newRecipe.Difficulty).
		Set(constants.GetVegetarian(), &newRecipe.Vegetarian).
		Set(constants.GetUpdateBy(), userId).
		Where(sq.Eq{constants.GetUniqueIdDbColumn(): &recipe.UniqueId}).
		RunWith(db).
		Query()
	newRecipe.UniqueId = recipe.UniqueId
	newRecipe.UpdatedBy = serializer.JsonNullFloat64{sql.NullFloat64{Valid:true, Float64:*userId}}

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, newRecipe)
}

/**
Delete the recipe
 */
func DeleteRecipe(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// Authenticate
	userId := GetUserOr401(w, r)
	if userId == nil {
		return
	}

	// Get the recipe
	recipe := getRecipeOr404(db, w, r)
	if recipe == nil {
		return
	}

	// Delete from database
	_, err := sq.
		StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Delete(constants.GetTRecipe()).
		Where(sq.Eq{constants.GetUniqueIdDbColumn(): recipe.UniqueId}).
		RunWith(db).
		Exec()

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}
