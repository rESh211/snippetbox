package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

func main() {
	//Создаем новый флаг командной строки, значение по умолчанию: ":400.
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
	flag.Parse()

	// Открываем файл для записи логов.
	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Создаем логгер для записи информационных сообщений.
	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	// Создаем логгер для записи сообщений об ошибках.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Регистрируем два новых обработчика и соответствующие URL-шаблоны.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Настройка файлового сервера.
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Создаем сервер.
	srv := &http.Server{
		// Адрес, на котором сервер будет слушать.
		Addr: *addr,
		// Лог для ошибок.
		ErrorLog: errorLog,
		// Мультиплексор для маршрутизации запросов.
		Handler: mux,
	}

	// Запускаем веб-сервер и записываем лог.
	infoLog.Printf("Запуск сервера на %s", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
