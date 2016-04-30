package models

import (
	"github.com/jinzhu/gorm"
	"github.com/YagoCarballo/kumquat-academy-api/database"
	. "github.com/YagoCarballo/kumquat-academy-api/constants"
)

type (
	PermissionsModel struct {}

	PermissionsTable struct {
		ModuleCode		string `json:"-"`
		ModuleId		uint32 `json:"-"`
		RoleId			uint32 `json:"id"`
		RoleName		string `json:"name"`
		RoleDescription	string `json:"description"`
		Admin			bool `json:"admin"`
		Read			bool `json:"read"`
		Write			bool `json:"write"`
		Delete			bool `json:"delete"`
		Update			bool `json:"update"`
	}
)

var DBPermissions PermissionsModel

func (model PermissionsModel) DB () *gorm.DB {
	return database.DB
}

func (model PermissionsModel) IsActionPermittedOnModule(userId uint32, id interface{}, action AccessRight) bool {
	moduleId := uint32(id.(uint32))
	permissions, err := model.GetPermissionsForModule(&userId, &moduleId, nil)
	if err != nil {
		return false
	}

	if permissions.Admin {
		return true
	}

	switch action {
	case ReadPermission:
		return permissions.Read
	case WritePermission:
		return permissions.Write
	case DeletePermission:
		return permissions.Delete
	case UpdatePermission:
		return permissions.Update
	}

	return false
}

func (model PermissionsModel) IsActionPermittedOnModuleWithCode(userId uint32, code interface{}, action AccessRight) bool {
	var moduleId uint32
	moduleCode := code.(string)
	permissions, err := model.GetPermissionsForModule(&userId, &moduleId, &moduleCode)
	if err != nil {
		return false
	}

	if permissions.Admin {
		return true
	}

	switch action {
	case ReadPermission:
		return permissions.Read
	case WritePermission:
		return permissions.Write
	case DeletePermission:
		return permissions.Delete
	case UpdatePermission:
		return permissions.Update
	}

	return false
}

func (model PermissionsModel) IsActionPermittedOnCourse(userId uint32, courseId interface{}, action AccessRight) bool {
	permissions, err := model.GetPermissionsForCourse(userId, courseId.(uint32))
	if err != nil {
		return false
	}

	if permissions.Admin {
		return true
	}

	switch action {
	case ReadPermission:
		return permissions.Read
	case WritePermission:
		return permissions.Write
	case DeletePermission:
		return permissions.Delete
	case UpdatePermission:
		return permissions.Update
	}

	return false
}

func (model PermissionsModel) IsAdmin(userId uint32) bool {
	var isAdmin []bool

	// Run the query
	query := model.DB().Table("users").Where("id = ?", userId).Select("admin").Pluck("admin", &isAdmin)
	if query.Error != nil || len(isAdmin) <= 0 {
		return false
	}

	return isAdmin[0]
}

func (model PermissionsModel) IsUsernameAnAdmin(username string) bool {
	var isAdmin []bool

	// Run the query
	query := model.DB().Table("users").Where("username = ?", username).Select("admin").Pluck("admin", &isAdmin)
	if query.Error != nil || len(isAdmin) <= 0 {
		return false
	}

	return isAdmin[0]
}

func (model PermissionsModel) CanUserSubmitAssignment(username string, assignmentId uint32) bool {
	var canSubmit []uint32

	query := model.DB().Table("users").Select("assignments.id").Joins(
		"inner join user_modules on user_modules.`user_id` = users.id " +
		"left outer join assignments on assignments.`module_code` = user_modules.`module_code` ",
	).Where(
		"users.username = ? and assignments.status = ? and assignments.id = ?",
		username,
		AssignmentAvailable,
		assignmentId,
	).Having(
		"(select count(*) from submissions " +
		"where submissions.`status` != 'canceled' " +
		"and submissions.`assignment_id` = assignments.`id` " +
		"and submissions.`user_id` = users.`id`) = 0",
	).Pluck("assignments.id", &canSubmit)

	// If there are any errors or no results, then the studen't cannot submit the assignment
	if query.Error != nil || len(canSubmit) <= 0 {
		return false

	// If there are results, then the student can submit the assignment
	} else {
		return true
	}
}

