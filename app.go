package main

import (
	// Go Libs
	"fmt"
	"os"
	"strconv"

	// Third Party Libs
	"github.com/unrolled/secure"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	// My Libs
	"crypto/rsa"
	"github.com/YagoCarballo/kumquat-academy-api/api"
	"github.com/YagoCarballo/kumquat-academy-api/database"
	"github.com/YagoCarballo/kumquat-academy-api/tools"
)

// The path to the Settings file
const SETTINGS_PATH = "./settings.toml"

// Creates the Endpoint routes to start listening to
func loadRoutes(router *web.Mux, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) {
	// Creates the /api routes
	apiRoutes := api.New(privateKey, publicKey)
	apiRoutes.SetupRoutes(router)
}

// Adds some global middlewares to the server
func setupMiddlewares(router *web.Mux) {
	router.Use(middleware.EnvInit)   // Creates an Environment Map if the current one is Nil
	router.Use(middleware.RealIP)    // Finds the Real IP of the Client
	router.Use(middleware.Recoverer) // Recovers from Panics and throws a 500 error

	if tools.GetSettings().Server.Debug {
		router.Use(middleware.Logger) // Enables Logging of routes
		//		router.Use(middleware.NoCache) // Disables Caching
	}

	// Sets the Secure Middleware Options
	secureMiddleware := secure.New(secure.Options{
		IsDevelopment: tools.GetSettings().Server.Debug,
	})

	// Adds the Secure Middleware
	router.Use(secureMiddleware.Handler)
}

func StartServer() {
	err := tools.LoadSettings(SETTINGS_PATH)
	if err != nil {
		panic(err)
	}

	err, _ = database.InitDatabase()
	if err != nil {
		panic(err)
	}

	// Loads the Server Settings
	serverSettings := tools.GetSettings().Server

	// Loads the Keys
	privateKey, publicKey, err := tools.LoadKey(serverSettings.PrivateKey, serverSettings.PublicKey)
	if err != nil {
		panic(err)
	}

	router := web.New()

	// Sets up the Global Middlewares
	setupMiddlewares(router)

	// Loads the Routes
	loadRoutes(router, privateKey, publicKey)

	// Sets the Port found in the Environment variables, (fallback to the settings)
	port := os.Getenv("APIPORT")
	if port == "" {
		port = strconv.Itoa(serverSettings.Port)
	}

	// Prints the port in use
	fmt.Println("Listening on Port: ", port)

	// Starts the Server
	err = graceful.ListenAndServe(":"+port, router)
	if err != nil {
		panic(err)
	}
}

// Starts the Server
func main() {
	StartServer()
}
