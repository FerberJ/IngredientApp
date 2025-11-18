package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"gotth/template/backend/auth"
	"gotth/template/backend/dao"
	"io"
	"net/http"
)

/*
func Export(w http.ResponseWriter, r *http.Request) {
	recipes, err := dao.ListRecipes(w, r, []string{}, []string{})
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)

	// Encode and send JSON response
	json.NewEncoder(w).Encode(recipes)
}
*/

func Export(w http.ResponseWriter, r *http.Request) {
	recipes, err := dao.ListRecipes(w, r, []string{}, []string{}, true)
	if err != nil {
		return
	}
	id := ""
	user, _ := auth.GetUser(w, r)
	if user != nil {
		id = user.Id
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, recipe := range recipes {
		recipeJSON, err := json.MarshalIndent(recipe, "", "  ")
		if err != nil {
			http.Error(w, "Error encoding recipe", http.StatusInternalServerError)
			return
		}

		folderName := fmt.Sprintf("%s", recipe.ID)
		filename := fmt.Sprintf("%s/data.json", folderName)

		f, err := zipWriter.Create(filename)
		if err != nil {
			http.Error(w, "Error creating ZIP file", http.StatusInternalServerError)
			return
		}

		_, err = f.Write(recipeJSON)
		if err != nil {
			http.Error(w, "Error writing to ZIP file", http.StatusInternalServerError)
			return
		}

		if recipe.Image != "" {
			imageName := fmt.Sprintf("%s/%s", folderName, recipe.Image)
			img, err := dao.GetImage(recipe.Image, id)
			if err != nil {
				http.Error(w, "Error getting image", http.StatusInternalServerError)
				return
			}

			imgFile, err := zipWriter.Create(imageName)
			if err != nil {
				http.Error(w, "Error creating ZIP file image", http.StatusInternalServerError)
				return
			}

			_, err = io.Copy(imgFile, img)
			if err != nil {
				http.Error(w, "Error writing image to ZIP", http.StatusInternalServerError)
				return
			}
		}
	}

	err = zipWriter.Close()
	if err != nil {
		http.Error(w, "Error closing ZIP file", http.StatusInternalServerError)
		return
	}

	// Set headers for ZIP download
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=recipes.zip")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

	// Write ZIP to response
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
