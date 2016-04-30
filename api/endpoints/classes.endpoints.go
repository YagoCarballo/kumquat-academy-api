package endpoints

import (
	"fmt"
	"strings"
	"net/http"
	"github.com/YagoCarballo/kumquat-academy-api/database/models"
	"time"
)

// CRUD

func CreateClass(courseId uint32, title string, start, end time.Time, levels []*models.CourseLevel) (int, map[string]interface{}) {
	class, err := models.DBClass.CreateClass(courseId, title, start, end, levels)
	if err != nil {
		dbError := err.Error()
		errorCode := "Unknown"
		errorMessage := "Error creating the class."
		duplicated := strings.HasPrefix(dbError, "Error 1062")
		if duplicated {
			errorCode = "Duplicated"
			errorMessage = "There is already a class with that title."
		}

		return http.StatusExpectationFailed, map[string]interface{}{
			"error": errorCode,
			"message": errorMessage,
		}
	}

	levelErrors := []map[string]interface{}{}
	for index, level := range levels {
		dbLevel, err := models.DBLevel.CreateLevel(courseId, class.ID, level.Level, level.Start, level.End)
		if err != nil {
			dbError := err.Error()
			errorCode := "Unknown"
			errorMessage := "Error creating the level."
			duplicated := strings.HasPrefix(dbError, "Error 1062")
			if duplicated {
				errorCode = "Duplicated"
				errorMessage = fmt.Sprintf("This class already has the level %d", level.Level)
			}

			levelErrors = append(levelErrors, map[string]interface{}{
				"error": errorCode,
				"message": errorMessage,
			});

			continue;
		}

		// Update the level int he class
		class.Levels[index] = dbLevel
	}

	return http.StatusCreated, map[string]interface{}{
		"class": class,
		"errors": levelErrors,
	}
}

func GetClass(id uint32) (int, map[string]interface{}) {
	class, err := models.DBClass.ReadClass(id)
	if err != nil || class == nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "NotFound",
			"message": "Class not found.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"class": class,
	}
}

func GetClassForYear(courseId uint32, year string) (int, map[string]interface{}) {
	class, err := models.DBClass.GetClassWithTitle(courseId, year)
	if err != nil || class == nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "NotFound",
			"message": "Class not found.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"class": class,
	}
}

func GetClassesForCourse(courseId uint32) (int, map[string]interface{}) {
	class, err := models.DBClass.GetClassesForCourse(courseId)
	if err != nil || class == nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "NotFound",
			"message": "Class not found.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"classes": class,
	}
}

func UpdateClass(id, courseId uint32, title string, start, end time.Time) (int, map[string]interface{}) {
	class, err := models.DBClass.UpdateClass(id, courseId, title, start, end)
	if err != nil || class == nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Class not updated.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"class": class,
	}
}

func DeleteClass(id uint32) (int, map[string]interface{}) {
	rows, err := models.DBClass.DeleteClass(id)
	if err != nil || rows <= 0 {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error deleting the Class",
		}
	}

	return http.StatusAccepted, map[string]interface{}{
		"message": fmt.Sprintf("Class %d removed", id),
	}
}
