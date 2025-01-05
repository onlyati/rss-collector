package services

import (
	"errors"

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
	err = gormDB.Unscoped().Where("user_id = ? and name = ?", user.ID, category).Find(&deletedCategory).Error
	if err != nil {
		return nil, err
	}

	if deletedCategory.Name == "" {
		newCategory := db.FavoriteCategory{
			Name:   category,
			UserID: user.ID,
		}
		err = gormDB.Create(&newCategory).Error
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

func ListCategoriesService(gormDB *gorm.DB, index string, username string) ([]db.FavoriteCategory, error) {
	var categories []db.FavoriteCategory
	err := gormDB.Where("id > ?", index).Limit(10).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func ReadCategoryService(gormDB *gorm.DB, category, username string) ([]db.FavoriteCategory, error) {
	var user db.User
	err := gormDB.Where(&db.User{UserName: username}).First(&user).Error
	if err != nil {
		return nil, err
	}

	var cat db.FavoriteCategory
	err = gormDB.Where("user_id = ? and name = ?", user.ID, category).Find(&cat).Error
	if err != nil {
		return nil, err
	}
	if cat.Name == "" {
		return []db.FavoriteCategory{}, nil
	}
	return []db.FavoriteCategory{cat}, nil
}

//
// ===> Delete
//

func DeleteCategoriesHard(gormDB *gorm.DB) error {
	var categories []db.FavoriteCategory
	err := gormDB.Unscoped().Where("deleted_at is not null").Find(&categories).Error
	if err != nil {
		return err
	}

	if len(categories) == 0 {
		return nil
	}

	err = gormDB.Unscoped().Delete(&categories).Error
	return err
}

func DeleteCategoryService(gormDB *gorm.DB, category, username string) error {
	var user db.User
	err := gormDB.Where("user_name = ?", username).Find(&user).Error
	if err != nil {
		return err
	}

	if user.UserName == "" {
		return errors.New("user does not exists")
	}

	err = gormDB.Where("user_id = ? and name = ?", user.ID, category).Delete(&db.FavoriteCategory{}).Error
	return err
}
