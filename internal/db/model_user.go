package db

import "gorm.io/gorm"

type User struct {
	*gorm.Model
	UserName  string             `json:"user_name" gorm:"not null;unique;index"`
	Favorites []FavoriteCategory `json:"favorites" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type FavoriteCategory struct {
	*gorm.Model
	UserID uint   `json:"user_id"`
	Name   string `json:"name" gorm:"no null"`
}
