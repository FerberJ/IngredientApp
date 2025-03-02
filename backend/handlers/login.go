package handlers

import (
	"gotth/template/backend/auth"
	"gotth/template/backend/store"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	s := store.GetStore()
	s.DeleteToken(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	c := auth.GetCasdoor()
	baseAddress := "http://localhost:3000" // TODO: add to config
	redirect := baseAddress + "/callback"
	loginUrl := c.Client.GetSigninUrl(redirect)
	http.Redirect(w, r, loginUrl, http.StatusFound)
}

func HandleLoginCallback(w http.ResponseWriter, r *http.Request) {
	c := auth.GetCasdoor()
	baseAddress := "http://localhost:3000" // TODO: add to config
	redirect := baseAddress + "/callback"
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
	}

	token, err := c.Client.GetOAuthToken(code, redirect)
	if err != nil {
		http.Error(w, "Failed to get token: "+err.Error(), http.StatusInternalServerError)
	}

	s := store.GetStore()
	s.SaveToken(token.AccessToken, w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
