package handlers

import (
	"gotth/template/backend/auth"
	"gotth/template/backend/configuration"
	"gotth/template/backend/store"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	s := store.GetStore()
	s.DeleteToken(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}

func HandleLogin(w http.ResponseWriter, r *http.Request, cfg configuration.Configutration) {
	c := auth.GetCasdoor()
	loginUrl := c.Client.GetSigninUrl(cfg.CallbackAddress)
	http.Redirect(w, r, loginUrl, http.StatusFound)
}

func HandleLoginCallback(w http.ResponseWriter, r *http.Request, cfg configuration.Configutration) {
	c := auth.GetCasdoor()
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
	}

	token, err := c.Client.GetOAuthToken(code, cfg.CallbackAddress)
	if err != nil {
		http.Error(w, "Failed to get token: "+err.Error(), http.StatusInternalServerError)
	}

	s := store.GetStore()
	s.SaveToken(token.AccessToken, w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
