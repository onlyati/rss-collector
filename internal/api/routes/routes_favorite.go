package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onlyati/rss-collector/internal/db/services"
)

//
// ===> CREATE endpoints
//

// Create a new favorite category for specific user
func (app *App) CreateCategory(c *gin.Context) {
	username := c.GetString("username")
	category := c.Query("category")

	if username == "" || category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and category are mandatory query parameters"})
		return
	}

	newCategory, err := services.CreateFavoriteService(app.Db, category, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newCategory)
}

//
// ===> READ endpoints
//

// List categories
func (app *App) ListCategories(c *gin.Context) {
	cursor := c.Query("index")
	if cursor == "" {
		cursor = "0"
	}
	username := c.GetString("username")
	category := c.Query("category")

	if category != "" {
		categories, err := services.ReadCategoryService(app.Db, category, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, categories)
	} else {
		categories, err := services.ListCategoriesService(app.Db, cursor, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, categories)
	}
}

//
// ===> DELETE endpoints
//

// Soft delete category
func (app *App) DeleteCategory(c *gin.Context) {
	username := c.GetString("username")
	category := c.Query("category")
	if username == "" || category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing username and/or category parameter"})
		return
	}

	err := services.DeleteCategoryService(app.Db, category, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// Permanently delete soft deleted categories
func (app *App) DeleteCategoryHard(c *gin.Context) {
	err := services.DeleteCategoriesHard(app.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
