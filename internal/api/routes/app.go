package routes

import "gorm.io/gorm"

type App struct {
	Db *gorm.DB
}
