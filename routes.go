package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *App) routes() http.Handler {
	router := httprouter.New()
	fileServer := http.FileServer(http.Dir("./templates/static"))

	dynamic := alice.New(app.sessionManager.LoadAndSave)
	dynamic = dynamic.Append(app.loggingMiddleware)

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.indexHandler))

	router.Handler(http.MethodGet, "/error", dynamic.ThenFunc(app.errorHandler))

	return router
}
