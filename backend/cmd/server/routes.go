package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/Books", app.listBooksHandler)                                               ///
	router.HandlerFunc(http.MethodPost, "/Books", app.requirePermission("movies:write", app.createBookHandler))      ////
	router.HandlerFunc(http.MethodGet, "/Books/:id", app.showBookHandler)                                            ////
	router.HandlerFunc(http.MethodPatch, "/Books/:id", app.requirePermission("movies:write", app.updateBookHandler)) ///
	router.HandlerFunc(http.MethodDelete, "/Books/:id", app.deleteBookHandler)                                       ////

	router.HandlerFunc(http.MethodGet, "/WritingVariant", app.listWritingVariantHandler)                                                 ///
	router.HandlerFunc(http.MethodPost, "/WritingVariant", app.requirePermission("movies:write", app.createWritingVariantHandler))       ////
	router.HandlerFunc(http.MethodGet, "/WritingVariant/:id", app.showWritingVariantHandler)                                             ////
	router.HandlerFunc(http.MethodPatch, "/WritingVariant/:id", app.requirePermission("movies:write", app.updateWritingVariantHandler))  ///
	router.HandlerFunc(http.MethodDelete, "/WritingVariant/:id", app.requirePermission("movies:write", app.deleteWritingVariantHandler)) ////

	router.HandlerFunc(http.MethodPost, "/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)

	return app.enableCORS(app.recoverPanic(app.rateLimit(app.authenticate(router))))

}
