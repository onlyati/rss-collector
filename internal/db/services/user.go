package services

import (
	"github.com/onlyati/rss-collector/internal/db"
	"gorm.io/gorm"
)

func CreateUser(gormDb *gorm.DB, name string) (*db.User, error) {
	newUser := db.User{
		UserName: name,
	}
	err := gormDb.Where(&db.User{UserName: name}).FirstOrCreate(&newUser)
	if err.Error != nil {
		return nil, err.Error
	}
	return &newUser, nil
}
