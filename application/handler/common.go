package handler

import (
	"encoding/json"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"github.com/gorilla/mux"
	"strconv"
	"jsotogaviard-api-test/application/constants"
	"database/sql"
	"jsotogaviard-api-test/application/model"
	"jsotogaviard-api-test/application/security"
	"strings"
)

/**
Answer with json and status
 */
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

/**
Answer with error
 */
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

/**
Get user from authorization header or send unauthorized
 */
func GetUserOr401(w http.ResponseWriter, r *http.Request) (*float64){
	tokenString := r.Header.Get(constants.GetAuthorization())
	tokenString = strings.Split(tokenString, " ")[1];
	userId, err := security.GetUserId(tokenString)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return nil
	} else {
		return userId
	}

}

/**
Get a value from the path in the query or send a bad request
 */
func GetVarAsNumberOr400(key string, w http.ResponseWriter, r *http.Request) (*uint64) {
	vars := mux.Vars(r)
	valueString := vars[key]
	val, err := strconv.ParseUint(valueString, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return nil
	} else {
		return &val
	}
}

/**
Get a parameter with a key
 */
func GetParam(key string, r *http.Request) (string, bool) {
	vars :=  r.URL.Query()
	valueString, ok := vars[key]
	if ok {
		return valueString[0], ok
	} else {
		return "", ok
	}
}

/**
Get a parameter with a key or send a bad request
 */
func GetParamAsNumberOr400(key string, w http.ResponseWriter, r *http.Request) (*uint64) {
	valString, ok := GetParam(key, r)
	val, err := strconv.ParseUint(valString, 10, 64)
	if !ok || err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return nil
	} else {
		return &val
	}
}

/**
Get a recipe or send a not found exception
 */
func getRecipeOr404(db *sql.DB, w http.ResponseWriter, r *http.Request) (*model.Recipe) {

	// Get unique id
	uniqueId := GetVarAsNumberOr400(constants.GetUniqueId(), w, r)
	if uniqueId == nil {
		return nil
	}

	// Get recipe from database
	recipe := model.Recipe{}
	err := sq.
		StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select(constants.GetStar()).
		From(constants.GetTRecipe()).
		Where(sq.Eq{constants.GetUniqueIdDbColumn(): uniqueId}).
		RunWith(db).
		QueryRow().
		Scan(
		&recipe.UniqueId,
		&recipe.Name,
		&recipe.PreparationTime,
		&recipe.Difficulty,
		&recipe.Vegetarian,
		&recipe.CreatedBy,
		&recipe.UpdatedBy)

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	} else {
		return &recipe
	}
}
