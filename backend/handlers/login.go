package handlers

import (
	"fmt"
	"gotth/template/backend/auth"
	"gotth/template/backend/store"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	store := store.GetStore()
	session, err := store.Get(r, "session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound) //
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Determine the scheme (http or https)
	c := auth.GetCasdoor()

	// Construct base address from request
	baseAddress := "http://localhost:3000"

	redirect := baseAddress + "/callback"
	fmt.Println(redirect)

	loginUrl := c.Client.GetSigninUrl(redirect)

	http.Redirect(w, r, loginUrl, http.StatusFound)
}

func HandleLoginCallback(w http.ResponseWriter, r *http.Request) {
	c := auth.GetCasdoor()

	// Construct base address from request
	baseAddress := "http://localhost:3000"

	redirect := baseAddress + "/callback"

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
	}

	// Exchange the code for a token
	token, err := c.Client.GetOAuthToken(code, redirect)
	if err != nil {
		http.Error(w, "Failed to get token: "+err.Error(), http.StatusInternalServerError)
	}

	store := store.GetStore()
	session, err := store.Get(r, "session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	session.Values["token"] = token.AccessToken
	err = session.Save(r, w)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Redirect to the home page or another desired page
	http.Redirect(w, r, "/", http.StatusFound)
}
