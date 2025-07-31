package handlers

import (
	"encoding/json"
	"gotth/template/backend/dao"
	"gotth/template/backend/store"
	"gotth/template/view/components"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func HandlePrepareBringRequest(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")

	recipe, err := dao.GetRecipe(w, r, recipeID, true)
	if err != nil {
		return
	}

	url, err := dao.GetImageTTL(recipe.Image)
	if err != nil {
		return
	}

	recipe.Image = url.String()

	s := store.GetStore()

	newId := uuid.NewString()

	s.AddTempRecipe(newId, recipe)

	returnValue := make(map[string]any)
	returnValue["id"] = newId
	returnValueByte, err := json.Marshal(returnValue)

	w.Write(returnValueByte)
}

func HandleBringRequest(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")
	servingSiceStr := r.URL.Query().Get("serving")

	s := store.GetStore()
	recipe, err := s.GetTempRecipe(recipeID)
	if err != nil {
		return
	}

	// Get the servingSice. The Serving sice can not be below 1
	servingSice := recipe.Nutrition.ServingSize
	if servingSiceStr != "" {
		num, err := strconv.ParseInt(servingSiceStr, 10, 64)
		if err == nil {
			if num >= 1 {
				servingSice = int(num)
			}
		}
	}

	dao.GetServing(servingSice, &recipe)

	components.BringRecipe(recipe).Render(r.Context(), w)
}
