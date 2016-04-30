package tools

import (
	"testing"

	. "github.com/franela/goblin"

	"os"
	"reflect"
)

// The path to the Settings file
const SETTINGS_PATH = "test-settings.toml"

func Test_LoadSettings(t *testing.T) {
	g := Goblin(t)

	g.Describe("When loading the settings File, ", func() {

		g.It("loading an Invalid Path does not return an error and creates the file.", func() {
			err := LoadSettings("Invalid Path")
			g.Assert(err == nil).IsTrue()

			// Opens the Database File
			_, err = os.Open("Invalid Path")
			g.Assert(err).Equal(nil)

			// Remove the invalid file
			os.Remove("Invalid Path")
		})
	})

	g.Describe("Settings File", func() {
		// Loads the Settings
		LoadSettings(SETTINGS_PATH)
		settings := GetSettings()

		g.It("Settings should have a valid properties", func() {
			g.Assert(reflect.TypeOf(settings.Title).String()).Equal("string")
			g.Assert(reflect.TypeOf(settings.Description).String()).Equal("string")
		})

		g.It("Should have a valid Server Object", func() {
			server := settings.Server
			g.Assert(reflect.TypeOf(server.Port).String()).Equal("int")
			g.Assert(reflect.TypeOf(server.Debug).String()).Equal("bool")
			g.Assert(reflect.TypeOf(server.Production).String()).Equal("bool")
			g.Assert(reflect.TypeOf(server.PrivateKey).String()).Equal("string")
			g.Assert(reflect.TypeOf(server.PublicKey).String()).Equal("string")
			g.Assert(reflect.TypeOf(server.UploadsPath).String()).Equal("string")
		})

		g.It("Should have a valid Database Object", func() {
			database := settings.Database
			g.Assert(reflect.TypeOf(database.Type).String()).Equal("string")
			g.Assert(reflect.TypeOf(database.Mysql.Username).String()).Equal("string")
			g.Assert(reflect.TypeOf(database.Mysql.Password).String()).Equal("string")
			g.Assert(reflect.TypeOf(database.Mysql.Host).String()).Equal("string")
			g.Assert(reflect.TypeOf(database.Mysql.Name).String()).Equal("string")
			g.Assert(reflect.TypeOf(database.Sqlite.Path).String()).Equal("string")
		})

		g.It("Should have a valid API Object", func() {
			api := settings.Api
			g.Assert(reflect.TypeOf(api.Prefix).String()).Equal("string")
			g.Assert(reflect.TypeOf(api.Version).String()).Equal("int")
		})

		g.It("Should have a valid Email Object", func() {
			email := settings.Email
			g.Assert(reflect.TypeOf(email.User).String()).Equal("string")
			g.Assert(reflect.TypeOf(email.Password).String()).Equal("string")
			g.Assert(reflect.TypeOf(email.Server).String()).Equal("string")
			g.Assert(reflect.TypeOf(email.Port).String()).Equal("int")
			g.Assert(reflect.TypeOf(email.Sender).String()).Equal("string")
		})
	})

	g.Describe("Environment Variable Settings - Success", func() {
		os.Setenv("DB_USER", "johndoe")
		os.Setenv("DB_PASS", "password")
		os.Setenv("DB_HOST", "example.com:3306")
		os.Setenv("EMAIL_USER", "johndoe")
		os.Setenv("EMAIL_PASS", "password")
		os.Setenv("EMAIL_SMTP_SERVER", "smtp.example.com")
		os.Setenv("EMAIL_SMTP_PORT", "587")
		os.Setenv("EMAIL_SENDER", "Test <test@gmail.com>")
		os.Setenv("UPLOADS_PATH", "./somewhere")
		os.Setenv("APIPORT", "3333")
		os.Setenv("PRODUCTION", "true")

		LoadSettings(SETTINGS_PATH)

		// Loads the Settings
		settings := GetSettings()

		g.It("Settings match Environment Variables", func() {
			g.Assert(settings.Server.Port).Equal(3333)
			g.Assert(settings.Server.Production).IsTrue()
			g.Assert(settings.Database.Type).Equal("MySQL")
			g.Assert(settings.Database.Mysql.Username).Equal("johndoe")
			g.Assert(settings.Database.Mysql.Password).Equal("password")
			g.Assert(settings.Database.Mysql.Host).Equal("example.com:3306")
			g.Assert(settings.Database.Sqlite.Path).Equal("./default.sqlite")
			g.Assert(settings.Email.User).Equal("johndoe")
			g.Assert(settings.Email.Password).Equal("password")
			g.Assert(settings.Email.Server).Equal("smtp.example.com")
			g.Assert(settings.Email.Port).Equal(587)
			g.Assert(settings.Email.Sender).Equal("Test <test@gmail.com>")
			g.Assert(settings.Server.UploadsPath).Equal("./somewhere")
		})
	})

	g.Describe("Environment Variable Settings - Fail", func() {
		os.Setenv("APIPORT", "NaN")

		LoadSettings(SETTINGS_PATH)

		// Loads the Settings
		settings := GetSettings()

		g.It("Settings should fall back the server port to 3000", func() {
			g.Assert(settings.Server.Port).Equal(3000)
		})
	})

	// Remove the Temporary Settings
	os.Remove(SETTINGS_PATH)
}
