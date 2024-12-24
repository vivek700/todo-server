package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"github.com/vivek700/todo-server/internal/database"
)

type Server struct {
	port int
	db   *database.Queries
}

var (
	url = fmt.Sprintf("%s%s%s", os.Getenv("TURSO_DATABASE_URL"), "?authToken=", os.Getenv("TURSO_AUTH_TOKEN"))
)

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	fmt.Println("Connecting to libsql database...")
	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	fmt.Println("Connected")

	queries := database.New(db)

	NewServer := &Server{
		port: port,
		db:   queries,
	}

	//declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}
