package main

// Thanks [echo](https://github.com/labstack/echo) web framework for this useful example of CRUD API
import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

type User struct {
	ID   int
	Name string
}

var (
	users = map[int]*User{
		1: &User{1, "First User"},
		2: &User{2, "Second User"},
	}
	seq = 3
)

//----------
// Handlers
//----------

func createUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		u := &User{}
		if err := c.Bind(u); err != nil {
			return err
		}
		u.ID = seq
		users[u.ID] = u
		seq++
		return c.JSON(http.StatusCreated, u)
	}
}

func getUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		user, ok := users[id]
		if !ok {
			return c.NoContent(http.StatusNotFound)
		}
		return c.JSON(http.StatusOK, user)
	}
}

func updateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		id, _ := strconv.Atoi(c.Param("id"))
		users[id].Name = u.Name
		return c.JSON(http.StatusOK, users[id])
	}
}
func deleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		delete(users, id)
		return c.NoContent(http.StatusNoContent)
	}
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Routes
	e.Post("/users", createUser())
	e.Get("/users/:id", getUser())
	e.Patch("/users/:id", updateUser())
	e.Delete("/users/:id", deleteUser())

	// Start server
	e.Run(standard.New(":1323"))
}
