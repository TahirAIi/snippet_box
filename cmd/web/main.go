package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippet_box/pkg/models/db"
	"time"
)

type application struct {
	errorLogger   *log.Logger
	infoLogger    *log.Logger
	snippets      *db.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}
	address := os.Getenv("NETWORK_PORT")
	dns := os.Getenv("DB_DNS")

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbConnection, err := openDbConnection(dns)
	if err != nil {
		fmt.Println(err)
	}
	templateCache, err := newTemplateCache("./ui/html/")
	app := &application{
		snippets:      &db.SnippetModel{DB: dbConnection},
		templateCache: templateCache,
		infoLogger:    infoLogger,
		errorLogger:   errorLogger,
	}
	fmt.Println("Starting server")
	server := &http.Server{
		Addr:         ":4000",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Handler:      app.routes(),
	}
	fmt.Sprintf("Starting server at port %s", address)
	err = server.ListenAndServe()
	fmt.Println(err)
}

func openDbConnection(dns string) (*gorm.DB, error) {
	dbConnection, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}
	if err = dbConnection.Ping(); err != nil {
		return nil, err
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: dbConnection,
	}), &gorm.Config{})
	return gormDB, err
}
