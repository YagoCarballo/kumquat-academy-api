package endpoints

import (
	"net/http"
	"github.com/YagoCarballo/kumquat.academy.api/database/models"
)

func GetModulesForUser(username string) (int, map[string]interface{}) {
	modules, err := models.DBModule.FindModulesForUser(username)
	if err != nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Error fetching the modules for that user.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"modules": modules,
	}
}

func CreateModule(title, description, icon, color string, duration uint32) (int, map[string]interface{}) {
	module := models.Module{
		Title: title,
		Description: description,
		Icon: icon,
		Color: color,
		Duration: duration,
	}

	dbModule, err := models.DBModule.CreateModule(module)
	if err != nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Error creating the module.",
		}
	}

	return http.StatusCreated, map[string]interface{}{
		"message": "Module created successfully",
		"module": dbModule,
	}
}

func FindRawModules(query string, page int) (int, map[string]interface{}) {
	modules, err := models.DBModule.FindRawModules(query, page)
	if err != nil {
		return http.StatusPreconditionFailed, map[string]interface{}{
			"error": "Error fetching the list of modules.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"modules": modules,
	}
}

func GetModulesForLevel(classId, lvl uint32) (int, map[string]interface{}) {
	modules, err := models.DBModule.FindModulesForLevel(classId, lvl)
	if err != nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Error fetching the modules for that level.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"modules": modules,
	}
}

func GetStudentsForModule(moduleCode string) (int, map[string]interface{}) {
	students, err := models.DBModule.FindStudentsForModule(moduleCode, "Student")
	if err != nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Error fetching the students for that module.",
		}
	}

	parsedStudents := []map[string]interface{}{}
	for _, student := range students {
		studentMap := map[string]interface{}{
			"id": 				student.ID,
			"first_name":		student.FirstName,
			"last_name":		student.LastName,
			"username":			student.Username,
			"email":			student.Email,
			"matric_number":	student.MatricNumber,
			"matric_date":		student.MatricDate,
			"date_of_birth":	student.DateOfBirth,
			"admin":			student.Admin,
			"avatar_id":		student.AvatarId,
			"avatar":			nil,
		}

		if student.Avatar != nil {
			studentMap["avatar"] = student.Avatar.Url
		}

		parsedStudents = append(parsedStudents, studentMap)
	}

	return http.StatusOK, map[string]interface{}{
		"students": parsedStudents,
	}
}

func AddStudentToModule(moduleCode string, userId uint32) (int, map[string]interface{}) {
	student, err := models.DBUser.FindUserWithId(userId)
	if err != nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Student does not exist",
		}
	}

	dbUserModule, err := models.DBModule.AddStudentToModule(moduleCode, userId)
	if err != nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Error adding student to the module.",
		}
	}

	if dbUserModule == nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Student is already in this module.",
		}
	}

	studentMap := map[string]interface{}{
		"id": 				student.ID,
		"first_name":		student.FirstName,
		"last_name":		student.LastName,
		"username":			student.Username,
		"email":			student.Email,
		"matric_number":	student.MatricNumber,
		"matric_date":		student.MatricDate,
		"date_of_birth":	student.DateOfBirth,
		"admin":			student.Admin,
		"avatar_id":		student.AvatarId,
		"avatar":			nil,
	}

	if student.Avatar != nil {
		studentMap["avatar"] = student.Avatar.Url
	}

	return http.StatusCreated, map[string]interface{}{
		"message": "Student added to the module successfully",
		"module": dbUserModule,
		"student": studentMap,
	}
}
