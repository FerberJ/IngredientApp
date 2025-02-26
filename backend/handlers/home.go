package handlers

import (
	"fmt"
	"gotth/template/backend/auth"
	"gotth/template/view/home"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(w, r)
	if err != nil {
		home.Index("", false).Render(r.Context(), w)
		return
	}

	fmt.Println(user)

	home.Index(user.Avatar, true).Render(r.Context(), w)
}
