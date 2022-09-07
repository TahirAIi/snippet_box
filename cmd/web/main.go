package main

import (
	"database/sql"
	"fmt"
	"github.com/golangcollege/sessions"
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
	users         *db.UserModel
	templateCache map[string]*template.Template
	session       *sessions.Session
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	address, doAddressExists := os.LookupEnv("NETWORK_PORT")
	if !(doAddressExists == true) {
		panic("Unable to find value of NETWORK_PORT in env")
	}

	dns, dbDnsExists := os.LookupEnv("DB_DNS")
	if !dbDnsExists {
		panic("Unable to find  db DNS")
	}

	appId, appSecretExists := os.LookupEnv("APP_ID")

	if !appSecretExists {
		panic("Unable to find app secret")
	}

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbConnection, err := openDbConnection(dns)
	if err != nil {
		fmt.Println(err)
	}
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(templateCache)
	session := sessions.New([]byte(appId))
	session.Lifetime = 12 * time.Hour

	app := &application{
		snippets:      &db.SnippetModel{DB: dbConnection},
		users:         &db.UserModel{DB: dbConnection},
		templateCache: templateCache,
		infoLogger:    infoLogger,
		errorLogger:   errorLogger,
		session:       session,
	}
	fmt.Println("Starting server")
	server := &http.Server{
		Addr:         address,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Handler:      app.routes(),
	}
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
