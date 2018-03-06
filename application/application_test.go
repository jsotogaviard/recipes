package application_test

import (
	//"bytes"
	//"encoding/json"
	//"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	//"jsotogaviard-api-test/app/handler"
	"jsotogaviard-api-test/application"
	"jsotogaviard-api-test/configuration"
	"encoding/json"
	"bytes"
	"log"
	"jsotogaviard-api-test/application/security"
	"jsotogaviard-api-test/application/model"
	"github.com/rubenv/sql-migrate"
	"jsotogaviard-api-test/application/constants"
)

var a *application.Application

func TestMain(m *testing.M) {

	c := configuration.GetConfig()

	s := security.GetSecurity()

	a = application.GetApplication(c, s)

	DropDatabase()

	code := m.Run()

	os.Exit(code)
}


func TestGetToken(t *testing.T) {
	CreateDatabase()

	login(t)

	DropDatabase()
}

func TestSearchRecipes(t *testing.T) {

	CreateDatabase()

	addRecipes(20)

	req, _ := http.NewRequest(http.MethodGet, constants.GetSlash() + constants.GetRecipes() + "?name=name&limit=5", nil)
	r := executeRequest(req)
	checkResponseCode(t, http.StatusOK, r.Code)

	recipeSlice := model.RecipeSlice{}
	decoder := json.NewDecoder(r.Body)
	if err0 := decoder.Decode(&recipeSlice); err0 != nil {
		t.Errorf("Cannot decode in recipe slice: %v", r.Body)
	}
	if len(recipeSlice.Data) != 5 {
		t.Errorf("No data should be there")
	}

	if recipeSlice.NextUri != "/recipes?limit=5&offset=5" {
		t.Errorf("No next uri")
	}

	DropDatabase()
}

func TestSearchRecipesVegetarian(t *testing.T) {

	CreateDatabase()

	addRecipes(20)

	req, _ := http.NewRequest(http.MethodGet, constants.GetSlash() + constants.GetRecipes() + "?vegetarian=true&name=name&limit=5&offset=4", nil)
	r := executeRequest(req)

	checkResponseCode(t, http.StatusOK, r.Code)

	recipeSlice := model.RecipeSlice{}
	decoder := json.NewDecoder(r.Body)
	if err0 := decoder.Decode(&recipeSlice); err0 != nil {
		t.Errorf("Cannot decode in recipe slice: %v", r.Body)
	}
	if len(recipeSlice.Data) != 5 {
		t.Errorf("No data should be there")
	}

	if recipeSlice.NextUri != "/recipes?limit=5&offset=9" {
		t.Errorf("No next uri")
	}

	DropDatabase()
}

func TestSearchRecipesDifficulty(t *testing.T) {

	CreateDatabase()

	addRecipes(20)

	req, _ := http.NewRequest(http.MethodGet, constants.GetSlash() + constants.GetRecipes() + "?vegetarian=true&name=name&limit=5&offset=4&difficulty=1", nil)
	r := executeRequest(req)

	checkResponseCode(t, http.StatusOK, r.Code)

	recipeSlice := model.RecipeSlice{}
	decoder := json.NewDecoder(r.Body)
	if err0 := decoder.Decode(&recipeSlice); err0 != nil {
		t.Errorf("Cannot decode in recipe slice: %v", r.Body)
	}
	if len(recipeSlice.Data) != 5 {
		t.Errorf("No data should be there")
	}

	if recipeSlice.NextUri != "/recipes?limit=5&offset=9" {
		t.Errorf("No next uri")
	}

	DropDatabase()
}

func TestSearchRecipesDifficulty1(t *testing.T) {

	CreateDatabase()

	addRecipes(20)

	req, _ := http.NewRequest(http.MethodGet, constants.GetSlash() + constants.GetRecipes() + "?vegetarian=true&name=name&limit=5&offset=4&difficulty=2", nil)
	r := executeRequest(req)

	checkResponseCode(t, http.StatusOK, r.Code)

	recipeSlice := model.RecipeSlice{}
	decoder := json.NewDecoder(r.Body)
	if err0 := decoder.Decode(&recipeSlice); err0 != nil {
		t.Errorf("Cannot decode in recipe slice: %v", r.Body)
	}
	if len(recipeSlice.Data) != 0 {
		t.Errorf("No data should be there")
	}

	if recipeSlice.NextUri != "" {
		t.Errorf("No next uri")
	}

	DropDatabase()
}

