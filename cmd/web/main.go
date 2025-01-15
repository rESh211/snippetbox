package main

import (
	"log"
	"net/http"
)

func main() {
	// Регистрируем два новых обработчика и соответствующие URL-шаблоны в
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8484")
	err := http.ListenAndServe(":8484", mux)
	log.Fatal(err)
}
