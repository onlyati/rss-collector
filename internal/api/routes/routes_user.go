package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onlyati/rss-collector/internal/db/services"
)

//
// ===> CREATE endpoints
//

// Create new user where the user ID is get from the access token
func (app *App) CreateUser(c *gin.Context) {
	username := c.GetString("username")

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

// List existing users
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

// Soft delete user
func (app *App) DeleteUser(c *gin.Context) {
	username := c.GetString("username")

	err := services.DeleteUserService(app.Db, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// Permanently delete soft deleted users
func (app *App) DeleteUserHard(c *gin.Context) {
	err := services.DeleteUsersHard(app.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
