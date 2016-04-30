package api

import (
	"net/http"

	"github.com/zenazn/goji/web"

	"github.com/YagoCarballo/kumquat.academy.api/tools"
	"github.com/YagoCarballo/kumquat.academy.api/api/middlewares"
	"github.com/YagoCarballo/kumquat.academy.api/api/endpoints"

	. "github.com/YagoCarballo/kumquat.academy.api/constants"
	"github.com/YagoCarballo/kumquat.academy.api/database/models"
)

func (api *API) LoadLectureEndpoints() {
	api.routes.Put("/module/:moduleCode/lecture", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		// Does the user have enough access rights?
		status, err := tools.VerifyAccess(moduleCode, cookieData.UserId, WritePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the JSON Body
		var lecture models.Lecture
		status, errMessage := tools.ParseBody(r.Body, &lecture)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		// Process the action and Give the response
		status, message := endpoints.CreateLecture(
			moduleCode,
			lecture.Location,
			lecture.Topic,
			lecture.Description,
			lecture.Start,
			lecture.End,
			lecture.Canceled,
			lecture.LectureSlotID,
		)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Get("/module/:moduleCode/lecture/:lectureId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		lectureId, status, err := tools.ParseID(c.URLParams["lectureId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(moduleCode, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.GetLecture(uint32(lectureId))
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Post("/module/:moduleCode/lecture/:lectureId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		lectureId, status, err := tools.ParseID(c.URLParams["lectureId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(moduleCode, cookieData.UserId, UpdatePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the JSON Body
		var lecture models.Lecture
		status, errMessage := tools.ParseBody(r.Body, &lecture)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		// Process the action and Give the response
		status, message := endpoints.UpdateLecture(
			lectureId,
			moduleCode,
			lecture.Location,
			lecture.Topic,
			lecture.Description,
			lecture.Start,
			lecture.End,
			lecture.Canceled,
			lecture.LectureSlotID,
		)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Delete("/module/:moduleCode/lecture/:lectureId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		lectureId, status, err := tools.ParseID(c.URLParams["lectureId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(moduleCode, cookieData.UserId, DeletePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.DeleteLecture(uint32(lectureId))
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Get("/module/:moduleCode/lectures", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		// Does the user have enough access rights?
		status, err := tools.VerifyAccess(moduleCode, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.FindLecturesForModule(moduleCode)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Get("/module/:moduleCode/lecture-weeks", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		// Does the user have enough access rights?
		status, err := tools.VerifyAccess(moduleCode, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.FindLectureWeeksForModule(moduleCode)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))


	api.routes.Get("/module/:moduleCode/lectures-overview", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		// Does the user have enough access rights?
		status, err := tools.VerifyAccess(moduleCode, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.FindLectureWeeksAndSlotsForModule(moduleCode)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Put("/module/:moduleCode/lecture/:lectureId/attachment", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		file, header, err := r.FormFile("file")
		moduleCode := c.URLParams["moduleCode"]

		lectureId, status, errMsg := tools.ParseID(c.URLParams["lectureId"])
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
		status, message := endpoints.UploadLectureAttachments(lectureId, file, header)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Delete("/module/:moduleCode/lecture/:lectureId/attachment/:attachmentId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]
		lectureId, status, err := tools.ParseID(c.URLParams["lectureId"])
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
		status, message := endpoints.RemoveLectureAttachments(lectureId, attachmentId)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Get("/schedule", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)

		// Process the action and Give the response
		status, message := endpoints.FindLectureWeeksForUser(cookieData.UserId)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))
}
