package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Помощник записывает сообщение об ошибке отправляет пользователю ответ.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Помощник отправляет определенный код состояния и соответствующее описание пользователю.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Помощник отправляет пользователю ответ "404 Страница не найдена".
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
