package main

import (
	"log"
	"net/http"
)

// Функция-обработчик "home", которая записывает байтовый слайс, содержащий.
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello website"))
}

func main() {
	//Функция http.NewServeMux() для инициализации нового рутера.
	mux := http.NewServeMux()
	//Функцию "home" регистрируется как обработчик для URL-шаблона "/".
	mux.HandleFunc("/", home)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8484")
	//Используется функция http.ListenAndServe() для запуска нового веб-сервера.
	err := http.ListenAndServe(":8484", mux)
	//Функцию log.Fatal() для логирования ошибок.
	log.Fatal(err)
}
