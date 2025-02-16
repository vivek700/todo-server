package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vivek700/todo-server/internal/database"

	_ "github.com/joho/godotenv/autoload"
)

var frontendUrl string = fmt.Sprintf("%s", os.Getenv("FRONTEND_URL"))

type Task struct {
	Description string `json:"description"`
}

type TaskID struct {
	ID int `query:"id"`
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{frontendUrl},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Accept", "Content-Type", "X-CSRF-Token"},
		MaxAge:           300,
	}))

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		c.String(http.StatusNotFound, "Error: Not found")
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello! I am the server.")

	})

	e.GET("/tasks", s.listTasksHandler)

	e.POST("/tasks", s.createTaskHandler)

	e.DELETE("/tasks", s.deleteTaskHandler)

	e.PUT("/tasks", s.updateTaskHandler)
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
		cookie.Secure = true   //set to true if using https
		cookie.SameSite = http.SameSiteNoneMode
		cookie.Path = "/"
		cookie.Expires = time.Now().AddDate(1, 0, 0)
		c.SetCookie(cookie)

		res, err := s.db.CreateUser(c.Request().Context(), newUUID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		fmt.Println(res)
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

func (s *Server) deleteTaskHandler(c echo.Context) error {
	userID, err := c.Cookie("access_code")
	if err != nil || userID.Value == "" {
		return c.NoContent(http.StatusForbidden)
	}

	if c.QueryParam("id") == "" {
		return c.String(http.StatusBadRequest, "Error: Task ID is required")
	}

	req := new(TaskID)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	user, _ := s.db.GetUser(c.Request().Context(), userID.Value)

	err = s.db.DeleteTask(c.Request().Context(), database.DeleteTaskParams{UserID: int64(user), ID: int64(req.ID)})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to delete task: " + err.Error(),
			"status":  "error",
		})
	}

	return c.String(http.StatusOK, "Task deleted successfully")

}

type TaskUpdate struct {
	Id     int  `json:"id" validate:"required"`
	Status bool `json:"status" validate:"required"`
}

func (s *Server) updateTaskHandler(c echo.Context) error {
	userID, err := c.Cookie("access_code")
	if err != nil || userID.Value == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	req := new(TaskUpdate)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}
	if req.Id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid task ID",
		})
	}
	fmt.Println(req.Id, req.Status)
	user, _ := s.db.GetUser(c.Request().Context(), userID.Value)
	err = s.db.UpdateTask(c.Request().Context(), database.UpdateTaskParams{ID: int64(req.Id), UserID: user, Status: req.Status})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to update the task" + err.Error(),
			"status":  "error",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Task updated successfully",
	})

}
