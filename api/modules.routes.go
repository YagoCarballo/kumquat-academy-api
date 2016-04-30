package api

import (
	"net/http"
	"encoding/json"

	"github.com/zenazn/goji/web"

	"github.com/YagoCarballo/kumquat-academy-api/tools"
	"github.com/YagoCarballo/kumquat-academy-api/api/middlewares"
	"github.com/YagoCarballo/kumquat-academy-api/api/endpoints"

	"github.com/YagoCarballo/kumquat-academy-api/database/models"

	. "github.com/YagoCarballo/kumquat-academy-api/constants"
)

func (api *API) LoadModulesEndpoints() {

	// Creates the GET -> /user/:username/modules endpoint
	api.routes.Get("/user/:username/modules", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		var username = c.URLParams["username"]

		if cookieData.Username != username {
			api.renderer.JSON(w, http.StatusForbidden, map[string]interface{}{
				"error":   "AccessDenied",
				"message": "Not enough permissions to access this area.",
			})
		} else {
			status, message := endpoints.GetModulesForUser(username)
			api.renderer.JSON(w, status, message)
		}
	}, api.privateKey, api.publicKey))

	// Creates the GET -> /modules/raw endpoint
	api.routes.Get("/modules/raw", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)

		if !cookieData.Admin {
			api.renderer.JSON(w, http.StatusForbidden, map[string]interface{}{
				"error": "AccessDenied",
				"message": "Not enough permissions to access this area.",
			}); return
		}

		params := r.URL.Query()
		query := params.Get("q")
		page, status, _ := tools.ParseID(params.Get("page"))
		if status != http.StatusOK {
			page = 0
		}

		status, message := endpoints.FindRawModules(query, int(page))
		api.renderer.JSON(w, status, message)

	}, api.privateKey, api.publicKey))

	// Creates the GET -> /modules endpoint
	api.routes.Get("/modules", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		status, message := endpoints.GetModulesForUser(cookieData.Username)
		api.renderer.JSON(w, status, message)

	}, api.privateKey, api.publicKey))

	api.routes.Put("/module", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)

		if cookieData.Admin == false {
			api.renderer.JSON(w, http.StatusForbidden, map[string]interface{}{
				"error":   "AccessDenied",
				"message": "Not enough permissions to create a module.",
			})
			return
		}

		decoder := json.NewDecoder(r.Body)
		var module models.Module
		err := decoder.Decode(&module)
		if err != nil {
			api.renderer.JSON(w, http.StatusForbidden, map[string]interface{}{
				"error":   "InvalidData",
				"message": "The data provided is not valid.",
			})
			return
		}

		status, message := endpoints.CreateModule(
			module.Title,
			module.Description,
			module.Icon,
			module.Color,
			module.Duration,
		)

		api.renderer.JSON(w, status, message)

	}, api.privateKey, api.publicKey))

	api.routes.Get("/module/:moduleCode/students", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		var moduleCode = c.URLParams["moduleCode"]

		// Does the user have enough access rights?
		status, err := tools.VerifyAccess(moduleCode, cookieData.UserId, ReadPermission, models.DBPermissions.IsActionPermittedOnModuleWithCode)
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		status, message := endpoints.GetStudentsForModule(moduleCode)
		api.renderer.JSON(w, status, message)

	}, api.privateKey, api.publicKey))

	api.routes.Put("/module/:moduleCode/student/:studentId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
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

		status, message := endpoints.AddStudentToModule(moduleCode, studentId)
		api.renderer.JSON(w, status, message)

	}, api.privateKey, api.publicKey))
}
