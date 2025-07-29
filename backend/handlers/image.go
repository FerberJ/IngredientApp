package handlers

import (
	"gotth/template/backend/auth"
	"gotth/template/backend/dao"
	"io"
	"net/http"

	"github.com/go-chi/chi"
)

func HandleImageGet(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "image")
	id := ""
	user, _ := auth.GetUser(w, r)
	if user != nil {
		id = user.Id
	}

	img, _ := dao.GetImage(filename, id)

	if img == nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpg")
	if _, err := io.Copy(w, img); err != nil {
		return
	}
}
