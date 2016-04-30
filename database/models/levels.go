package models

import (
	"github.com/jinzhu/gorm"
	"github.com/YagoCarballo/kumquat-academy-api/database"
	"time"
)

type LevelsModel struct{}
var DBLevel LevelsModel

func (model LevelsModel) DB() *gorm.DB {
	return database.DB
}

func (model LevelsModel) CreateLevel(courseId uint32, classId, lvl uint32, start, end time.Time) (*CourseLevel, error) {
	level := CourseLevel{
		Level: lvl,
		CourseID: courseId,
		ClassID: classId,
		Start: start,
		End: end,
	}

	query := model.DB().Create(&level)
	if query.Error != nil {
		return nil, query.Error
	}

	return &level, nil
}


func (model LevelsModel) ReadLevel(courseId, classId, lvl uint32) (*CourseLevel, error) {
	var level CourseLevel

	query := model.DB().Where("level = ? AND class_id = ? AND course_id = ?", lvl, classId, courseId).First(&level)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.RecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	return &level, nil
}

func (model LevelsModel) UpdateLevel(courseId, classId, lvl uint32, start, end time.Time) (*CourseLevel, error) {
	level := CourseLevel{
		Level: lvl,
		Start: start,
		End: end,
	}
	query := model.DB().Table("course_levels").Where("level = ?", lvl).Update(&level)
	if query.Error != nil {
		return nil, query.Error
	}

	if query.RowsAffected <= 0 {
		return nil, nil
	}

	return &level, nil
}


func (model LevelsModel) DeleteLevel(courseId, classId, level uint32) (int64, error) {
	query := model.DB().
			Table("course_levels").
			Where("level = ? AND class_id = ? AND course_id = ?", level, classId, courseId).
			Delete(CourseLevel{})
	if query.Error != nil {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}

func (model LevelsModel) AddModule(code string, lvl, classId, moduleId uint32, start time.Time) (*OutputModule, error) {
	levelModule := LevelModule{
		Code: code,
		Level: lvl,
		ClassID: classId,
		ModuleID: moduleId,
		Status: ModuleDraft,
		Start: start,
	}

	query := model.DB().Create(&levelModule)
	if query.Error != nil {
		return nil, query.Error
	}

	// Get the full module
	module, err := DBModule.GetLevelModel(classId, lvl, moduleId)
	if err != nil {
		return nil, err
	}

	return module, nil
}


func (model LevelsModel) RemoveModule(code string, classId uint32) (int64, error) {
	query := model.DB().
		Table("level_modules").
		Where("code = ? AND class_id = ?", code, classId).
		Delete(LevelModule{})
	if query.Error != nil {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}
