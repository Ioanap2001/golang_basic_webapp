package main

import (
	"bytes"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
)

type TemplateData struct {
	Form      any
	Flash     string
	CSRFToken string
}

func (app *App) newTemplateData(r *http.Request, w http.ResponseWriter) *TemplateData {
	return &TemplateData{
		Flash:     app.sessionManager.PopString(r.Context(), "flash"),
		CSRFToken: nosurf.Token(r),
	}
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/*.tmpl". This will essentially gives
	// us a slice of all the filepaths for our application 'page' templates
	// like: [ui/index.tmpl ui/view.tmpl]
	generalPages, err := filepath.Glob("./templates/*.html")
	if err != nil {
		return nil, err
	}

	allPages := [][]string{
		generalPages,
	}
	var pages []string
	for _, r := range allPages {
		pages = append(pages, r...)
	}
	// Loop through the page filepaths one-by-one.
	for _, page := range pages {

		name := filepath.Base(page)
		// Parse the base template file into a template set.
		ts, err := template.New(name).ParseFiles("./templates/base.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./templates/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the  page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}

func (app *App) render(w http.ResponseWriter, r *http.Request, status int, page string, data *TemplateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		app.errorLog.Printf("The template %s does not exist", page)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	app.checkErr(err, w, r)

	w.WriteHeader(status)
	_, err = buf.WriteTo(w)
	app.checkErr(err, w, r)
}
