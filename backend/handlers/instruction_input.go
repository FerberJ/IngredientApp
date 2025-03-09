package handlers

import (
	add_recipe_components "gotth/template/view/components/addRecipe"
	"net/http"
)

func HandleAddInstructionInput(w http.ResponseWriter, r *http.Request) {
	add_recipe_components.InstructionInput().Render(r.Context(), w)
}

func HandleRemoveInstructionInput(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
