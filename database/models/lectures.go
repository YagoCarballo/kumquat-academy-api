package models

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/YagoCarballo/kumquat-academy-api/database"
	"fmt"
)

type LecturesModel struct{}

var DBLecture LecturesModel

func (model LecturesModel) DB() *gorm.DB {
	return database.DB
}

func (model LecturesModel) CreateLecture(moduleId uint32, location, topic, description string, start, end time.Time, canceled bool, lectureSlotId *uint32) (*Lecture, error) {
	lecture := Lecture{
		Description: description,
		ModuleID: moduleId,
		Location: location,
		Topic: topic,
		Start: start,
		End: end,
		Canceled: canceled,
	}

	if lectureSlotId != nil {
		lecture.LectureSlotID = lectureSlotId
	}

	query := model.DB().Create(&lecture)
	if query.Error != nil {
		return nil, query.Error
	}

	return &lecture, nil
}

func (model LecturesModel) ReadLecture(lectureId uint32) (*Lecture, error) {
	var lecture Lecture

	query := model.DB().Preload("LectureSlot").Preload("Attachments").Where("id = ?", lectureId).First(&lecture)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	return &lecture, nil
}

func (model LecturesModel) UpdateLecture(lectureId, moduleId uint32, location, topic, description string, start, end time.Time, canceled bool, lectureSlotId *uint32) (*Lecture, error) {
	lecture := Lecture{}

	query := model.DB().Table("lectures").Where("id = ?", lectureId).Updates(map[string]interface{}{
		"description": description,
		"location": location,
		"topic": topic,
		"start": start,
		"end": end,
		"canceled": canceled,
		"lecture_slot_id": lectureSlotId,
	})
	if query.Error != nil {
		return nil, query.Error
	}

	if query.RowsAffected <= 0 {
		return nil, nil
	}

	query = model.DB().Where("id = ?", lectureId).First(&lecture)
	if query.Error != nil {
		return nil, query.Error
	}

	return &lecture, nil
}

func (model LecturesModel) DeleteLecture(lectureId uint32) (int64, error) {
	query := model.DB().Table("lectures").Where("id = ?", lectureId).Delete(Lecture{})
	if query.Error != nil {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}

func (model LecturesModel) FindLecturesForModule(moduleCode string) ([]Lecture, error) {
	var lectures []Lecture

	query := model.DB().Table("lectures").Preload("LectureSlot").Preload("Attachments").Joins(
		"inner join level_modules on level_modules.module_id = lectures.module_id",
	).Where("level_modules.code = ?", moduleCode).Find(&lectures)
	if query.Error != nil {
		return nil, query.Error
	}

	return lectures, nil
}

func (model LecturesModel) FindLecturesWeeksForModule(moduleCode string) (map[string][]Lecture, error) {
	var lectures []Lecture

	query := model.DB().Table("lectures").Preload("LectureSlot").Preload("Attachments").Joins(
		"inner join level_modules on level_modules.module_id = lectures.module_id",
	).Where("level_modules.code = ?", moduleCode).Find(&lectures)
	if query.Error != nil {
		return nil, query.Error
	}

	weeks := map[string][]Lecture{}
	for _, lecture := range lectures {
		year, week := lecture.Start.ISOWeek()
		key := fmt.Sprintf("%d-%d", year, week)

		if weeks[key] == nil {
			weeks[key] = []Lecture{}
		}

		weeks[key] = append(weeks[key], lecture);
	}

	return weeks, nil
}

func (model LecturesModel) FindLecturesForModuleInRange(moduleCode string, start, end *time.Time, limit int) ([]Lecture, error) {
	var lectures []Lecture

	query := model.DB().Table("lectures").Preload("LectureSlot").Preload("Attachments").Limit(limit).Order("start").Joins(
		"inner join level_modules on level_modules.module_id = lectures.module_id",
	).Where("level_modules.code = ? and lectures.start >= ? and lectures.start <= ?", moduleCode, start, end).Find(&lectures)
	if query.Error != nil {
		return nil, query.Error
	}

	return lectures, nil
}

func (model LecturesModel) FindLecturesForUserInRange(userId uint32, start, end *time.Time, limit int) ([]Lecture, error) {
	var lectures []Lecture

	query := model.DB().Table("lectures").Preload("LectureSlot").Preload("Attachments").Preload("Module").Limit(limit).Order("start").Joins(
		"inner join level_modules on level_modules.module_id = lectures.module_id " +
		"inner join user_modules on user_modules.module_code = level_modules.code",
	).Where("user_modules.user_id = ? and lectures.start >= ? and lectures.start <= ? and canceled = 0", userId, start, end).Find(&lectures)
	if query.Error != nil {
		return nil, query.Error
	}

	return lectures, nil
}

func GroupLecturesInWeeks(lectures *[]Lecture) map[string][]Lecture {
	// Create an empty list per day of the Week
	week := map[string][]Lecture{
		time.Monday.String(): 	[]Lecture{},
		time.Tuesday.String(): 	[]Lecture{},
		time.Wednesday.String():[]Lecture{},
		time.Thursday.String():	[]Lecture{},
		time.Friday.String():	[]Lecture{},
		time.Saturday.String():	[]Lecture{},
		time.Sunday.String():	[]Lecture{},
	}

	// Group the lectures in days of the Week
	for _, lecture := range *lectures {
		weekDay := lecture.Start.Weekday().String()
		week[weekDay] = append(week[weekDay], lecture)
	}

	return week
}

func (model LecturesModel) AddAttachmentToLecture(lectureId, attachmentId uint32) (int64, error) {
	lectureAttachment := LectureAttachments{
		LectureID: lectureId,
		AttachmentID: attachmentId,
	}

	query := model.DB().Create(&lectureAttachment)
	if query.Error != nil {
		return query.RowsAffected, query.Error
	}

	return query.RowsAffected, nil
}

func (model LecturesModel) RemoveAttachmentFromLecture(lectureId, attachmentId uint32) (int64, error) {
	query := model.DB().Table("lecture_attachments").
				Where("lecture_id = ? and attachment_id = ?", lectureId, attachmentId).
				Delete(LectureAttachments{})
	if query.Error != nil {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}

