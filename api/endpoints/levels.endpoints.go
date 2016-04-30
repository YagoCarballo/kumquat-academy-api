package endpoints

import (
	"net/http"
	"github.com/YagoCarballo/kumquat.academy.api/database/models"
	"fmt"
	"strings"
	"time"
)

// CRUD

func CreateLevel(lvl, classId, courseId uint32, start, end time.Time) (int, map[string]interface{}) {
	level, err := models.DBLevel.CreateLevel(courseId, classId, lvl, start, end)
	if err != nil {
		dbError := err.Error()
		errorCode := "Unknown"
		errorMessage := "Error creating the level."
		duplicated := strings.HasPrefix(dbError, "Error 1062")
		if duplicated {
			errorCode = "Duplicated"
			errorMessage = "There is already a level in that class."
		}

		return http.StatusExpectationFailed, map[string]interface{}{
			"error": errorCode,
			"message": errorMessage,
		}
	}

	return http.StatusCreated, map[string]interface{}{
		"level": level,
	}
}

func GetLevel(courseId, classId, lvl uint32) (int, map[string]interface{}) {
	level, err := models.DBLevel.ReadLevel(courseId, classId, lvl)
	if err != nil || level == nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "NotFound",
			"message": "Level not found.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"level": level,
	}
}

func UpdateLevel(lvl, classId, courseId uint32, start, end time.Time) (int, map[string]interface{}) {
	level, err := models.DBLevel.UpdateLevel(lvl, classId, courseId, start, end)
	if err != nil || level == nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Level not updated.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"level": level,
	}
}

func DeleteLevel(courseId, classId, lvl uint32) (int, map[string]interface{}) {
	rows, err := models.DBLevel.DeleteLevel(courseId, classId, lvl)
	if err != nil || rows <= 0 {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error deleting the Level",
		}
	}

	return http.StatusAccepted, map[string]interface{}{
		"message": fmt.Sprintf("Level %d removed from Class %d", lvl, classId),
	}
}

func AddModuleToLevel(code string, lvl, classId, moduleId uint32, start time.Time) (int, map[string]interface{}) {
	levelModule, err := models.DBLevel.AddModule(code, lvl, classId, moduleId, start)
	if err != nil {
		dbError := err.Error()
		errorCode := "Unknown"
		errorMessage := "Error adding the module to the level."
		duplicated := strings.HasPrefix(dbError, "Error 1062")
		if duplicated {
			errorCode = "Duplicated"
			errorMessage = "That module is already in that level"
		}

		return http.StatusExpectationFailed, map[string]interface{}{
			"error": errorCode,
			"message": errorMessage,
		}
	}

	return http.StatusCreated, map[string]interface{}{
		"module": levelModule,
	}
}