func (model PermissionsModel) GetPermissionsForModule (userId, moduleId *uint32, code *string) (*PermissionsTable, error) {
	var permissions []PermissionsTable
	var moduleIdentifier interface{}

	// Use the module code if provided or the moduleId instead.
	moduleIdentifier = *moduleId
	filter := "user_modules.user_id = ? and level_modules.module_id = ?";
	if code != nil {
		filter = "user_modules.user_id = ? and level_modules.code = ?";
		moduleIdentifier = *code
	}

	rows, err := model.DB().Table("user_modules").
		Select(`
			level_modules.code,
			level_modules.module_id,
			users.admin,
			moduleRole.id 'role_id',
			moduleRole.name 'role_name',
			moduleRole.description 'role_description',
			(moduleRole.can_read + coalesce(courseRole.can_read, 0) != 0) 'read',
			(moduleRole.can_write + coalesce(courseRole.can_write, 0) != 0) 'write',
			(moduleRole.can_delete + coalesce(courseRole.can_delete, 0) != 0) 'delete',
			(moduleRole.can_update + coalesce(courseRole.can_update, 0) != 0) 'update'
		`).Joins(`
			inner join users on users.id = user_modules.user_id
			inner join level_modules on level_modules.code = user_modules.module_code
			inner join modules as module on module.id = level_modules.module_id
			inner join roles as moduleRole on moduleRole.id = user_modules.role_id
			inner join classes on classes.id = level_modules.class_id
			inner join courses as course on course.id = classes.course_id
			left outer join user_courses on user_courses.course_id = classes.course_id and user_courses.user_id = user_modules.user_id
			left outer join roles as courseRole on courseRole.id = user_courses.role_id
		`).Where(filter, *userId, moduleIdentifier).Rows()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		permission := PermissionsTable{}

		err = rows.Scan(
			&permission.ModuleCode,
			&permission.ModuleId,
			&permission.Admin,
			&permission.RoleId,
			&permission.RoleName,
			&permission.RoleDescription,
			&permission.Read,
			&permission.Write,
			&permission.Delete,
			&permission.Update,
		)

		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if len(permissions) <= 0 {
		isAdmin := model.IsAdmin(*userId)
		return &PermissionsTable{
			ModuleId: *moduleId,
			Admin: isAdmin,
			RoleName: "admin",
		}, nil
	}

	return &permissions[0], nil
}

func (model PermissionsModel) GetPermissionsForCourse (userId, courseId uint32) (*PermissionsTable, error) {
	var permissions []PermissionsTable

	rows, err := model.DB().Table("user_courses").
	Select(`
			users.admin,
			courseRole.id 'role_id',
			courseRole.name 'role_name',
			courseRole.description 'role_description',
			(courseRole.can_read != 0) 'read',
			(courseRole.can_write != 0) 'write',
			(courseRole.can_delete != 0) 'delete',
			(courseRole.can_update != 0) 'update'
		`).Joins(`
			inner join users on users.id = user_courses.user_id
			inner join roles as courseRole on courseRole.id = user_courses.role_id
			inner join courses as course on course.id = user_courses.course_id
		`).Where("user_courses.user_id = ? and user_courses.course_id = ?", userId, courseId).Rows()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		permission := PermissionsTable{}

		err = rows.Scan(
			&permission.Admin,
			&permission.RoleId,
			&permission.RoleName,
			&permission.RoleDescription,
			&permission.Read,
			&permission.Write,
			&permission.Delete,
			&permission.Update,
		)

		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if len(permissions) <= 0 {
		isAdmin := model.IsAdmin(userId)
		return &PermissionsTable{
			Admin: isAdmin,
			RoleName: "admin",
		}, nil
	}

	return &permissions[0], nil
}
