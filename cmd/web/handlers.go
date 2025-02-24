package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Обработчик главной странице.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Проверяется, если текущий путь URL запроса точно совпадает с шаблоном "/".
	//Если нет, вызывается функция http.NotFound() для возвращения клиенту ошибки 404.
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Инициализируем срез содержащий пути к двум файлам.
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Используем функцию template.ParseFiles() для чтения файла шаблона.
	// Если возникла ошибка,ответ: 500 Internal Server Error (Внутренняя ошибка на сервере).
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Затем мы используем метод Execute() для записи содержимого шаблона в тело HTTP ответа.
	//Последний параметр в Execute() предоставляет возможность отправки динамических данных в шаблон.
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// Обработчик для отображения содержимого заметки.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Извлекаем значение параметра id из URL
	//Если его нельзя конвертировать в integer, или значение меньше 1, возвращаем ответ 404 - страница не найдена!
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Используем функцию fmt.Fprintf() для вставки значения из id в строку ответа и записываем его в http.ResponseWriter.
	fmt.Fprintf(w, "Отображение выбранной заметки с ID %d...", id)
}

// Обработчик для создания новой заметки.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Используем r.Method для проверки, использует ли запрос метод POST или нет.
	if r.Method != http.MethodPost {
		// Используем метод Header().Set() для добавления заголовка 'Allow: POST' в карту HTTP-заголовков.
		w.Header().Set("Allow", http.MethodPost)
		// Используем функцию http.Error() для отправки кода состояния 405 с соответствующим сообщением.
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Создание новой заметки..."))
}
