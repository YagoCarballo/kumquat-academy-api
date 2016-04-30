package endpoints

import (
	"net/http"
	"github.com/YagoCarballo/kumquat-academy-api/database/models"
	"fmt"
	"strings"
)

// CRUD

func CreateCourse(title, description string) (int, map[string]interface{}) {
	course, err := models.DBCourse.CreateCourse(title, description)
	if err != nil {
		dbError := err.Error()
		errorCode := "Unknown"
		errorMessage := "Error creating the course."
		duplicated := strings.HasPrefix(dbError, "Error 1062")
		if duplicated {
			errorCode = "Duplicated"
			errorMessage = "There is already a course with that title."
		}

		return http.StatusExpectationFailed, map[string]interface{}{
			"error": errorCode,
			"message": errorMessage,
		}
	}

	return http.StatusCreated, map[string]interface{}{
		"course": course,
	}
}

func GetCourse(id uint32) (int, map[string]interface{}) {
	course, err := models.DBCourse.ReadCourse(id)
	if err != nil || course == nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "NotFound",
			"message": "Course not found.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"course": course,
	}
}

func GetCourseWithTitle(title string) (int, map[string]interface{}) {
	course, err := models.DBCourse.GetCourseWithTitle(title)
	if err != nil || course == nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "NotFound",
			"message": "Course not found.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"course": course,
	}
}

func UpdateCourse(id uint32, title, description string) (int, map[string]interface{}) {
	course, err := models.DBCourse.UpdateCourse(id, title, description)
	if err != nil || course == nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Course not updated.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"course": course,
	}
}

func DeleteCourse(id uint32) (int, map[string]interface{}) {
	rows, err := models.DBCourse.DeleteCourse(id)
	if err != nil || rows <= 0 {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error deleting the Course",
		}
	}

	return http.StatusAccepted, map[string]interface{}{
		"message": fmt.Sprintf("Course %d removed", id),
	}
}

func GetCoursesForUser(userId uint32) (int, map[string]interface{}) {
	courses, err := models.DBCourse.FindCoursesForUser(userId)
	if err != nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Unknown",
			"message": "Error fetching the courses for that user.",
		}
	}

	for courseIndex, course := range courses {
		for moduleIndex, module := range course.Modules {
			permissions, err := models.DBPermissions.GetPermissionsForModule(&userId, &module.Id, nil)
			if err != nil {
				courses[courseIndex].Modules[moduleIndex].Role = &models.OutputRole{}
				continue
			}

			courses[courseIndex].Modules[moduleIndex].Role = &models.OutputRole{
				Id: permissions.RoleId,
				Name: permissions.RoleName,
				Description: permissions.RoleDescription,
				Admin: permissions.Admin,
				Read: permissions.Read,
				Write: permissions.Write,
				Delete: permissions.Delete,
				Update: permissions.Update,
			}
		}
	}

	return http.StatusOK, map[string]interface{}{
		"courses": courses,
	}
}


