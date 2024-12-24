package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vivek700/todo-server/internal/database"
)

type Task struct {
	Description string `json:"description"`
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://*", "http://*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Accept", "Content-Type", "X-CSRF-Token"},
		MaxAge:       300,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello from todo-server")
	})

	e.GET("/tasks", s.listTaskHandler)

	e.POST("/addtask", s.taskhandler)

	return e
}

var ctx = context.Background()

func (s *Server) listTaskHandler(c echo.Context) error {

	data, err := s.db.ListTasks(ctx)

	if err != nil {
		log.Fatal("error in listing item")
	}

	fmt.Println(data)

	return c.JSON(http.StatusOK, data)
}

func (s *Server) taskhandler(c echo.Context) error {

	task := new(Task)
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid payload",
		})
	}

	if task.Description == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "description is required"})

	}
	res, _ := s.db.CreateTask(ctx, database.CreateTaskParams{Description: task.Description, Status: false})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Task created successfully",
		"task":    res,
	})

}
