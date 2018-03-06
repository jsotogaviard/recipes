package application

import (
	"fmt"
	"log"
	"net/http"
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"jsotogaviard-api-test/application/handler"
	"jsotogaviard-api-test/application/security"
	"jsotogaviard-api-test/configuration"
	"jsotogaviard-api-test/application/constants"
	"jsotogaviard-api-test/application/model"
	"github.com/rubenv/sql-migrate"
)

// Application is a router and a database
type Application struct {
	Router   *mux.Router
	Database *sql.DB
}

/**
Initialize the application
 - init the db connection
 - run the migration
 - init the router
 - set the routes
 */
func GetApplication(c *configuration.Config, s *security.Security) (*Application) {
	dbURI := fmt.Sprintf(
		constants.GetHost() + "=%s " +
			constants.GetUser() + "=%s " +
			"sslmode=disable " +
			constants.GetPassword() + "=%s " +
			constants.GetDbname() + "=%s",
		c.Database.Host,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName)

	db, err := sql.Open(constants.GetPostgres(), dbURI)
	if err != nil {
		log.Fatal("Could not connect database : %v", err.Error())
	}

	model.MigrateDatabase(db, migrate.Up)

	a := Application{
		Database: db,
		Router :  mux.NewRouter(),
	}
	a.setRouters(s)
	return &a
}

/**
Set the routes
 */
func (a *Application) setRouters(s *security.Security) {

	// Security token
	a.Post(constants.GetSlash() + constants.GetToken(), http.HandlerFunc(a.GetToken)) // Not protected

	// Recipe endpoints
	// Not protected
	a.Get(constants.GetSlash() + constants.GetRecipes() ,http.HandlerFunc(a.GetRecipes))
	a.Get(constants.GetSlash() + constants.GetRecipes() + constants.GetSlash() + "{" + constants.GetUniqueId() + "}", http.HandlerFunc(a.GetRecipe))

	// Protected
	a.Post(constants.GetSlash() + constants.GetRecipes(), s.Jwt.Handler(http.HandlerFunc(a.CreateRecipe)))
	a.Put(constants.GetSlash() + constants.GetRecipes() + constants.GetSlash() + "{" + constants.GetUniqueId() + "}", s.Jwt.Handler(http.HandlerFunc(a.UpdateRecipe)))
	a.Delete(constants.GetSlash() + constants.GetRecipes() + constants.GetSlash() + "{" + constants.GetUniqueId() + "}", s.Jwt.Handler(http.HandlerFunc(a.DeleteRecipe)))

	// Recipe rating
	a.Post(constants.GetSlash() + constants.GetRecipes() + constants.GetSlash() + "{" + constants.GetUniqueId() + "}" + constants.GetSlash() + constants.GetRating(), http.HandlerFunc(a.RateRecipe)) // Not protected

}

func (a *Application) Get(path string, handler http.Handler) {
	a.Router.Handle(path, handler).Methods(http.MethodGet)
}

func (a *Application) Post(path string, handler http.Handler) {
	a.Router.Handle(path, handler).Methods(http.MethodPost)
}

func (a *Application) Put(path string, handler http.Handler) {
	a.Router.Handle(path, handler).Methods(http.MethodPut)
}

func (a *Application) Delete(path string, handler http.Handler) {
	a.Router.Handle(path, handler).Methods(http.MethodDelete)
}

/*
** Recipe Handlers
 */
func (a *Application) GetToken(w http.ResponseWriter, r *http.Request) {
	handler.GetToken(a.Database, w, r)
}

func (a *Application) GetRecipes(w http.ResponseWriter, r *http.Request) {
	handler.GetRecipes(a.Database, w, r)
}

func (a *Application) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	handler.CreateRecipe(a.Database, w, r)
}

func (a *Application) GetRecipe(w http.ResponseWriter, r *http.Request) {
	handler.GetRecipe(a.Database, w, r)
}

func (a *Application) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	handler.UpdateRecipe(a.Database, w, r)
}

func (a *Application) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	handler.DeleteRecipe(a.Database, w, r)
}

func (a *Application) RateRecipe(w http.ResponseWriter, r *http.Request) {
	handler.RateRecipe(a.Database, w, r)
}

// Run server
func (a *Application) Run(host string) {
	log.Print("Starting server")
	log.Fatal(http.ListenAndServe(host, a.Router))
}
