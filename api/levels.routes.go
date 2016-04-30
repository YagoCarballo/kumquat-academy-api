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

func (api *API) LoadLevelsEndpoints() {
	api.routes.Put("/course/:courseId/class/:classId/level", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		courseId, status, err := tools.ParseID(c.URLParams["courseId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		classId, status, err := tools.ParseID(c.URLParams["classId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		if cookieData.Admin == false {
			api.renderer.JSON(w, http.StatusForbidden, map[string]interface{}{
				"error":   "AccessDenied",
				"message": "Not enough permissions to create a course.",
			})
			return
		}

		// Parse the JSON Body
		var level models.CourseLevel
		status, errMessage := tools.ParseBody(r.Body, &level)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		// Process the action and Give the response
		status, message := endpoints.CreateLevel(level.Level, classId, courseId, level.Start, level.End)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	// Creates the GET -> /course/:id endpoint
	api.routes.Get("/course/:courseId/class/:classId/level/:level", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)

		courseId, status, err := tools.ParseID(c.URLParams["courseId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		classId, status, err := tools.ParseID(c.URLParams["classId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		level, status, err := tools.ParseID(c.URLParams["level"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(courseId, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnCourse)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.GetLevel(uint32(courseId), uint32(classId), uint32(level))
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	// Creates the POST -> /course endpoint
	api.routes.Post("/course/:courseId/class/:classId/level/:level", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		courseId, status, err := tools.ParseID(c.URLParams["courseId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		classId, status, err := tools.ParseID(c.URLParams["classId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		lvl, status, err := tools.ParseID(c.URLParams["level"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(courseId, cookieData.UserId, UpdatePermission, models.DBPermissions.IsActionPermittedOnCourse)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the JSON Body
		var level models.CourseLevel
		status, errMessage := tools.ParseBody(r.Body, &level)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		// Process the action and Give the response
		status, message := endpoints.UpdateLevel(lvl, classId, courseId, level.Start, level.End)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	// Creates the POST -> /course endpoint
	api.routes.Delete("/course/:courseId/class/:classId/level/:level", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		courseId, status, err := tools.ParseID(c.URLParams["courseId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		classId, status, err := tools.ParseID(c.URLParams["classId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		lvl, status, err := tools.ParseID(c.URLParams["level"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(courseId, cookieData.UserId, DeletePermission, models.DBPermissions.IsActionPermittedOnCourse)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.DeleteLevel(courseId, classId, lvl)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Put("/course/:courseId/class/:classId/level/:level/module", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		classId, status, err := tools.ParseID(c.URLParams["classId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		lvl, status, err := tools.ParseID(c.URLParams["level"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		if cookieData.Admin == false {
			api.renderer.JSON(w, http.StatusForbidden, map[string]interface{}{
				"error":   "AccessDenied",
				"message": "Not enough permissions to create a course.",
			})
			return
		}

		// Parse the JSON Body
		var levelModule models.LevelModule
		status, errMessage := tools.ParseBody(r.Body, &levelModule)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		// Process the action and Give the response
		status, message := endpoints.AddModuleToLevel(levelModule.Code, lvl, classId, levelModule.ModuleID, levelModule.Start)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Get("/course/:courseId/class/:classId/level/:level/modules", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)

		courseId, status, err := tools.ParseID(c.URLParams["courseId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		classId, status, err := tools.ParseID(c.URLParams["classId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		level, status, err := tools.ParseID(c.URLParams["level"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(courseId, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnCourse)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.GetModulesForLevel(uint32(classId), uint32(level))
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))
}
