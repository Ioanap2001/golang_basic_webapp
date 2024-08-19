package main

import "net/http"

func (app *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r, w)
	app.render(w, r, http.StatusOK, "index.html", data)
}