func TestGetRecipesEmpty(t *testing.T) {

	CreateDatabase()

	req, _ := http.NewRequest(http.MethodGet, constants.GetSlash() + constants.GetRecipes(), nil)
	r := executeRequest(req)

	checkResponseCode(t, http.StatusOK, r.Code)

	recipeSlice := model.RecipeSlice{}
	decoder := json.NewDecoder(r.Body)
	if err0 := decoder.Decode(&recipeSlice); err0 != nil {
		t.Errorf("Cannot decode in recipe slice: %v", r.Body)
	}

	if len(recipeSlice.Data) != 0 {
		t.Errorf("No data should be there")
	}

	if recipeSlice.NextUri != "" {
		t.Errorf("No next uri")
	}

	DropDatabase()
}

func TestGetRecipesNotEmpty(t *testing.T) {

	CreateDatabase()

	addRecipes(22)

	req, _ := http.NewRequest(http.MethodGet, constants.GetSlash() + constants.GetRecipes(), nil)
	r := executeRequest(req)

	checkResponseCode(t, http.StatusOK, r.Code)

	recipeSlice := model.RecipeSlice{}
	decoder := json.NewDecoder(r.Body)
	if err0 := decoder.Decode(&recipeSlice); err0 != nil {
		t.Errorf("Cannot decode in recipe slice: %v", r.Body)
	}

	if len(recipeSlice.Data) != 20 {
		t.Errorf("No data should be there")
	}

	if recipeSlice.NextUri == "" {
		t.Errorf("next uri")
	}

	DropDatabase()
}

func TestGetNotExistentRecipe(t *testing.T) {

	CreateDatabase()

	addRecipes(1)

	req, _ := http.NewRequest(http.MethodGet, constants.GetSlash() + constants.GetRecipes() + constants.GetSlash() + "1", nil)
	r := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, r.Code)

	DropDatabase()

}

func TestGetExistentRecipe(t *testing.T) {

	CreateDatabase()

	addRecipes(1)

	req, _ := http.NewRequest(http.MethodGet, constants.GetSlash() + constants.GetRecipes() + constants.GetSlash() + "0", nil)
	r := executeRequest(req)

	checkResponseCode(t, http.StatusOK, r.Code)

	recipe := model.Recipe{}
	decoder := json.NewDecoder(r.Body)
	if err0 := decoder.Decode(&recipe); err0 != nil {
		t.Errorf("Cannot decode in recipe: %v", err0.Error())
	}

	DropDatabase()
}

func TestCreateRecipe(t *testing.T) {

	CreateDatabase()

	token := login(t)

	recipe := model.Recipe{
		Name: "name",
		PreparationTime: 1,
		Difficulty:2,
		Vegetarian:true,
		}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(recipe)

	req, _ := http.NewRequest(http.MethodPost,  constants.GetSlash() + constants.GetRecipes() , b)
	req.Header.Set(constants.GetAuthorization(), "Bearer " + token)
	r := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, r.Code)

	newRecipe := model.Recipe{}
	decoder := json.NewDecoder(r.Body)
	if err0 := decoder.Decode(&newRecipe); err0 != nil {
		t.Errorf("Cannot decode in recipe: %v", r.Body)
	}

	b1 := new(bytes.Buffer)
	json.NewEncoder(b1).Encode(recipe)
	req1, _ := http.NewRequest(http.MethodPost,  constants.GetSlash() + constants.GetRecipes() , b1)
	req1.Header.Set(constants.GetAuthorization(), "Bearer " + token)
	r1 := executeRequest(req1)

	checkResponseCode(t, http.StatusCreated, r1.Code)

	newRecipe1 := model.Recipe{}
	decoder1 := json.NewDecoder(r1.Body)
	if err0 := decoder1.Decode(&newRecipe1); err0 != nil {
		t.Errorf("Cannot decode in recipe: %v", r.Body)
	}

	DropDatabase()

}

