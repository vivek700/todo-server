package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Service interface {
	QueryTasks() map[string]string
}

type service struct {
	db *sql.DB
}

var (
	url        = fmt.Sprintf("%s%s%s", os.Getenv("TURSO_DATABASE_URL"), "?authToken=", os.Getenv("TURSO_AUTH_TOKEN"))
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	// fmt.Println(url)
	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	defer db.Close()
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) QueryTasks() map[string]string {
	tasks := make(map[string]string)
	tasks["task1"] = "something"
	return tasks
}
