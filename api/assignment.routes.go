package api

import (
	"net/http"

	"github.com/zenazn/goji/web"

	"github.com/YagoCarballo/kumquat-academy-api/tools"
	"github.com/YagoCarballo/kumquat-academy-api/api/middlewares"
	"github.com/YagoCarballo/kumquat-academy-api/api/endpoints"

	. "github.com/YagoCarballo/kumquat-academy-api/constants"
	"github.com/YagoCarballo/kumquat-academy-api/database/models"
	"fmt"
)

func (api *API) LoadAssignmentsEndpoints() {
	api.routes.Put("/module/:moduleCode/assignment", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		// Does the user have enough access rights?
		status, err := tools.VerifyAccess(moduleCode, cookieData.UserId, WritePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the JSON Body
		var assignment models.Assignment
		status, errMessage := tools.ParseBody(r.Body, &assignment)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		// Process the action and Give the response
		status, message := endpoints.CreateAssignment(
			assignment.Title,
			assignment.Description,
			assignment.Status,
			assignment.Weight,
			assignment.Start,
			assignment.End,
			assignment.ModuleCode,
		)

		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Get("/module/:moduleCode/assignment/:assignmentId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		assignmentId, status, err := tools.ParseID(c.URLParams["assignmentId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(moduleCode, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.GetAssignment(uint32(assignmentId))
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Post("/module/:moduleCode/assignment/:assignmentId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]
		assignmentId, status, err := tools.ParseID(c.URLParams["assignmentId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(moduleCode, cookieData.UserId, UpdatePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the JSON Body
		var assignment models.Assignment
		status, errMessage := tools.ParseBody(r.Body, &assignment)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		// Process the action and Give the response
		status, message := endpoints.UpdateAssignment(
			assignmentId,
			assignment.Title,
			assignment.Description,
			assignment.Status,
			assignment.Weight,
			assignment.Start,
			assignment.End,
			assignment.ModuleCode,
		)

		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Delete("/module/:moduleCode/assignment/:assignmentId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]
		assignmentId, status, err := tools.ParseID(c.URLParams["assignmentId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(moduleCode, cookieData.UserId, DeletePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.DeleteAssignment(assignmentId)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Get("/module/:moduleCode/assignments", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		// Does the user have enough access rights?
		status, err := tools.VerifyAccess(moduleCode, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		canWrite := models.DBPermissions.IsActionPermittedOnModuleWithCode(cookieData.UserId, moduleCode, WritePermission)

		// Process the action and Give the response
		status, message := endpoints.FindAssignmentsForModule(cookieData.Username, moduleCode, !canWrite)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Put("/module/:moduleCode/assignment/:assignmentId/attachment", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		file, header, err := r.FormFile("file")
		moduleCode := c.URLParams["moduleCode"]

		assignmentId, status, errMsg := tools.ParseID(c.URLParams["assignmentId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMsg); return
		}

		// Does the user have enough access rights?
		status, errMsg = tools.VerifyAccess(moduleCode, cookieData.UserId, WritePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMsg); return
		}

		if err != nil {
			api.renderer.JSON(w, http.StatusConflict, map[string]interface{}{
				"error": "Conflict",
				"message": "Invalid or Missing File",
			}); return
		}

		// Process the action and Give the response
		status, message := endpoints.UploadAssignmentAttachments(assignmentId, file, header)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Delete("/module/:moduleCode/assignment/:assignmentId/attachment/:attachmentId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]
		assignmentId, status, err := tools.ParseID(c.URLParams["assignmentId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		attachmentId, status, err := tools.ParseID(c.URLParams["attachmentId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(moduleCode, cookieData.UserId, DeletePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.RemoveAssignmentAttachments(assignmentId, attachmentId)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Put("/module/:moduleCode/assignment/:assignmentId/submit", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		assignmentId, status, errMsg := tools.ParseID(c.URLParams["assignmentId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMsg); return
		}

		err := r.ParseMultipartForm(200000) // grab the multipart form
		if err != nil {
			fmt.Println(err)
			api.renderer.JSON(w, http.StatusConflict, map[string]interface{}{
				"error": "Conflict",
				"message": "Request is not formatted properly",
			}); return
		}

		// Does the user have enough access rights?
		if (!models.DBPermissions.CanUserSubmitAssignment(cookieData.Username, assignmentId)) {
			api.renderer.JSON(w, http.StatusForbidden, map[string]interface{}{
				"error":   "AccessDenied",
				"message": "Not enough permissions to submit this assignment.",
			})
			return
		}

		formData := r.MultipartForm // ok, no problem so far, read the Form data

		//get the *fileHeaders
		fileHeaders := formData.File["files[]"] // grab the filenames

		// Process the action and Give the response
		status, message := endpoints.SubmitAssignment(cookieData.Username, formData.Value["description"][0], assignmentId, fileHeaders)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Post("/module/:moduleCode/assignment/:assignmentId/grade", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		// Does the user have enough access rights?
		status, errMsg := tools.VerifyAccess(moduleCode, cookieData.UserId, WritePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMsg); return
		}

		// Parse the JSON Body
		var submission map[string]interface{}
		status, errMessage := tools.ParseBody(r.Body, &submission)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		submissionId, status, errMsg := tools.ParseID(fmt.Sprint("", submission["id"].(float64)))
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMsg); return
		}

		grade, status, errMsg := tools.ParseID(fmt.Sprint("", submission["grade"].(float64)))
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMsg); return
		}

		// Process the action and Give the response
		status, message := endpoints.GradeAssignment(submissionId, grade)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))
}
