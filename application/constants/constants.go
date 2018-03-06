package constants

/**
Class to stock all constants
 */

func GetPostgres() string {
	return "postgres"
}

func GetUser() string {
	return "user"
}

func GetTUser() string {
	return "t_user"
}

func GetHashedPassword() string {
	return "hashed_password"
}

func GetPassword() string {
	return "password"
}

func GetDbname() string {
	return "dbname"
}

func GetHost() string {
	return "host"
}

func GetPort() string {
	return "port"
}

func GetCharset() string {
	return "charset"
}

func GetRecipes() string {
	return "recipes"
}

func GetTRecipe() string {
	return "t_recipe"
}

func GetSlash() string {
	return "/"
}

func GetToken() string {
	return "token"
}

func GetDoubleDot() string {
	return ":"
}

func GetUniqueIdDbColumn() string {
	return "unique_id"
}

func GetUniqueId() string {
	return "uniqueId"
}

func GetRating() string {
	return "rating"
}

func GetLogin() string {
	return "login"
}

func GetUserId() string {
	return "userId"
}

func GetExpiration() string {
	return "rating"
}

func GetTRecipeRating() string {
	return "t_recipe_rating"
}

func GetRecipeUniqueId() string {
	return "recipe_unique_id"
}

func GetLimit() string {
	return "limit"
}

func GetOffset() string {
	return "offset"
}

func GetLimitDefaultValue() *uint64 {
	var t = uint64(20)
	return &t
}

func GetOffsetDefaultValue() *uint64 {
	var t = uint64(0)
	return &t
}

func GetQuestionMark() string {
	return "?"
}

func GetEquals() string {
	return "="
}

func GetAnd() string {
	return "&"
}

func GetName() string {
	return "name"
}

func GetPreparationTime() string {
	return "preparation_time"
}

func GetDifficulty() string {
	return "difficulty"
}

func GetVegetarian() string {
	return "vegetarian"
}

func GetCreatedBy() string {
	return "created_by"
}

func GetUpdateBy() string {
	return "updated_by"
}

func GetStar() string {
	return "*"
}

func GetAuthorization() string {
	return "Authorization"
}