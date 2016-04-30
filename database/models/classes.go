package models

import (
	"github.com/jinzhu/gorm"
	"github.com/YagoCarballo/kumquat.academy.api/database"
	"time"
)

type ClassesModel struct{}
var DBClass ClassesModel

func (model ClassesModel) DB() *gorm.DB {
	return database.DB
}

func (model ClassesModel) CreateClass(courseId uint32, title string, start, end time.Time, levels []*CourseLevel) (*Class, error) {
	class := Class{
		Title: title,
		CourseID: courseId,
		Start: start,
		End: end,
		Levels: levels,
	}

	query := model.DB().Create(&class)
	if query.Error != nil {
		return nil, query.Error
	}

	return &class, nil
}

func (model ClassesModel) ReadClass(id uint32) (*Class, error) {
	var class Class
	levels := []*CourseLevel{}

	query := model.DB().Where("id = ?", id).First(&class)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.RecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	// Fetch the levels for this class
	query = model.DB().Where("class_id = ? AND course_id = ?", class.ID, class.CourseID).Find(&levels)
	if query.Error != nil {
		// If no levels found, continue without errors
	}

	class.Levels = levels
	return &class, nil
}

func (model ClassesModel) UpdateClass(id, courseId uint32, title string, start, end time.Time) (*Class, error) {
	class := Class{
		Title: title,
		CourseID: courseId,
		Start: start,
		End: end,
	}

	query := model.DB().Table("classes").Where("id = ?", id).Update(&class)
	if query.Error != nil {
		return nil, query.Error
	}

	if query.RowsAffected <= 0 {
		return nil, nil
	}

	class.ID = id
	levels := []*CourseLevel{}

	// Fetch the levels for this class
	query = model.DB().Where("class_id = ? AND course_id = ?", class.ID, class.CourseID).Find(&levels)
	if query.Error != nil {
		// If no levels found, continue without errors
	}

	class.Levels = levels
	return &class, nil
}

func (model ClassesModel) DeleteClass(id uint32) (int64, error) {
	query := model.DB().Table("classes").Where("id = ?", id).Delete(Class{})
	if query.Error != nil {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}

func (model ClassesModel) GetClassWithTitle(courseId uint32, year string) (*Class, error) {
	// Creates empty Class
	var class Class
	levels := []*CourseLevel{}

	query := model.DB().First(&class, "title = ? AND course_id = ?", year, courseId)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.RecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	// Fetch the levels for this class
	query = model.DB().Where("class_id = ? AND course_id = ?", class.ID, class.CourseID).Find(&levels)
	if query.Error != nil {
		// If no levels found, continue without errors
	}

	class.Levels = levels
	return &class, nil
}


func (model ClassesModel) GetClassesForCourse(courseId uint32) ([]Class, error) {
	// Creates empty Class
	var classes []Class

	query := model.DB().Find(&classes, "course_id = ?", courseId)
	if query.Error != nil {
		return classes, query.Error
	}

	for index, class := range classes  {
		levels := []*CourseLevel{}

		// Fetch the levels for this class
		query = model.DB().Where("class_id = ? AND course_id = ?", class.ID, class.CourseID).Find(&levels)
		if query.Error != nil {
			// If no levels found, continue without errors
		}

		classes[index].Levels = levels
	}

	return classes, nil
}
