package models

import (
	"github.com/jinzhu/gorm"
	"github.com/YagoCarballo/kumquat.academy.api/database"
	"database/sql"
	"time"
	"fmt"
	"strings"
)

type ModulesModel struct {}
var DBModule ModulesModel

func (model ModulesModel) DB () *gorm.DB {
	return database.DB
}

type (
	OutputModule struct {
		Id 		uint32		`json:"id"`
		Level		uint32		`json:"level"`
		CourseId	uint32		`json:"course_id"`
		ClassId		uint32		`json:"class_id"`
		Code		string		`json:"code"`
		Title		string		`json:"title"`
		Description	string		`json:"description"`
		Color		string		`json:"color"`
		Icon		string		`json:"icon"`
		Duration	uint32		`json:"duration"`
		Role		*OutputRole	`json:"role,omitempty"`
		Year		string		`json:"year"`
		Course		*OutputCourse	`json:"course,omitempty"`
		Status		ModuleStatus	`json:"status"`
	}

	OutputRole struct {
		Id			uint32	`json:"id"`
		Name		string	`json:"name"`
		Description	string	`json:"description"`
		Read		bool	`json:"read"`
		Write		bool	`json:"write"`
		Delete		bool	`json:"delete"`
		Update		bool	`json:"update"`
		Admin		bool	`json:"admin"`
	}

	OutputCourse struct {
		Id			uint32	`json:"id"`
		Title		string	`json:"name"`
		Description	string	`json:"description"`
	}
)

func (model ModulesModel) CreateModule(module Module) (*Module, error) {
	// Query the Module
	query := model.DB().Create(&module)
	if query.Error != nil {
		return nil, query.Error
	}

	// Returns the Module
	return &module, nil
}

func (model ModulesModel) FindModule(id uint32) (*Module, error) {
	// Creates empty Module
	var module Module

	// Query the Module
	query := model.DB().First(&module, "id = ?", id)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.RecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	// Returns the Module
	return &module, nil
}

