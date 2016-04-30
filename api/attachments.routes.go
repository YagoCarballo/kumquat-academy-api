package api

import (
	"net/http"

	"github.com/zenazn/goji/web"

	"github.com/YagoCarballo/kumquat.academy.api/api/middlewares"
	"github.com/YagoCarballo/kumquat.academy.api/api/endpoints"
	"github.com/YagoCarballo/kumquat.academy.api/tools"
)

func (api *API) LoadAttachmentsEndpoints() {
	api.routes.Get("/attachment/:name", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		fileType, bytes, err := endpoints.ServeFile(c.URLParams["name"])
		if err != nil || fileType == nil {
			api.renderer.Text(w, http.StatusNotFound, err.Error())
			return
		}

		headers := w.Header()
		headers["Content-Type"] = []string{ *fileType }
		api.renderer.Data(w, http.StatusOK, *bytes)
	}, api.privateKey, api.publicKey))

	api.routes.Get("/attachment/:token/:name", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		fileType, bytes, err := endpoints.ServeFile(c.URLParams["token"])
		if err != nil || fileType == nil {
			api.renderer.Text(w, http.StatusNotFound, err.Error())
			return
		}

		headers := w.Header()
		headers["Content-Type"] = []string{ *fileType }
		api.renderer.Data(w, http.StatusOK, *bytes)
	}, api.privateKey, api.publicKey))

	api.routes.Put("/attachment", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		file, header, err := r.FormFile("file")

		if err != nil {
			api.renderer.JSON(w, http.StatusConflict, map[string]interface{}{
				"error": "Conflict",
				"message": "Invalid or Missing File",
			}); return
		}

		// Does the user have enough access rights?
		if (!cookieData.Admin) {
			api.renderer.JSON(w, http.StatusForbidden, map[string]interface{}{
				"error":   "AccessDenied",
				"message": "Not enough permissions to remove an attachment.",
			})
			return
		}

		// Process the action and Give the response
		status, message := endpoints.UploadFile(file, header)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Put("/attachments", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		err := r.ParseMultipartForm(200000) // grab the multipart form
		if err != nil {
			api.renderer.JSON(w, http.StatusConflict, map[string]interface{}{
				"error": "Conflict",
				"message": "Request is not formatted properly",
			}); return
		}

		// Does the user have enough access rights?
		if (!cookieData.Admin) {
			api.renderer.JSON(w, http.StatusForbidden, map[string]interface{}{
				"error":   "AccessDenied",
				"message": "Not enough permissions to remove an attachment.",
			})
			return
		}

		formData := r.MultipartForm // ok, no problem so far, read the Form data

		//get the *fileHeaders
		fileHeaders := formData.File["files[]"] // grab the filenames

		// Process the action and Give the response
		status, message := endpoints.UploadFiles(fileHeaders)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	api.routes.Delete("/attachment/:attachmentId", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Get and Parse the parameters
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		attachmentId, status, err := tools.ParseID(c.URLParams["attachmentId"])
		if status != http.StatusOK {
			api.renderer.JSON(w, status, err); return
		}

		// Does the user have enough access rights?
		if (!cookieData.Admin) {
			api.renderer.JSON(w, http.StatusForbidden, map[string]interface{}{
				"error":   "AccessDenied",
				"message": "Not enough permissions to remove an attachment.",
			})
			return
		}

		// Process the action and Give the response
		status, message := endpoints.DeleteAttachment(attachmentId)
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))
}
