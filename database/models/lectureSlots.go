package models

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/YagoCarballo/kumquat.academy.api/database"
)

type LectureSlotsModel struct{}

var DBLectureSlot LectureSlotsModel

func (model LectureSlotsModel) DB() *gorm.DB {
	return database.DB
}

func (model LectureSlotsModel) CreateLectureSlot(moduleId uint32, location, lecType string, start, end time.Time) (*LectureSlot, error) {
	lectureSlot := LectureSlot{
		ModuleID: moduleId,
		Location: location,
		Type: lecType,
		Start: start,
		End: end,
	}

	query := model.DB().Create(&lectureSlot)
	if query.Error != nil {
		return nil, query.Error
	}

	return &lectureSlot, nil
}

func (model LectureSlotsModel) ReadLectureSlot(lectureSlotId uint32) (*LectureSlot, error) {
	var lectureSlot LectureSlot

	query := model.DB().Where("id = ?", lectureSlotId).First(&lectureSlot)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.RecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	return &lectureSlot, nil
}

func (model LectureSlotsModel) UpdateLectureSlot(lectureSlotId, moduleId uint32, location, lecType string, start, end time.Time) (*LectureSlot, error) {
	lectureSlot := LectureSlot{}

	query := model.DB().Table("lecture_slots").Where("id = ?", lectureSlotId).Updates(map[string]interface{}{
		"location": location,
		"type": lecType,
		"start": start,
		"end": end,
	})
	if query.Error != nil {
		return nil, query.Error
	}

	if query.RowsAffected <= 0 {
		return nil, nil
	}

	query = model.DB().Where("id = ?", lectureSlotId).First(&lectureSlot)
	if query.Error != nil {
		return nil, query.Error
	}

	return &lectureSlot, nil
}

func (model LectureSlotsModel) DeleteLectureSlot(lectureSlotId uint32) (int64, error) {
	query := model.DB().Table("lecture_slots").Where("id = ?", lectureSlotId).Delete(LectureSlot{})
	if query.Error != nil {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}

func (model LectureSlotsModel) FindLectureSlotsForModule(moduleCode string) (map[string][]LectureSlot, error) {
	var lectureSlots []LectureSlot

	query := model.DB().Table("lecture_slots").Joins(
		"inner join level_modules on level_modules.module_id = lecture_slots.module_id",
	).Where("level_modules.code = ?", moduleCode).Find(&lectureSlots)
	if query.Error != nil {
		return nil, query.Error
	}

	// Create an empty list per day of the Week
	lectureSlotsMap := map[string][]LectureSlot{
		time.Monday.String(): 	[]LectureSlot{},
		time.Tuesday.String(): 	[]LectureSlot{},
		time.Wednesday.String():[]LectureSlot{},
		time.Thursday.String():	[]LectureSlot{},
		time.Friday.String():	[]LectureSlot{},
		time.Saturday.String():	[]LectureSlot{},
		time.Sunday.String():	[]LectureSlot{},
	}

	// Group the lecture Slots in days of the Week
	for _, slot := range lectureSlots {
		weekDay := slot.Start.Weekday().String()
		lectureSlotsMap[weekDay] = append(lectureSlotsMap[weekDay], slot)
	}

	return lectureSlotsMap, nil
}
