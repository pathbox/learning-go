package controller

import (
	"app/shared/view"
	"net/http"
)

func AboutGET(w http.ResponseWriter, r *http.Request) {
	v := view.New(r)
	v.Name = "about/about"
	v.Render(w)
}
