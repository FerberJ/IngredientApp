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
		return
	}

	token, err := c.Client.GetOAuthToken(code, cfg.CallbackAddress)
	if err != nil {
		return
	}

	s := store.GetStore()

	// Add error handling for SaveToken
	err = s.SaveToken(token.AccessToken, w, r)
	if err != nil {
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
