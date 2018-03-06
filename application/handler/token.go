package handler

import (
	"net/http"
	"database/sql"
	"encoding/json"
	sq "github.com/Masterminds/squirrel"

	"jsotogaviard-api-test/application/security"
	"jsotogaviard-api-test/application/constants"
	"jsotogaviard-api-test/application/model"
)

/**
Get the token
 */
func GetToken(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// Decode the payload
	userBody := model.User{}
	decoder := json.NewDecoder(r.Body)
	if err0 := decoder.Decode(&userBody); err0 != nil {
		respondError(w, http.StatusBadRequest, err0.Error())
		return
	}
	defer r.Body.Close()

	// Check user exists and has correct hashed password
	var userId int
	err1 := sq.
		StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select(constants.GetUniqueIdDbColumn()).
		From(constants.GetTUser()).
		Where(sq.Eq{
			constants.GetLogin():          userBody.Login,
			constants.GetHashedPassword(): userBody.HashedPassword}).
		RunWith(db).
		QueryRow().
		Scan(&userId)
	if err1 != nil {
		respondError(w, http.StatusInternalServerError, err1.Error())
		return
	}

	// Compute the token
	tokenString, err2 := security.ComputeToken(userBody.Login, userId)
	if err2 != nil {
		respondError(w, http.StatusInternalServerError, err2.Error())
		return
	}

	// Provide the token
	respondJSON(w, http.StatusOK, model.Token{Token:tokenString})
}