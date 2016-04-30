package endpoints

import (
	"fmt"
	"time"
	"net/http"
	"github.com/YagoCarballo/kumquat.academy.api/database/models"
	"github.com/YagoCarballo/kumquat.academy.api/tools"
	"mime/multipart"
	"os"
)

// CRUD

func CreateLecture(moduleCode, location, topic, description string, start, end time.Time, canceled bool, lectureSlotId *uint32) (int, map[string]interface{}) {
	module, err := models.DBModule.FindModuleWithCode(moduleCode)
	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error reading the module",
		}
	}

	lecture, err := models.DBLecture.CreateLecture(module.ModuleID, location, topic, description, start, end, canceled, lectureSlotId)
	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error creating the lecture.",
		}
	}

	return http.StatusCreated, map[string]interface{}{
		"lecture": lecture,
	}
}

func GetLecture(lectureId uint32) (int, map[string]interface{}) {
	lecture, err := models.DBLecture.ReadLecture(lectureId)
	if err != nil || lecture == nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "NotFound",
			"message": "Lecture not found.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"lecture": lecture,
	}
}

func UpdateLecture(lectureId uint32, moduleCode, location, topic, description string, start, end time.Time, canceled bool, lectureSlotId *uint32) (int, map[string]interface{}) {
	module, err := models.DBModule.FindModuleWithCode(moduleCode)
	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error reading the module",
		}
	}

	lecture, err := models.DBLecture.UpdateLecture(lectureId, module.ModuleID, location, topic, description, start, end, canceled, lectureSlotId)
	if err != nil || lecture == nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Lecture not updated.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"lecture": lecture,
	}
}

func DeleteLecture(lectureId uint32) (int, map[string]interface{}) {
	rows, err := models.DBLecture.DeleteLecture(lectureId)
	if err != nil || rows <= 0 {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error deleting the Lecture",
		}
	}

	return http.StatusAccepted, map[string]interface{}{
		"message": fmt.Sprintf("Lecture %d removed.", lectureId),
	}
}

func FindLecturesForModule(moduleId string) (int, map[string]interface{}) {
	lectures, err := models.DBLecture.FindLecturesForModule(moduleId)
	if err != nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "Unknown",
			"message": "Error fetching the lectures.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"lectures": lectures,
	}
}

func FindLectureWeeksForModule(moduleId string) (int, map[string]interface{}) {
	lectures, err := models.DBLecture.FindLecturesWeeksForModule(moduleId)
	if err != nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "Unknown",
			"message": "Error fetching the lectures.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"lectures": lectures,
	}
}

func FindLectureWeeksAndSlotsForModule(moduleCode string) (int, map[string]interface{}) {
	// Create the Weeks Array
	weeks := []map[string]interface{}{}

	// Get the GMT Timezone to use as base
	gmt := time.FixedZone("GMT", 0)

	// Get the Module Information
	levelModule, err := models.DBModule.FindModuleWithCode(moduleCode)
	if err != nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "Unknown",
			"message": "Error fetching the module",
		}
	}

	// Get Lecture Slots
	lectureSlots, err := models.DBLectureSlot.FindLectureSlotsForModule(moduleCode)
	if err != nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "Unknown",
			"message": "Error fetching the lecture slots.",
		}
	}

	// Get the duration of the Day
	dayDuration, err := time.ParseDuration("23h59m")
	if err != nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "Unknown",
			"message": "Error parsing dates.",
		}
	}

	// Get the Module Start Date
	startYear, startWeek := levelModule.Start.ISOWeek()
	parsedStart := tools.FirstDayOfISOWeek(startYear, startWeek, gmt)

	// Get the Course End Year
	endYear, endWeek := levelModule.Start.AddDate(0, 0, int(levelModule.Module.Duration * 7)).ISOWeek()
	parsedEnd := tools.FirstDayOfISOWeek(endYear, endWeek, gmt).AddDate(0, 0, 6)

	// Loop through each week
	count := 0
	for current := parsedStart; current.Before(parsedEnd); current = current.AddDate(0, 0, 7) {
		count++;

		// Get the start of the Day and end of the Day
		weekStart := current
		weekEnd := current.AddDate(0, 0, 6).Add(dayDuration)

		// Get the lectures for the week and group them into week days
		lectures, _ := models.DBLecture.FindLecturesForModuleInRange(moduleCode, &weekStart, &weekEnd, -1);
		lecturesMap := models.GroupLecturesInWeeks(&lectures)

		weeks = append(weeks, map[string]interface{}{
			"week": count,
			"start": weekStart,
			"end": weekEnd,
			"lectures": lecturesMap,
		})

		// Once the duration is reached, Stop
		if count >= int(levelModule.Module.Duration) {
			break;
		}
	}

	return http.StatusOK, map[string]interface{}{
		"slots": lectureSlots,
		"weeks": weeks,
	}
}


func UploadLectureAttachments(lectureId uint32, file multipart.File, header *multipart.FileHeader) (int, FileResponseMessage) {
	status, response := UploadFile(file, header)
	if status != http.StatusOK {
		return http.StatusExpectationFailed, FileResponseMessage{
			Error: "Unknown",
			Message: "Error uploading the file.",
		}
	}

	count, err := models.DBLecture.AddAttachmentToLecture(lectureId, response.Attachment.ID)
	if err != nil || count <= 0 {
		return http.StatusExpectationFailed, FileResponseMessage{
			Error: "Unknown",
			Message: "Error adding the attachment to the lecture.",
		}
	}

	return status, response
}

func RemoveLectureAttachments(lectureId, attachmentId uint32) (int, map[string]interface{}) {
	var err error

	attachment, _ := models.DBAttachment.ReadAttachment(attachmentId)
	if attachment != nil {
		serverSettings := tools.GetSettings().Server
		path := fmt.Sprintf("%s/%s", serverSettings.UploadsPath, attachment.Url)
		os.Remove(path)
	}

	count, err := models.DBLecture.RemoveAttachmentFromLecture(lectureId, attachmentId)
	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error deleting the attachment from the lecture.",
		}
	}

	if count <= 0 {
		return http.StatusNotFound, map[string]interface{}{
			"error": "NotFound",
			"message": "An attachment with that Id was not found inside the lecture.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"message": "Attachment removed.",
	}
}


func FindLectureWeeksForUser(userId uint32) (int, map[string]interface{}) {
	// Get the GMT Timezone to use as base
	gmt := time.FixedZone("GMT", 0)

	// Get the start of the current Week
	startYear, startWeek := time.Now().ISOWeek()
	parsedStart := tools.FirstDayOfISOWeek(startYear, startWeek, gmt)

	// Get the end of the current Week
	parsedEnd := parsedStart.AddDate(0, 0, 6)

	// Get the lectures for the week and group them into week days
	lectures, _ := models.DBLecture.FindLecturesForUserInRange(userId, &parsedStart, &parsedEnd, -1);
	lecturesMap := models.GroupLecturesInWeeks(&lectures)

	return http.StatusOK, map[string]interface{}{
		"schedule": lecturesMap,
	}
}
