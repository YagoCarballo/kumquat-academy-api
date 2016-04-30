package endpoints

import (
	"fmt"
	"time"
	"net/http"
	"github.com/YagoCarballo/kumquat.academy.api/database/models"
)

// CRUD

func CreateLectureSlot(moduleCode, location, lecType string, start, end time.Time) (int, map[string]interface{}) {
	module, err := models.DBModule.FindModuleWithCode(moduleCode)
	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error reading the module",
		}
	}

	lectureSlot, err := models.DBLectureSlot.CreateLectureSlot(module.ModuleID, location, lecType, start, end)
	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error creating the lectureSlot.",
		}
	}

	return http.StatusCreated, map[string]interface{}{
		"lecture_slot": lectureSlot,
	}
}

func GetLectureSlot(lectureSlotId uint32) (int, map[string]interface{}) {
	lectureSlot, err := models.DBLectureSlot.ReadLectureSlot(lectureSlotId)
	if err != nil || lectureSlot == nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "NotFound",
			"message": "Lecture slot not found.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"lecture_slot": lectureSlot,
	}
}

func UpdateLectureSlot(lectureSlotId uint32, moduleCode, location, lecType string, start, end time.Time) (int, map[string]interface{}) {
	module, err := models.DBModule.FindModuleWithCode(moduleCode)
	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error reading the module",
		}
	}

	lectureSlot, err := models.DBLectureSlot.UpdateLectureSlot(lectureSlotId, module.ModuleID, location, lecType, start, end)
	if err != nil || lectureSlot == nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Lecture slot not updated.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"lecture_slot": lectureSlot,
	}
}

func DeleteLectureSlot(lectureSlotId uint32) (int, map[string]interface{}) {
	rows, err := models.DBLectureSlot.DeleteLectureSlot(lectureSlotId)
	if err != nil || rows <= 0 {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error deleting the Lecture slot",
		}
	}

	return http.StatusAccepted, map[string]interface{}{
		"message": fmt.Sprintf("Lecture slot %d removed.", lectureSlotId),
	}
}

func FindLectureSlotsForModule(moduleId string) (int, map[string]interface{}) {
	lectureSlots, err := models.DBLectureSlot.FindLectureSlotsForModule(moduleId)
	if err != nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "Unknown",
			"message": "Error fetching the lecture slots.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"lecture_slots": lectureSlots,
	}
}
