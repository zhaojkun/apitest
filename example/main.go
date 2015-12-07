package main

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/octokit/go-octokit/octokit"
)

func main() {
	e := echo.New()

	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Routes
	e.Get("/hello", hello)
	e.Get("/user/:name", getUser)

	// Start server
	e.Run(":1323")
}

func hello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!\n")
}

func getUser(c *echo.Context) error {
	username := c.Param("name")
	if username == "" {
		return errors.New("parameter 'name' must be provided")
	}

	user, found, err := fetchUserFromGithub(username)
	if err != nil {
		return err
	} else if !found {
		return c.String(http.StatusNotFound, "user %s not found", username)
	}

	return c.JSON(http.StatusOK, user)
}

func fetchUserFromGithub(username string) (user *octokit.User, found bool, err error) {
	if username == "BadGuy" {
		return nil, false, errors.New("BadGuy failed me :(")
	}
	client := octokit.NewClient(nil)
	userURL, _ := octokit.UserURL.Expand(octokit.M{"user": username})

	var result *octokit.Result
	user, result = client.Users(userURL).One()

	found = true
	if result.Err != nil {
		err = result.Err
		if responseErr, ok := result.Err.(*octokit.ResponseError); ok {
			found = responseErr.Type != octokit.ErrorNotFound
			if !found {
				err = nil
			}
		}

	}

	return user, found, err
}
