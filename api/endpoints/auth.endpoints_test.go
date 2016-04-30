package endpoints

import (
	"os"
	"time"
	"testing"
	"net/http"

	. "github.com/franela/goblin"
	"github.com/jinzhu/gorm"

	"github.com/YagoCarballo/kumquat-academy-api/database"
	"github.com/YagoCarballo/kumquat-academy-api/database/models"
	"github.com/YagoCarballo/kumquat-academy-api/tools"
)

// The path to the Settings file
const SETTINGS_PATH = "test-settings.toml"

var (
	db *gorm.DB
)

func init() {
	tools.LoadSettings(SETTINGS_PATH)
	_, db = database.InitDatabase()
}

func Test_Auth(t *testing.T) {
	g := Goblin(t)
	var sessionData *tools.JWTSession
	g.Describe("Tests Sign In endpoint", func() {

		g.It("Should get an access token and some user's data", func() {
			var status int
			var data map[string]interface{}
			status, data, sessionData = SignIn(
				"admin",
				"c7ad44cbad762a5da0a452f9e854fdc1e0e7a52a38015f23f3eab1d80b931dd472634dfac71cd34ebc35d16ab7fb8a90c81f975113d6c7538dc69dd8de9077ec",
				"127.0.0.1",
			)
			g.Assert(status).Equal(http.StatusAccepted)
			g.Assert(data["user"] != nil).IsTrue()
			g.Assert(sessionData != nil).IsTrue()

			var user map[string]interface{} = data["user"].(map[string]interface{});
			g.Assert(user["first_name"] != nil).IsTrue()
			g.Assert(user["last_name"] != nil).IsTrue()
			g.Assert(user["username"] != nil).IsTrue()
			g.Assert(user["email"] != nil).IsTrue()
			g.Assert(user["matric_number"] != nil).IsTrue()
			g.Assert(user["admin"] != nil).IsTrue()
		})

		g.It("Should get an access denied", func() {
			status, data, sessionData := SignIn("admin", "invalid_pass", "127.0.0.1")
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(data["error"] != nil).IsTrue()
			g.Assert(data["message"] != nil).IsTrue()
			g.Assert(data["error"]).Equal("InvalidPassword")
			g.Assert(sessionData != nil).IsTrue()
		})

		g.It("Should get invalid credentials", func() {
			status, data, sessionData := SignIn("", "", "127.0.0.1")
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(data["error"] != nil).IsTrue()
			g.Assert(data["message"] != nil).IsTrue()
			g.Assert(data["error"]).Equal("InvalidCredentials")
			g.Assert(sessionData != nil).IsTrue()
		})
	})

	g.Describe("Tests get User info", func() {
		g.It("Should get the User's Info", func() {
			status, data := UserInfo(sessionData.UserId)

			g.Assert(status).Equal(http.StatusOK)
			g.Assert(data["user"] != nil).IsTrue()
			g.Assert(sessionData != nil).IsTrue()

			var user map[string]interface{} = data["user"].(map[string]interface{});
			g.Assert(user["first_name"] != nil).IsTrue()
			g.Assert(user["last_name"] != nil).IsTrue()
			g.Assert(user["username"] != nil).IsTrue()
			g.Assert(user["email"] != nil).IsTrue()
			g.Assert(user["matric_number"] != nil).IsTrue()
			g.Assert(user["admin"] != nil).IsTrue()

			status, data = UserInfo(0)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(data["error"] != nil).IsTrue()
		})
	})

	g.Describe("Logs out a user", func() {
		g.It("Should get the User's Info", func() {
			session, err := models.DBSession.FindSession(sessionData.AccessToken)
			g.Assert(err == nil).IsTrue()
			g.Assert(session != nil).IsTrue()

			status, data := LogOut(sessionData.AccessToken)
			g.Assert(status).Equal(http.StatusAccepted)
			g.Assert(data["message"] != nil).IsTrue()

			session, err = models.DBSession.FindSession(sessionData.AccessToken)
			g.Assert(err == nil).IsTrue()
			g.Assert(session == nil).IsTrue()

			status, data = LogOut(sessionData.AccessToken)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(data["error"] != nil).IsTrue()
		})
	})

	g.Describe("Tests Sign Up endpoint", func() {
		g.It("Should get an Invalid User", func() {
			status, data := SignUp(nil)
			g.Assert(status).Equal(http.StatusExpectationFailed)
			g.Assert(data["error"] != nil).IsTrue()
			g.Assert(data["message"] != nil).IsTrue()
			g.Assert(data["error"]).Equal("InvalidUser")
		})

		g.It("Should get an Invalid User 2", func() {
			status, data := SignUp(&models.User{})
			g.Assert(status).Equal(http.StatusExpectationFailed)
			g.Assert(data["error"] != nil).IsTrue()
			g.Assert(data["message"] != nil).IsTrue()
			g.Assert(data["error"]).Equal("InvalidUser")
		})

		// Disable Logger
		db.LogMode(false)

		g.It("Should get an Duplicated User", func() {
			status, data := SignUp(&models.User{
				ID:				0,
				FirstName:		"John",
				LastName:		"Doe",
				Username:		"admin",
				Email:			"john@doe.com",
				Password:		string([]byte{}),
				Admin:			false,
				Active:			true,
				DateOfBirth:	time.Now(),
				MatricNumber:	"123456789",
				MatricDate:		time.Now(),

				Sessions:		[]models.Session{},
				CreatedAt:		time.Now(),
				UpdatedAt:		time.Now(),
			})

			g.Assert(status).Equal(http.StatusConflict)
			g.Assert(data["error"] != nil).IsTrue()
			g.Assert(data["message"] != nil).IsTrue()
			g.Assert(data["error"]).Equal("DuplicatedUser")
		})
	})

	// Enable Logger
	db.LogMode(false)

	// Remove the Temporary Settings
	os.Remove(SETTINGS_PATH)
}
