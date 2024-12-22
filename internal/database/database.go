package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Task struct {
	Description string
	Status      bool
	CreatedAt   int
}

type Service interface {
	QueryTasks() map[string]string
	CreateNewTask(des string) string
	Close() error
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

func (s *service) CreateNewTask(des string) string {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		status BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := s.db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	insertTaskQuery := `
	INSERT INTO tasks (description, status) VALUES (?,?)
	`
	_, err = s.db.Exec(insertTaskQuery, des, false)
	if err != nil {
		log.Fatalf("Failed to insert tasks: %v", err)
	}

	fmt.Println("Task inserted successfully")

	return "created"

}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", url)
	return s.db.Close()
}
