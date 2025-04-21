package handlers

import (
	"gotth/template/backend/auth"
	"gotth/template/view/home"
	"net/http"
)

func HandleListPage(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(w, r)
	if err != nil {
		home.ListIndex("", false).Render(r.Context(), w)
		return
	}

	home.ListIndex(user.Name, true).Render(r.Context(), w)
}
