package services

import (
	"github.com/onlyati/rss-collector/internal/db"
	"gorm.io/gorm"
)

//
// ===> Create
//

func CreateUserService(gormDb *gorm.DB, name string) (*db.User, error) {
	var deletedUser db.User
	err := gormDb.Unscoped().Where("user_name = ?", name).Find(&deletedUser).Error
	if err != nil {
		return nil, err
	}

	if deletedUser.UserName == "" {
		newUser := db.User{
			UserName: name,
		}
		err := gormDb.Where("user_name = ?", name).FirstOrCreate(&newUser).Error

		if err != nil {
			return nil, err
		}
		return &newUser, nil
	} else {
		err := gormDb.Unscoped().Model(&db.User{}).Where("id = ?", deletedUser.ID).Update("deleted_at", nil).Error
		if err != nil {
			return nil, err
		}
		return &deletedUser, nil
	}
}

//
// ===> Read
//

func ListUsersService(gormDB *gorm.DB, index string) ([]db.User, error) {
	var users []db.User
	err := gormDB.Where("id > ?", index).Limit(10).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func ReadUserService(gormDB *gorm.DB, username string) ([]db.User, error) {
	var user db.User
	err := gormDB.Where("user_name = ?", username).Find(&user).Error
	if err != nil {
		return nil, err
	}
	if user.UserName == "" {
		return []db.User{}, nil
	}

	return []db.User{user}, nil
}

//
// ===> Delete
//

func DeleteUsersHard(gormDB *gorm.DB) error {
	var users []db.User
	err := gormDB.Unscoped().Where("deleted_at is not null").Find(&users).Error
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	err = gormDB.Unscoped().Delete(&users).Error
	return err
}

func DeleteUserService(gormDB *gorm.DB, username string) error {
	err := gormDB.Where("user_name = ?", username).Delete(&db.User{}).Error
	return err
}