func TestRateRecipe(t *testing.T) {

	CreateDatabase()

	addRecipes(1)

	recipeRating := model.Rating{
		Rating: 1,
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(recipeRating)

	uri := constants.GetSlash() + constants.GetRecipes() + constants.GetSlash() + "0" + constants.GetSlash() + constants.GetRating()
	req, _ := http.NewRequest(http.MethodPost, uri, b)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNoContent, response.Code)

	DropDatabase()
}

func TestUpdateRecipe(t *testing.T) {

	CreateDatabase()

	token := login(t)

	recipe := model.Recipe{
		Name: "name",
		PreparationTime: 1,
		Difficulty:2,
		Vegetarian:true,
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(recipe)

	req, _ := http.NewRequest(http.MethodPost,  constants.GetSlash() + constants.GetRecipes() , b)
	req.Header.Set(constants.GetAuthorization(), "Bearer " + token)
	r := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, r.Code)

	newRecipe := model.Recipe{}
	decoder := json.NewDecoder(r.Body)
	if err0 := decoder.Decode(&newRecipe); err0 != nil {
		t.Errorf("Cannot decode in recipe: %v", r.Body)
	}

	recipeToUpdate := model.Recipe{
		Name: "name2",
		PreparationTime: 1,
		Difficulty:2,
		Vegetarian:true,
	}
	b1 := new(bytes.Buffer)
	json.NewEncoder(b1).Encode(recipeToUpdate)

	id := strconv.Itoa(int(newRecipe.UniqueId))
	req1, _ := http.NewRequest(http.MethodPut,  constants.GetSlash() + constants.GetRecipes() + constants.GetSlash() + id , b1)
	req1.Header.Set(constants.GetAuthorization(), "Bearer " + token)
	r1 := executeRequest(req1)

	checkResponseCode(t, http.StatusOK, r1.Code)

	newRecipe1 := model.Recipe{}
	decoder1 := json.NewDecoder(r1.Body)
	if err1 := decoder1.Decode(&newRecipe1); err1 != nil {
		t.Errorf("Cannot decode in recipe: %v", r.Body)
	}

	if(newRecipe1.Name != "name2"){
		t.Errorf("Wrong update")
	}

	DropDatabase()

}

func TestDeleteRecipe(t *testing.T) {

	CreateDatabase()

	token := login(t)

	addRecipes(1)

	req, _ := http.NewRequest(http.MethodGet, constants.GetSlash() + constants.GetRecipes() + constants.GetSlash()  + "0", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest(http.MethodDelete, constants.GetSlash() + constants.GetRecipes() + constants.GetSlash() + "0", nil)
	req.Header.Set(constants.GetAuthorization(), "Bearer " + token)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusNoContent, response.Code)

	req, _ = http.NewRequest(http.MethodGet,  constants.GetSlash() + constants.GetRecipes() + "0", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	DropDatabase()
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// Database utility
func CreateDatabase(){
	model.MigrateDatabase(a.Database, migrate.Up)
}

func DropDatabase(){
	model.MigrateDatabase(a.Database, migrate.Down)
}

// Data utilities
func login(t *testing.T) string {

	login := "login"
	hashedPassword := "hashedPassword"

	addUser(login, hashedPassword)

	u := model.User{Login: login, HashedPassword: hashedPassword}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)

	req, _ := http.NewRequest(http.MethodPost, constants.GetSlash() + constants.GetToken(), b)
	r := executeRequest(req)

	checkResponseCode(t, http.StatusOK, r.Code)

	token := model.Token{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&token); err != nil {
		t.Errorf("Cannot decode token")
	}

	return token.Token
}


func addUser(login string, hashedPassword string) {
	a.Database.Exec("INSERT INTO t_user(login, hashed_password) VALUES($1, $2)", login, hashedPassword)
}

func addRecipes(count int) {
	addUser("a", "a")
	vegetarian := true
	for i := 0; i < count; i++ {
		if _, err := a.Database.Exec("INSERT INTO t_recipe VALUES($1, $2,$3, $4, $5, $6, $7)", i, "name_" + strconv.Itoa(i), 3, 1, vegetarian, 1, 1); err != nil{
			log.Fatal(err)
		}
		if _, err := a.Database.Exec("INSERT INTO t_recipe_rating VALUES($1, $2, $3)", i, i , 3); err != nil{
			log.Fatal(err)
		}
		vegetarian =! vegetarian
	}
}
