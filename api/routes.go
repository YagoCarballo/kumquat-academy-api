package api

import (
	"fmt"
	"net/http"
	"crypto/rsa"

	"github.com/albrow/forms"
	"github.com/unrolled/render"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"github.com/YagoCarballo/kumquat-academy-api/api/endpoints"
	"github.com/YagoCarballo/kumquat-academy-api/api/middlewares"
	"github.com/YagoCarballo/kumquat-academy-api/database/models"
	"github.com/YagoCarballo/kumquat-academy-api/tools"
)

type (
	API struct {
		routes       *web.Mux
		secureRoutes *web.Mux
		renderer     *render.Render
		privateKey	 *rsa.PrivateKey
		publicKey	 *rsa.PublicKey
	}
)

// Constructor
func New(privateKey	 *rsa.PrivateKey, publicKey	 *rsa.PublicKey) *API {
	// Creates an Instance of the Routes
	routes := web.New()
	secureRoutes := web.New()
	renderer := render.New()

	// Returns the created struct as pointer
	return &API{routes, secureRoutes, renderer, privateKey, publicKey}
}

func (api *API) SetupRoutes(baseRouter *web.Mux) {
	// Generates an API path with the current API version and the path set in settings
	apiPath := fmt.Sprintf("%s/v%d/*", tools.GetSettings().Api.Prefix, tools.GetSettings().Api.Version)

	// Sets the Route of the API
	baseRouter.Handle(apiPath, api.routes)

	// Routes the endpoints to /api
	api.routes.Use(middleware.SubRouter)

	// Starts listening to the Endpoints
	api.LoadAuthEndpoints()
	api.LoadEndpoints()
	api.LoadModulesEndpoints()
	api.LoadCoursesEndpoints()
	api.LoadClassesEndpoints()
	api.LoadLevelsEndpoints()
	api.LoadAssignmentsEndpoints()
	api.LoadAttachmentsEndpoints()
	api.LoadUsersEndpoints()
	api.LoadLectureEndpoints()
	api.LoadLectureSlotEndpoints()
}

func (api *API) LoadAuthEndpoints() {
	// Auth Login
	// Signs In a user
	api.routes.Post("/auth/sign-in/token", func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Parses the Body
		var data models.SignInUser
		err := data.ParseUser(&r.Body)
		if err != nil {
			api.renderer.JSON(w, http.StatusConflict, map[string]interface{}{
				"error":   "InvalidCredentials",
				"message": "The provided credentials are invalid.",
			})
			return
		}

		// Gets Device ID
		deviceId := r.Header.Get("Device")

		// Processes the Request
		status, message, sessionData := endpoints.SignIn(data.Username, data.Password, deviceId)

		// Attaches the Access token to the response
		message["access_token"] = sessionData.AccessToken

		// Returns the JSON
		api.renderer.JSON(w, status, message)
	})

	// Signs In a user
	api.routes.Post("/auth/sign-in", func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Parses the Body
		var data models.SignInUser
		err := data.ParseUser(&r.Body)
		if err != nil {
			api.renderer.JSON(w, http.StatusConflict, map[string]interface{}{
				"error":   "InvalidCredentials",
				"message": "The provided credentials are invalid.",
			})
			return
		}

		// Gets Device ID
		deviceId := r.Header.Get("Device")

		// Processes the Request
		status, message, sessionData := endpoints.SignIn(data.Username, data.Password, deviceId)

		// Sets the session cookie
		err = tools.SetJWTCookie("token", sessionData, w, api.publicKey)
		if err != nil {
			api.renderer.JSON(w, http.StatusConflict, map[string]interface{}{
				"error":   "InvalidSession",
				"message": "Error when trying to secure the session.",
			})
			return
		}

		// Returns the JSON
		api.renderer.JSON(w, status, message)
	})

	// Auth Register
	// Signs Up a user
	api.routes.Post("/auth/sign-up", func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Parses the Body
		var data models.User

		// Read the request data.
		userData, err := forms.Parse(r)
		if err != nil {
			api.renderer.JSON(w, http.StatusConflict, map[string]interface{}{
				"error":   "InvalidJSON",
				"message": "The data provided is Invalid.",
			})
			return
		}

		// Parse and validate the Data
		err, validationErrors := data.Parse(userData)
		if err != nil {
			api.renderer.JSON(w, http.StatusConflict, map[string]interface{}{
				"error":   "InvalidDate",
				"message": "Error parsing the provided date.",
			})
			return

		} else if validationErrors != nil {
			api.renderer.JSON(w, http.StatusConflict, map[string]interface{}{
				"error":   "ValidationError",
				"message": "The data provided is not valid.",
				"errors":  validationErrors,
			})
			return
		}

		// Processes the Request
		status, message := endpoints.SignUp(&data)

		// Returns the JSON
		api.renderer.JSON(w, status, message)
	})

	// Creates the GET -> /auth/logout endpoint
	api.routes.Get("/auth/logout", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		status, message := endpoints.LogOut(cookieData.AccessToken)

		// Returns the JSON
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))

	// Creates the GET -> /auth/info endpoint
	// This endpoint recovers information about a user linked to a valid session. (Used when rendering on the server)
	api.routes.Get("/auth/info", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		var cookieData *tools.JWTSession = c.Env["token"].(*tools.JWTSession)
		status, message := endpoints.UserInfo(cookieData.UserId)

		// Returns the JSON
		api.renderer.JSON(w, status, message)
	}, api.privateKey, api.publicKey))
}

func (api *API) LoadEndpoints() {
	// Creates the GET -> /api/about endpoint
	api.routes.Get("/hello/:name", func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Processes the Request
		status, message := endpoints.SayHello(c.URLParams["name"])

		// Returns the JSON
		api.renderer.JSON(w, status, message)
	})

	// Exposes some configuration properties used by react.js
	api.routes.Get("/conf", func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Gets the Settings
		settings := tools.GetSettings()

		// Creates the Output
		output := map[string]interface{}{
			"api":         settings.Api,
			"debug":       settings.Server.Debug,
			"port":        settings.Server.Port,
			"title":       settings.Title,
			"description": settings.Description,
		}

		// Returns the JSON
		api.renderer.JSON(w, 200, output)
	})

	api.routes.Get("/private", middlewares.CheckSession(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// Creates the Output
		output := map[string]interface{}{
			"message": "Sucess!!",
		}

		// Returns the JSON
		api.renderer.JSON(w, 200, output)
	}, api.privateKey, api.publicKey))

	api.routes.Get("/restricted", middlewares.Restricted(
		middlewares.Permissions{
			true,  // Read
			false, // Write
		},
		func(c web.C, w http.ResponseWriter, r *http.Request) {
			// Creates the Output
			output := map[string]interface{}{
				"message": "Sucess!!",
			}

			// Returns the JSON
			api.renderer.JSON(w, 200, output)
		}, api.privateKey, api.publicKey))
}
