package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://*", "http://*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Accept", "Content-Type", "X-CSRF-Token"},
		MaxAge:       300,
	}))

	e.GET("/", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   30 * 24 * 60 * 60,
			HttpOnly: true,
		}
		sess.Values["access_code"] = "akjkajfdsjfaljd"
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return err
		}
		return c.String(http.StatusOK, "hello from todo-server")
	})
	e.GET("/ses", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		return c.String(http.StatusOK, fmt.Sprintf("access_code=%v\n", sess.Values["access_code"]))
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
