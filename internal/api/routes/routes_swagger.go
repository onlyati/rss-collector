package routes

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

func (app *App) GetSwaggerYAML(c *gin.Context) {
	template, err := template.ParseFiles("openapi/openapi.yaml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	builder := &strings.Builder{}

	if app.AuthOptions != nil {
		err = template.Execute(builder, gin.H{
			"HostName":   app.Hostname,
			"Port":       app.Port,
			"AuthEnable": true,
			"JWKSauth":   app.AuthOptions.AuthorizationEndpoint,
			"JWKStoken":  app.AuthOptions.TokenEndpoint,
		})
	} else {
		err = template.Execute(builder, gin.H{
			"HostName":   app.Hostname,
			"Port":       app.Port,
			"AuthEnable": false,
		})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := builder.String()

	c.Data(http.StatusOK, "application/yaml", []byte(result))
}

func (app *App) GetSwaggerUI(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"HostName": app.Hostname,
		"Port":     app.Port,
	})
}

func (app *App) GetRedirect(c *gin.Context) {
	c.HTML(http.StatusOK, "oauth2-redirect.html", gin.H{})
}
