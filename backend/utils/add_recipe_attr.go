package utils

import "github.com/a-h/templ"

// Returns attribuets
// change the verb of htmx
func GetAttr(edit bool, id string) templ.Attributes {
	if edit {
		return templ.Attributes{"hx-put": "/recipe/" + id}
	} else {
		return templ.Attributes{"hx-post": "/recipe"}
	}
}
