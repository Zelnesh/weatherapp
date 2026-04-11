package main

import (

	"net/http"
)


func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/consultspace", app.consultSpace)

	return secureHeaders(blockStaticFolder(mux))
}
