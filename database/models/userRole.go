package models

import (
	"github.com/jinzhu/gorm"
	"github.com/YagoCarballo/kumquat-academy-api/database"
)

type (
	UserRoleModel struct {}
)
var DBUserRole UserRoleModel

func (model UserRoleModel) DB () *gorm.DB {
	return database.DB
}

func (model UserRoleModel) FindUserRole(roleName string) (*Role, error) {
	var userRole Role

	query := model.DB().First(&userRole, "name = ?", roleName)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.RecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	return &userRole, nil
}