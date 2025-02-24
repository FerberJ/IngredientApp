package repository

import "gotth/template/backend/db"

type RecipeRepository struct {
	*BaseRepository
}

func NewRecipeRepository(provider *db.MongoProvider) *RecipeRepository {
	return &RecipeRepository{
		BaseRepository: NewBaseRepository(provider, "recipes"),
	}
}
