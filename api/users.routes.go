package api

import (
	"time"
	"net/http"

	"github.com/zenazn/goji/web"

	"github.com/YagoCarballo/kumquat-academy-api/tools"
	"github.com/YagoCarballo/kumquat-academy-api/api/middlewares"
	"github.com/YagoCarballo/kumquat-academy-api/api/endpoints"
	"github.com/YagoCarballo/kumquat-academy-api/database/models"

	. "github.com/YagoCarballo/kumquat-academy-api/constants"
)

type (
	StudentForModule struct {
		FirstName string `json:"first_name"`
		LastName string `json:"last_name"`
		Username string `json:"username"`
		Email string `json:"email"`
		MatricNumber string `json:"matric_number"`
		MatricDate time.Time `json:"matric_date"`
		DateOfBirth time.Time `json:"date_of_birth"`
		AvatarId *uint32 `json:"avatar_id"`
	}

	ChangePassword struct {
		Password string `json:"password"`
	}
)

func (api *API) LoadUsersEndpoints() {
	api.routes.Put("/module/:moduleCode/student", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		var moduleCode = c.URLParams["moduleCode"]

		// Does the user have enough access rights?
		status, err := tools.VerifyAccess(moduleCode, cookieData.UserId, WritePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the JSON Body
		var studentData StudentForModule
		status, errMessage := tools.ParseBody(r.Body, &studentData)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		status, message := endpoints.CreateStudentAndAddToModule(
			moduleCode,
			studentData.FirstName,
			studentData.LastName,
			studentData.Username,
			studentData.Email,
			studentData.MatricNumber,
			studentData.MatricDate,
			studentData.DateOfBirth,
			studentData.AvatarId,
		)

		api.renderer.JSON(w, status, message)

	}, api.privateKey, api.publicKey))

	api.routes.Post("/module/:moduleCode/student/:studentId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		var moduleCode = c.URLParams["moduleCode"]

		studentId, status, err := tools.ParseID(c.URLParams["studentId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(moduleCode, cookieData.UserId, WritePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}


		// Parse the JSON Body
		var student StudentForModule
		status, errMessage := tools.ParseBody(r.Body, &student)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		status, message := endpoints.UpdateStudent(
			studentId,
			moduleCode,
			student.FirstName,
			student.LastName,
			student.Username,
			student.Email,
			student.MatricNumber,
			student.MatricDate,
			student.DateOfBirth,
			student.AvatarId,
		)

		api.renderer.JSON(w, status, message)

	}, api.privateKey, api.publicKey))


	api.routes.Get("/module/:moduleCode/students/search/:query", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		var moduleCode = c.URLParams["moduleCode"]
		var query = c.URLParams["query"]

		// Does the user have enough access rights?
		status, err := tools.VerifyAccess(moduleCode, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		status, message := endpoints.SearchUsers(query, moduleCode)
		api.renderer.JSON(w, status, message)

	}, api.privateKey, api.publicKey))

	api.routes.Get("/password/:email/reset", func(c web.C, w http.ResponseWriter, r *http.Request) {
		var email = c.URLParams["email"]

		status, message := endpoints.ForgotPassword(email)
		api.renderer.JSON(w, status, message)

	})

	api.routes.Post("/password/reset/:token", func(c web.C, w http.ResponseWriter, r *http.Request) {
		var token = c.URLParams["token"]

		// Parse the JSON Body
		var changePassword ChangePassword
		status, errMessage := tools.ParseBody(r.Body, &changePassword)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		status, message := endpoints.ChangePassword(token, changePassword.Password)
		api.renderer.JSON(w, status, message)

	})

	api.routes.Put("/user/:userId/avatar", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		file, header, err := r.FormFile("file")

		userId, status, errMsg := tools.ParseID(c.URLParams["userId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMsg); return
		}

		// Does the user have enough access rights?
		if userId != cookieData.UserId && !cookieData.Admin {
			api.renderer.JSON(w, http.StatusForbidden, map[string]interface{}{
				"error":   "AccessDenied",
				"message": "Not enough permissions to change the avatar for this user",
			})
			return
		}

		if err != nil {
			api.renderer.JSON(w, http.StatusConflict, map[string]interface{}{
				"error": "Conflict",
				"message": "Invalid or Missing File",
			}); return
		}

		// Process the action and Give the response
		status, message := endpoints.UploadAvatar(userId, file, header)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))
}
