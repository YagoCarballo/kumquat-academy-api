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

func (api *API) LoadClassesEndpoints() {
	api.routes.Get("/course/:courseId/classOf/:year", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		year := c.URLParams["year"]

		// Parse the course Id
		courseId, status, err := tools.ParseID(c.URLParams["courseId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(courseId, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnCourse)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Get the class
		status, message := endpoints.GetClassForYear(courseId, year)
		if status == http.StatusOK {
			course := message["class"].(*models.Class)

			// Does the user have enough access rights?
			accessStatus, accessError := tools.VerifyAccess(course.ID, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnCourse)
			if accessStatus != http.StatusOK {
				api.renderer.JSON(w, accessStatus, accessError); return
			}
		}

		// Give the response
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Put("/course/:courseId/class", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)

		// Parse the course Id
		courseId, status, err := tools.ParseID(c.URLParams["courseId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(courseId, cookieData.UserId, WritePermission, models.DBPermissions.IsActionPermittedOnCourse)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the JSON Body
		var class models.Class
		status, errMessage := tools.ParseBody(r.Body, &class)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		// Process the action and Give the response
		status, message := endpoints.CreateClass(courseId, class.Title, class.Start, class.End, class.Levels)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Get("/course/:courseId/class/:classId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		// Parse the course Id
		courseId, status, err := tools.ParseID(c.URLParams["courseId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the class Id
		classId, status, err := tools.ParseID(c.URLParams["classId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(courseId, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnCourse)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.GetClass(uint32(classId))
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Post("/course/:courseId/class/:classId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		// Parse the course Id
		courseId, status, err := tools.ParseID(c.URLParams["courseId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the class Id
		classId, status, err := tools.ParseID(c.URLParams["classId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(courseId, cookieData.UserId, WritePermission, models.DBPermissions.IsActionPermittedOnCourse)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the JSON Body
		var class models.Class
		status, errMessage := tools.ParseBody(r.Body, &class)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		// Process the action and Give the response
		status, message := endpoints.UpdateClass(classId, courseId, class.Title, class.Start, class.End)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Delete("/course/:courseId/class/:classId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		// Parse the course Id
		courseId, status, err := tools.ParseID(c.URLParams["courseId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the class Id
		classId, status, err := tools.ParseID(c.URLParams["classId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(courseId, cookieData.UserId, DeletePermission, models.DBPermissions.IsActionPermittedOnCourse)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.DeleteClass(classId)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Get("/course/:courseId/classes", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		// Parse the course Id
		courseId, status, err := tools.ParseID(c.URLParams["courseId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(courseId, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnCourse)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.GetClassesForCourse(courseId)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))
}
