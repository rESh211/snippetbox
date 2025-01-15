package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// Регистрируем два новых обработчика и соответствующие URL-шаблоны.
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Настройка файлового сервера.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
