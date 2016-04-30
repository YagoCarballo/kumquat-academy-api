package models

import (
	"github.com/jinzhu/gorm"
	"github.com/YagoCarballo/kumquat-academy-api/database"
	"fmt"
)

type CoursesModel struct {}
var DBCourse CoursesModel

func (model CoursesModel) DB () *gorm.DB {
	return database.DB
}

func (model CoursesModel) CreateCourse(title, description string) (*Course, error) {
	course := Course{
		Title: title,
		Description: description,
	}

	query := model.DB().Create(&course)
	if query.Error != nil {
		return nil, query.Error
	}

	return &course, nil
}

func (model CoursesModel) ReadCourse(id uint32) (*Course, error) {
	var course Course

	query := model.DB().Where("id = ?", id).First(&course)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.RecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	return &course, nil
}

func (model CoursesModel) UpdateCourse(id uint32, title, description string) (*Course, error) {
	course := Course{
		Title: title,
		Description: description,
	}

	query := model.DB().Table("courses").Where("id = ?", id).Update(&course)
	if query.Error != nil {
		return nil, query.Error
	}

	if query.RowsAffected <= 0 {
		return nil, nil
	}

	course.ID = id
	return &course, nil
}

func (model CoursesModel) DeleteCourse(id uint32) (int64, error) {
	query := model.DB().Where("id = ?", id).Delete(Course{})
	if query.Error != nil {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}

func (model CoursesModel) GetCourseWithTitle(title string) (*Course, error) {
	var course Course

	query := model.DB().Where("title = ?", title).First(&course)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.RecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	return &course, nil
}

func (model CoursesModel) FindCoursesForUser(userId uint32) ([]Course, error) {
	courses := []Course{}

	isAdmin := DBPermissions.IsAdmin(userId)
	if !isAdmin {
		query := model.DB().Raw(`
			select distinct courses.*
			from courses
				inner join classes on courses.id = classes.course_id
				inner join level_modules on classes.course_id = courses.id
				inner join modules on modules.id = level_modules.module_id
				inner join user_modules on user_modules.module_code = level_modules.code
			where user_modules.user_id = ?
				union
			select distinct courses.*
			from courses
				inner join user_courses on user_courses.course_id = courses.id
				inner join classes on user_courses.course_id = classes.id
				inner join level_modules on user_courses.course_id = courses.id
			where user_courses.user_id = ?
		`, userId, userId).Find(&courses)
		if query.Error != nil {
			return courses, query.Error
		}
	} else {
		query := model.DB().Raw(`
			select distinct courses.*
			from courses`).Find(&courses)
		if query.Error != nil {
			return courses, query.Error
		}
	}


	for courseIndex, course := range courses {
		courses[courseIndex].Role, _ = DBPermissions.GetPermissionsForCourse(userId, course.ID)
		levelModules, lvlErr := model.FindLevelModulesForCourseAndUser(course.ID, userId, isAdmin)
		if lvlErr != nil {
			courses[courseIndex].LevelModules = []LevelModule{}
			continue;
		}

		if levelModules != nil {
			for _, levelModule := range levelModules {
				courses[courseIndex].Modules = append(courses[courseIndex].Modules, &OutputModule{
					Id: levelModule.Module.ID,
					Level: levelModule.Level,
					CourseId: course.ID,
					ClassId: levelModule.ClassID,
					Code: levelModule.Code,
					Title: levelModule.Module.Title,
					Description: levelModule.Module.Description,
					Color: levelModule.Module.Color,
					Icon: levelModule.Module.Icon,
					Duration: levelModule.Module.Duration,
					Year: fmt.Sprintf("%d/%d", levelModule.Class.Start.Year(), levelModule.Class.End.Year()),
					Status: levelModule.Status,
				})
			}
			courses[courseIndex].LevelModules = levelModules
		}

	}

	return courses, nil
}


func (model CoursesModel) FindLevelModulesForCourseAndUser(courseId, userId uint32, isAdmin bool) ([]LevelModule, error) {
	levels := []LevelModule{}

	if isAdmin == false {
		query := model.DB().Preload("Module").Preload("Class").Raw(`
			select distinct level_modules.* from level_modules
			inner join course_levels on course_levels.class_id = level_modules.class_id
			left outer join classes on classes.id = level_modules.class_id
			left outer join user_modules on user_modules.class_id = level_modules.class_id and user_modules.module_code = level_modules.code
			left outer join user_courses on user_courses.course_id = classes.course_id
			where classes.course_id = ? and (user_modules.user_id = ? or user_courses.user_id = ?)
		`, courseId, userId, userId).Find(&levels)
		if query.Error != nil {
			return levels, query.Error
		}
	} else {
		query := model.DB().Preload("Module").Preload("Class").Raw(`
			select distinct level_modules.* from level_modules
			inner join course_levels on course_levels.class_id = level_modules.class_id
			where course_levels.course_id = ?
		`, courseId).Find(&levels)
		if query.Error != nil {
			return levels, query.Error
		}
	}

	return levels, nil
}

