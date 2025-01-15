package main

import (
	"flag"
	"log"
	"net/http"
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
	//Создаем новый флаг командной строки, значение по умолчанию: ":400
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
	flag.Parse()

	// Регистрируем два новых обработчика и соответствующие URL-шаблоны в
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Запускаем веб-сервер и записываем лог.
	log.Printf("Запуск сервера на %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