func (model ModulesModel) DeleteModule(id uint32) (int64, error) {
	query := model.DB().
		Table("modules").
		Where("id = ?", id).
		Delete(CourseLevel{})
	if query.Error != nil {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}

func (model ModulesModel) FindModuleWithCode(code string) (*LevelModule, error) {
	var levelModule LevelModule

	query := model.DB().
				Preload("Module").
				Preload("Class").
				Table("level_modules").
				Where("code = ?", code).
				First(&levelModule)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.RecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	return &levelModule, nil
}

func (model ModulesModel) FindModulesForUser(username string) ([]OutputModule, error) {
	modules := []OutputModule{}

	var rows *sql.Rows
	var err error

	isAdmin := DBPermissions.IsUsernameAnAdmin(username)
	if !isAdmin {
		rows, err = model.DB().Table("user_modules").
			Select(`
				modules.id, level_modules.level, classes.course_id, level_modules.class_id, level_modules.code,
				modules.title, modules.description, modules.color, modules.icon, modules.duration, level_modules.status,
				moduleRole.id 'roleId',
				moduleRole.name 'roleName',
				moduleRole.description 'roleDescription',
				(moduleRole.can_read + coalesce(courseRole.can_read, 0) != 0) 'read',
				(moduleRole.can_write + coalesce(courseRole.can_write, 0) != 0) 'write',
				(moduleRole.can_delete + coalesce(courseRole.can_delete, 0) != 0) 'delete',
				(moduleRole.can_update + coalesce(courseRole.can_update, 0) != 0) 'update',
				users.admin,
				classes.start,
				classes.end,
				course.id, course.title, course.description
			`).Joins(`
				inner join level_modules on level_modules.code = user_modules.module_code
				inner join users on users.id = user_modules.user_id
				inner join modules on modules.id = level_modules.module_id
				inner join roles as moduleRole on moduleRole.id = user_modules.role_id
				inner join classes on classes.id = level_modules.class_id
				inner join courses as course on course.id = classes.course_id
				left outer join user_courses on user_courses.course_id = course.id and user_courses.user_id = user_modules.user_id
				left outer join roles as courseRole on courseRole.id = user_courses.role_id
			`).Where("users.username = ? OR users.admin = 1", username).Rows()
	} else {
		rows, err = model.DB().Table("level_modules").
		Select(`  modules.id, level_modules.level, classes.course_id, level_modules.class_id, level_modules.code,
			  modules.title, modules.description, modules.color, modules.icon, modules.duration, level_modules.status,
			  0 'roleId',
			  'Admin' as roleName,
			  'Admin of a module / course.' as roleDescription,
			  0 'read',
			  0 'write',
			  0 'delete',
			  0 'update',
			  1 'admin',
			  classes.start,
			  classes.end,
			  course.id, course.title, course.description
		`).Joins(`
			inner join modules on level_modules.module_id = modules.id
			inner join classes on classes.id = level_modules.class_id
			inner join courses as course on course.id = classes.course_id
		`).Rows()
	}

	if err != nil {
		return modules, err
	}

	for rows.Next() {
		module := OutputModule{}
		role := OutputRole{}
		course := OutputCourse{}
		classStart := time.Now()
		classEnd := time.Now()

		err = rows.Scan(
			// Module Info
			&module.Id, &module.Level, &module.CourseId, &module.ClassId, &module.Code, &module.Title,
			&module.Description, &module.Color, &module.Icon, &module.Duration, &module.Status,

			// Role Info
			&role.Id, &role.Name, &role.Description, &role.Read, &role.Write, &role.Delete, &role.Update, &role.Admin,

			// Class
			&classStart,
			&classEnd,

			// Course
			&course.Id, &course.Title, &course.Description,
		)
		if err != nil {
			return modules, err
		}

		module.Year = fmt.Sprintf("%d/%d", classStart.Year(), classEnd.Year())
		module.Role = &role
		module.Course = &course
		modules = append(modules, module)
	}

	return modules, nil
}

func (model ModulesModel) FindRawModules (query string, page int) ([]Module, error) {
	modules := []Module{}

	dbQuery := model.DB().Limit(10).Offset(page * 10).Where("modules.title like ? or id = ?", fmt.Sprint("%", query, "%"), query).Find(&modules)
	if dbQuery.Error != nil {
		return modules, dbQuery.Error
	}

	return modules, nil
}


func (model ModulesModel) GetLevelModel(classId, level, moduleId uint32) (*OutputModule, error) {
	var module *OutputModule

	var rows *sql.Rows
	var err error

	rows, err = model.DB().Table("level_modules").
	Select(`
		modules.id, level_modules.level, classes.course_id, level_modules.class_id, level_modules.code,
		modules.title, modules.description, modules.color, modules.icon, modules.duration, level_modules.status,
		classes.start,
		classes.end,
		course.id, course.title, course.description
	`).Joins(`
		inner join modules on modules.id = level_modules.module_id
		inner join classes on classes.id = level_modules.class_id
		inner join courses as course on course.id = classes.course_id
	`).Limit(1).Where("level_modules.level = ? and classes.id = ? and modules.id = ?", level, classId, moduleId).Rows()

	if err != nil {
		return module, err
	}

	// Read only one
	rows.Next()

	module = &OutputModule{}
	course := OutputCourse{}
	classStart := time.Now()
	classEnd := time.Now()

	err = rows.Scan(
		// Module Info
		&module.Id, &module.Level, &module.CourseId, &module.ClassId, &module.Code, &module.Title,
		&module.Description, &module.Color, &module.Icon, &module.Duration, &module.Status,

		// Class
		&classStart,
		&classEnd,

		// Course
		&course.Id, &course.Title, &course.Description,
	)
	if err != nil {
		return module, err
	}

	module.Year = fmt.Sprintf("%d/%d", classStart.Year(), classEnd.Year())
	module.Course = &course
	module.Role = nil

	return module, nil
}

func (model ModulesModel) FindModulesForLevel(classId, level uint32) ([]OutputModule, error) {
	modules := []OutputModule{}

	var rows *sql.Rows
	var err error

	rows, err = model.DB().Table("level_modules").
	Select(`
		modules.id, level_modules.level, classes.course_id, level_modules.class_id, level_modules.code,
		modules.title, modules.description, modules.color, modules.icon, modules.duration, level_modules.status,
		classes.start,
		classes.end,
		course.id, course.title, course.description
	`).Joins(`
		inner join modules on modules.id = level_modules.module_id
		inner join classes on classes.id = level_modules.class_id
		inner join courses as course on course.id = classes.course_id
	`).Where("level_modules.level = ? and classes.id = ?", level, classId).Rows()

	if err != nil {
		return modules, err
	}

	for rows.Next() {
		module := OutputModule{}
		course := OutputCourse{}
		classStart := time.Now()
		classEnd := time.Now()

		err = rows.Scan(
			// Module Info
			&module.Id, &module.Level, &module.CourseId, &module.ClassId, &module.Code, &module.Title,
			&module.Description, &module.Color, &module.Icon, &module.Duration, &module.Status,

			// Class
			&classStart,
			&classEnd,

			// Course
			&course.Id, &course.Title, &course.Description,
		)
		if err != nil {
			return modules, err
		}

		module.Year = fmt.Sprintf("%d/%d", classStart.Year(), classEnd.Year())
		module.Course = &course
		module.Role = nil
		modules = append(modules, module)
	}

	return modules, nil
}

func (model ModulesModel) FindStudentsForModule(moduleCode, roleName string) ([]User, error) {
	var students []User

	query := model.DB().Table("users").Select("distinct users.*").Preload("Avatar").Joins(
		"inner join user_modules on user_modules.user_id = users.id " +
		"inner join roles on roles.id = user_modules.role_id",
	).Where("roles.name = ? and user_modules.module_code = ?", roleName, moduleCode).Find(&students)
	if query.Error != nil {
		return nil, query.Error
	}

	return students, nil
}


func (model ModulesModel) GetModuleStudent(studentId uint32, moduleCode string) (*User, error) {
	var students User

	query := model.DB().Table("users").Select("distinct users.*").Preload("Avatar").Joins(
		"inner join user_modules on user_modules.user_id = users.id",
	).Where("users.id = ? and user_modules.module_code = ?", studentId, moduleCode).First(&students)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.RecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	return &students, nil
}

func (model ModulesModel) AddStudentToModule(moduleCode string, userId uint32) (*UserModule, error) {
	// Preload the level module (so we can get the class ID of the module)
	levelModule, err := DBModule.FindModuleWithCode(moduleCode)
	if err != nil {
		return nil, err
	}

	// Preload the User role named Student
	userRole, err := DBUserRole.FindUserRole("Student")
	if err != nil {
		return nil, err
	}

	userModule := UserModule{
		UserID: userId,
		RoleID: userRole.ID,
		ModuleCode: moduleCode,
		ClassID: levelModule.ClassID,
	}

	query := model.DB().Create(&userModule)
	if query.Error != nil {
		isDuplicated := strings.HasPrefix(query.Error.Error(), "Error 1062")
		if isDuplicated {
			return nil, nil
		}

		return nil, query.Error
	}

	return &userModule, nil
}


func (model ModulesModel) RemoveStudentFromModule(moduleCode string, userId uint32) (int64, error) {
	query := model.DB().
					Table("user_modules").
					Where("user_id = ? and module_code = ?", userId, moduleCode).
					Delete(CourseLevel{})
	if query.Error != nil || query.RowsAffected <= 0 {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}
