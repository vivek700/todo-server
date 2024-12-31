package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
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
		userid, err := c.Cookie("access_code")
		if err != nil || userid.Value == "" {
			newUUID := uuid.New().String()

			cookie := new(http.Cookie)
			cookie.Name = "access_code"
			cookie.Value = newUUID
			cookie.HttpOnly = true // Secure: not accessible via JavaScript
			cookie.Secure = false  //set to true if using https
			cookie.Expires = time.Now().Add(30 * 24 * time.Hour)

			c.SetCookie(cookie)

			return c.String(http.StatusOK, "hello from todo-server")
		}
		return c.String(http.StatusOK, "Welcome back! Your user ID is: "+userid.Value)

	})

	e.GET("/tasks", s.listTasksHandler)

	e.POST("/tasks", s.createTaskHandler)

	return e
}

func (s *Server) listTasksHandler(c echo.Context) error {
	data, err := s.db.ListTasks(c.Request().Context(), 2)
	if err != nil {
		log.Fatal("error in listing item")
	}

	fmt.Println(data)

	return c.JSON(http.StatusOK, data)
}

func (s *Server) createTaskHandler(c echo.Context) error {
	task := new(Task)
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid payload",
		})
	}

	if task.Description == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "description is required"})
	}

	res, _ := s.db.CreateTask(
		c.Request().Context(),
		database.CreateTaskParams{Description: task.Description, Status: false},
	)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Task created successfully",
		"task":    res,
	})
}
