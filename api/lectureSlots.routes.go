package api

import (
	"net/http"

	"github.com/zenazn/goji/web"

	"github.com/YagoCarballo/kumquat-academy-api/tools"
	"github.com/YagoCarballo/kumquat-academy-api/api/middlewares"
	"github.com/YagoCarballo/kumquat-academy-api/api/endpoints"

	. "github.com/YagoCarballo/kumquat-academy-api/constants"
	"github.com/YagoCarballo/kumquat-academy-api/database/models"
)

func (api *API) LoadLectureSlotEndpoints() {
	api.routes.Put("/module/:moduleCode/lecture-slot", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		// Does the user have enough access rights?
		status, err := tools.VerifyAccess(moduleCode, cookieData.UserId, WritePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the JSON Body
		var lectureSlot models.LectureSlot
		status, errMessage := tools.ParseBody(r.Body, &lectureSlot)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		// Process the action and Give the response
		status, message := endpoints.CreateLectureSlot(
			moduleCode,
			lectureSlot.Location,
			lectureSlot.Type,
			lectureSlot.Start,
			lectureSlot.End,
		)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Get("/module/:moduleCode/lecture-slot/:lectureSlotId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		lectureSlotId, status, err := tools.ParseID(c.URLParams["lectureSlotId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(moduleCode, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.GetLectureSlot(uint32(lectureSlotId))
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Post("/module/:moduleCode/lecture-slot/:lectureSlotId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		lectureSlotId, status, err := tools.ParseID(c.URLParams["lectureSlotId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(moduleCode, cookieData.UserId, UpdatePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the JSON Body
		var lectureSlot models.LectureSlot
		status, errMessage := tools.ParseBody(r.Body, &lectureSlot)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		// Process the action and Give the response
		status, message := endpoints.UpdateLectureSlot(
			lectureSlotId,
			moduleCode,
			lectureSlot.Location,
			lectureSlot.Type,
			lectureSlot.Start,
			lectureSlot.End,
		)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Delete("/module/:moduleCode/lecture-slot/:lectureSlotId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		lectureSlotId, status, err := tools.ParseID(c.URLParams["lectureSlotId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(moduleCode, cookieData.UserId, DeletePermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.DeleteLectureSlot(uint32(lectureSlotId))
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Get("/module/:moduleCode/lecture-slots", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		moduleCode := c.URLParams["moduleCode"]

		// Does the user have enough access rights?
		status, err := tools.VerifyAccess(moduleCode, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.FindLectureSlotsForModule(moduleCode)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))
}
