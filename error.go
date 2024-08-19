package main

import "net/http"

func (app *App) checkErr(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		app.errorLog.Println(err)
		http.Redirect(w, r, "/error", http.StatusFound)
	}
}

func (app *App) errorHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r, w)
	app.render(w, r, http.StatusOK, "error.html", data)
}
