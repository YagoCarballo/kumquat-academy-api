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

func (api *API) LoadCoursesEndpoints() {
	// Creates the GET -> /course/titled/:title endpoint
	api.routes.Get("/course/titled/:title", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		title := c.URLParams["title"]

		// Get the course
		status, message := endpoints.GetCourseWithTitle(title)
		if status == http.StatusOK {
			course := message["course"].(*models.Course)

			// Does the user have enough access rights?
			accessStatus, accessError := tools.VerifyAccess(course.ID, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnCourse)
			if accessStatus != http.StatusOK {
				api.renderer.JSON(w, accessStatus, accessError); return
			}
		}

		// Give the response
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	// Creates the PUT -> /course endpoint
	api.routes.Put("/course", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)

		// Does the user have enough access rights?
		if cookieData.Admin == false {
			api.renderer.JSON(w, http.StatusForbidden, map[string]interface{}{
				"error":   "AccessDenied",
				"message": "Not enough permissions to create a course.",
			})
			return
		}

		// Parse the JSON Body
		var course models.Course
		status, errMessage := tools.ParseBody(r.Body, &course)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		// Process the action and Give the response
		status, message := endpoints.CreateCourse(course.Title, course.Description)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	// Creates the GET -> /course/:id endpoint
	api.routes.Get("/course/:id", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		courseId, status, err := tools.ParseID(c.URLParams["id"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(courseId, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnCourse)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.GetCourse(uint32(courseId))
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	// Creates the POST -> /course endpoint
	api.routes.Post("/course/:id", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		courseId, status, err := tools.ParseID(c.URLParams["id"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(courseId, cookieData.UserId, UpdatePermission, models.DBPermissions.IsActionPermittedOnCourse)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Parse the JSON Body
		var course models.Course
		status, errMessage := tools.ParseBody(r.Body, &course)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, errMessage); return
		}

		// Process the action and Give the response
		status, message := endpoints.UpdateCourse(courseId, course.Title, course.Description)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	// Creates the POST -> /course endpoint
	api.routes.Delete("/course/:id", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		courseId, status, err := tools.ParseID(c.URLParams["id"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		status, err = tools.VerifyAccess(courseId, cookieData.UserId, DeletePermission, models.DBPermissions.IsActionPermittedOnCourse)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Process the action and Give the response
		status, message := endpoints.DeleteCourse(courseId)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	// Creates the GET -> /courses endpoint
	api.routes.Get("/courses", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		status, message := endpoints.GetCoursesForUser(cookieData.UserId)
		api.renderer.JSON(w, status, message)

	}, api.privateKey, api.publicKey))
}
