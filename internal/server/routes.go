package server

import (
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

	e.RouteNotFound("/*", func(c echo.Context) error {
		return c.String(http.StatusNotFound, "Not found")
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello! I am the server.")

	})

	e.GET("/tasks", s.listTasksHandler)

	e.POST("/tasks", s.createTaskHandler)

	return e
}

func (s *Server) listTasksHandler(c echo.Context) error {

	userID, err := c.Cookie("access_code")
	if err != nil || userID.Value == "" {
		newUUID := uuid.New().String()

		cookie := new(http.Cookie)
		cookie.Name = "access_code"
		cookie.Value = newUUID
		cookie.HttpOnly = true // Secure: not accessible via JavaScript
		cookie.Secure = false  //set to true if using https
		cookie.Expires = time.Now().Add(30 * 24 * time.Hour)

		c.SetCookie(cookie)

		res, err := s.db.CreateUser(c.Request().Context(), newUUID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, res)
	}

	resUser, _ := s.db.GetUser(c.Request().Context(), userID.Value)

	tasks, err := s.db.ListTasks(c.Request().Context(), resUser)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch tasks.")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Tasks retrieved successfully",
		"status":  "success",
		"data":    tasks,
	})

}

func (s *Server) createTaskHandler(c echo.Context) error {

	userID, err := c.Cookie("access_code")
	if err != nil || userID.Value == "" {
		return c.NoContent(http.StatusForbidden)
	}

	task := new(Task)
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid payload",
		})
	}

	if task.Description == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "description is required"})
	}

	user, _ := s.db.GetUser(c.Request().Context(), userID.Value)

	res, err := s.db.CreateTask(
		c.Request().Context(),
		database.CreateTaskParams{UserID: user, Description: task.Description, Status: false},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to create task: " + err.Error(),
			"status":  "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Task created successfully",
		"task":    res,
	})
}

// func (s *Server) deleteTaskHandler(c echo.Context) error {

// }
