package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/onlyati/rss-collector/internal/rss_model"
)

func (app *App) GetRSS(c *gin.Context) {
	var rssFeeds []rss_model.RSS
	err := app.Db.Find(&rssFeeds).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rssFeeds)
}

func (app *App) GetItem(c *gin.Context) {
	from_date := c.Query("from")
	if from_date == "" {
		from_date = "9999-12-31 00:00:00+00"
	}
	selectMode := c.Query("select")
	if selectMode == "" {
		selectMode = "select"
	}

	categoriesRaw := c.Query("categories")
	categories := []string{}
	if categoriesRaw != "" {
		categories = strings.Split(categoriesRaw, ",")
	}

	var rssItem []rss_model.RSSItem

	q := app.Db.Where("pub_date < ?", from_date)

	if len(categories) > 0 {
		if selectMode == "unselect" {
			q.Where("not (category && ?)", pq.Array(categories))
		} else {
			q.Where("category && ?", pq.Array(categories))
		}

	}

	q.Order("pub_date desc").Limit(10)
	err := q.Find(&rssItem).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rssItem)
}

func (app *App) GetCategories(c *gin.Context) {

}
