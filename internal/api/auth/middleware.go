package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(auth *Authentication) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokens := strings.Split(authHeader, " ")
		if len(tokens) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication header is not enough long"})
			return
		}

		claims, err := auth.validate(tokens[1])
		if err != nil {
			if err := auth.refreshJWKS(); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}

			claims, err = auth.validate((tokens[1]))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
		}

		ok := false
		for key, claim := range *claims {
			if key == "preferred_username" {
				switch v := claim.(type) {
				case string:
					c.Set("username", v)
					ok = true
				}
			}
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "preferred_username is not in the token as string"})
			return
		}

		c.Next()
	}
}
