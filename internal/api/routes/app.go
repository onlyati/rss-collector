package routes

import (
	"github.com/onlyati/rss-collector/internal/api/auth"
	"gorm.io/gorm"
)

type App struct {
	Db          *gorm.DB
	Hostname    string
	Port        int
	AuthOptions *auth.KeycloakLinks
}
