package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onlyati/rss-collector/internal/db/services"
)

//
// ===> CREATE endpoints
//

func (app *App) RegisterUser(c *gin.Context) {
	username := c.GetString("username")

	user, err := services.CreateUserService(app.Db, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (app *App) CreateUser(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing username parameter"})
		return
	}

	user, err := services.CreateUserService(app.Db, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

//
// ===> READ endpoints
//

func (app *App) ListUsers(c *gin.Context) {
	cursor := c.Query("index")
	if cursor == "" {
		cursor = "0"
	}
	username := c.Query("username")

	if username != "" {
		user, err := services.ReadUserService(app.Db, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	} else {
		users, err := services.ListUsersService(app.Db, cursor)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

//
// ===> DELETE endpoints
//

func (app *App) DeleteUser(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing username parameter"})
		return
	}

	user, err := services.DeleteUserService(app.Db, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user.UserName == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "user does not exists"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (app *App) DeleteUserHard(c *gin.Context) {
	err := services.DeleteUsersHard(app.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
