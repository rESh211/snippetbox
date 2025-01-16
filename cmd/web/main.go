//Парсинг настроек конфигурации среды выполнения для приложения.
//Установление зависимостей для обработчиков.
//Запуск HTTP-сервера.

package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Создаем структуру `application` для хранения зависимостей всего веб-приложения.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

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

// Функция openDB() обертывает sql.Open() и возвращает пул соединений sql.DB для заданной строки подключения (DSN).
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	//Создаем новый флаг командной строки, значение по умолчанию: ":400.
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
	// Определение нового флага из командной строки для настройки MySQL подключения.
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "Название MySQL источника данных")
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
	errorLog := log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Функцию openDB().Мы передаем в нее полученный источник данных (DSN) из флага командной строки.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	// Мы также откладываем вызов db.Close(), чтобы пул соединений был закрыт до выхода из функции main().
	defer db.Close()

	// Инициализируем новую структуру с зависимостями приложения.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Создаем сервер.
	srv := &http.Server{
		// Адрес, на котором сервер будет слушать.
		Addr: *addr,
		// Лог для ошибок.
		ErrorLog: errorLog,
		// Мультиплексор для маршрутизации запросов.
		Handler: app.routes(),
	}

	// Запускаем веб-сервер и записываем лог.
	infoLog.Printf("Запуск сервера на %s", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
