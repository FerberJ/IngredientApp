package handlers

import (
	"gotth/template/backend/configuration"
	"gotth/template/backend/store"
	"gotth/template/backend/utils"
	"net/http"
	"slices"

	"github.com/go-chi/chi"
)

func HandleAddClosableBadge(w http.ResponseWriter, r *http.Request, cfg configuration.Configutration) {
	var valSlice []string
	s := store.GetStore()
	keyword := chi.URLParam(r, "keyword")

	val, err := s.GetValue("badgeList", w, r)
	if err == nil && val != nil {
		valSlice = val.([]string)
	}

	if !slices.Contains(valSlice, keyword) {
		valSlice = append(valSlice, keyword)
		s.AddValue("badgeList", valSlice, w, r)
	}

	HandleRecipes(w, r, cfg)
}

func HandleRemoveClosableBadge(w http.ResponseWriter, r *http.Request, cfg configuration.Configutration) {
	var valSlice []string
	s := store.GetStore()
	keyword := chi.URLParam(r, "keyword")

	val, err := s.GetValue("badgeList", w, r)
	if err == nil && val != nil {
		valSlice = val.([]string)

		valSlice = utils.RemoveString(valSlice, keyword)
		s.AddValue("badgeList", valSlice, w, r)
	}

	HandleRecipes(w, r, cfg)
}

func HandleRemoveAllClosableBadge(w http.ResponseWriter, r *http.Request, cfg configuration.Configutration) {
	var valSlice []string = make([]string, 0)
	s := store.GetStore()

	val, err := s.GetValue("badgeList", w, r)
	if err == nil && val != nil {
		s.AddValue("badgeList", valSlice, w, r)
	}

	HandleRecipes(w, r, cfg)
}
