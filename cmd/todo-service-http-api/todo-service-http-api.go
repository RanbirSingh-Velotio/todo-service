package main

import (
	"flag"
	"fmt"
	"github.com/RanbirSingh-Velotio/todo-service/pkg/config"
	"github.com/RanbirSingh-Velotio/todo-service/pkg/todo"
	todoHandler "github.com/RanbirSingh-Velotio/todo-service/pkg/todo/handler"
	todoService "github.com/RanbirSingh-Velotio/todo-service/pkg/todo/service"
	sqliteService "github.com/RanbirSingh-Velotio/todo-service/store/sqlite"
	"github.com/RanbirSingh-Velotio/todo-service/utils/handlerutil"
	"github.com/google/gops/agent"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// getConfigDir returns config path(string) based on environment
func getConfigDir(environ string) string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err)
		return ""
	}

	mainIniFile := "todo-main.ini"
	baseDir := strings.TrimSuffix(dir, "cmd\\todo-service-http-api")
	filesDir := "files\\etc\\"
	configDir := baseDir + filesDir + fmt.Sprintf("config\\%s\\", environ)
	if _, err := os.Stat(configDir + mainIniFile); os.IsNotExist(err) {
		// default config files located based on bin files path,
		// if not exists then the environment is not development or customized
		filesDir = "\\etc\\"
		configDir = filesDir + fmt.Sprintf("config\\%s\\", environ)
	}
	return configDir
}

func loadMainConfigFile(configDir string) {

	err := config.NewMainConfig(configDir + "todo-main.ini")
	if err != nil {
		log.Print("error")
	}
}

func initializeConfig() {
	configTest := flag.Bool("t", false, "config test")
	flag.Parse()

	environ := os.Getenv("APP_ENV")
	log.Println("Environment is : " + environ)
	if environ == "" {
		environ = "development"
	}

	configDir := getConfigDir(environ)
	loadMainConfigFile(configDir)

	//Exit if test-flag is given
	if *configTest {
		log.Println("Test flag is given for config test")
		os.Exit(0)
	}
}

// startServer start server using grace
func startServer() {
	// gops for profiling
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Printf("Error")
	}

	defer agent.Close()

	// Start server
	mConf, _ := config.GetConfig()
	serverPort := ":" + strconv.Itoa(mConf.Server.Port)

	err := http.ListenAndServe(serverPort, nil)

	if err != nil {
		log.Printf("unable to shutdown http server gracefully: %v\n", err)
	}
}

func handlerHealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Success")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func initializeTodoService() {
	db := initDatabase()
	sqlitSrv := sqliteService.New(db)
	todoSrv := todoService.New(sqlitSrv)
	todo.Init(todoSrv)
	handler := todoHandler.InitHandler(todoSrv)
	handlerutil.Add(handler)
}

// Database Connect function
func initDatabase() *sqlx.DB {
	var err error
	db, err := sqlx.Connect("sqlite3", "todo.db")
	if err != nil {
		panic("failed to connect database")
	}
	sqlStmt := `
	create table IF NOT EXISTS todo (id integer not null primary key, name text,completed bool);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)

	}
	db.SetMaxOpenConns(10)

	fmt.Println("Database successfully connected")
	return db
}

func main() {
	initializeConfig()

	initializeTodoService()

	handlerutil.Start()

	http.HandleFunc("/", notFoundHandler)

	http.HandleFunc("/health", handlerHealthCheck)

	startServer()
}
