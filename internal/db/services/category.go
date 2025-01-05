package services

import (
	"github.com/onlyati/rss-collector/internal/db"
	"gorm.io/gorm"
)

//
// ===> CREATE
//

func CreateFavoriteService(gormDB *gorm.DB, category, username string) (*db.FavoriteCategory, error) {
	var user db.User
	err := gormDB.Where(&db.User{UserName: username}).First(&user).Error
	if err != nil {
		return nil, err
	}

	var deletedCategory db.FavoriteCategory
	err = gormDB.Unscoped().Where("user_id = ? and id = ?", user.ID, category).Find(&deletedCategory).Error
	if err != nil {
		return nil, err
	}

	if deletedCategory.Name == "" {
		newCategory := db.FavoriteCategory{
			Name:   category,
			UserID: user.ID,
		}
		err = gormDB.FirstOrCreate(&newCategory).Error
		if err != nil {
			return nil, err
		}

		return &newCategory, err
	} else {
		err := gormDB.Unscoped().Model(&db.FavoriteCategory{}).Where("id = ?", deletedCategory.ID).Update("deleted_at", nil).Error
		if err != nil {
			return nil, err
		}
		return &deletedCategory, nil
	}
}

//
// ===> Read
//

func ReadCategoriesService(gormDB *gorm.DB, index uint, username string) ([]db.FavoriteCategory, error) {
	var categories []db.FavoriteCategory
	err := gormDB.Where("id > ?", index).Limit(10).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

//
// ===> Delete
//

func DeleteCategoryService(gormDB *gorm.DB, category, username string) (*db.FavoriteCategory, error) {
	var user db.User
	err := gormDB.Where(&db.User{UserName: username}).First(&user).Error
	if err != nil {
		return nil, err
	}

	var cat db.FavoriteCategory
	err = gormDB.Where(&db.FavoriteCategory{UserID: user.ID, Name: category}).Delete(&cat).Error
	if err != nil {
		return nil, err
	}

	return &cat, nil
}
