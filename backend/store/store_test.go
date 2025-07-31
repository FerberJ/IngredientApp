package store_test

import (
	"fmt"
	"gotth/template/backend/configuration"
	"gotth/template/backend/db"
	"gotth/template/backend/models"
	"gotth/template/backend/store"
	"testing"
)

func TestRecipe(t *testing.T) {
	rec := models.Recipe{
		Name: "TEST",
	}

	db.NewBadgerProvider(configuration.Configutration{})
	store.InitStore()
	s := store.GetStore()

	s.AddTempRecipe("abcdef", rec)

	newRec, err := s.GetTempRecipe("abcdef")

	fmt.Println(newRec, err)
}
