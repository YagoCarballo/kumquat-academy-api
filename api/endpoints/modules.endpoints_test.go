package endpoints

import (
	"testing"
	"net/http"

	. "github.com/franela/goblin"

	"github.com/YagoCarballo/kumquat-academy-api/database/models"
)


func Test_Modules(t *testing.T) {
	g := Goblin(t)

	// Disable Verbose Logger
	db.LogMode(false)

	g.Describe("Tests Get Modules For User endpoint", func() {
		g.It("Should get one module for user student", func() {
			status, data := GetModulesForUser("student")
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(data != nil).IsTrue()

			modules := data["modules"].([]models.OutputModule)
			g.Assert(len(modules) > 0).IsTrue()
			g.Assert(modules[0].Id == 1).IsTrue()
			g.Assert(modules[0].CourseId == 1).IsTrue()
			g.Assert(modules[0].Code == "AC31007").IsTrue()
			g.Assert(modules[0].Color == "#9C0098").IsTrue()
			g.Assert(modules[0].Title == "Big Data").IsTrue()
			g.Assert(modules[0].Description == "Introduction to the world of Big Data").IsTrue()
			g.Assert(modules[0].Icon == "fa-cloud").IsTrue()
			g.Assert(modules[0].Role.Id == 3).IsTrue()
			g.Assert(modules[0].Role.Name == "Student").IsTrue()
			g.Assert(modules[0].Role.Description == "Student of a module / course.").IsTrue()
			g.Assert(modules[0].Role.Read).IsTrue()
			g.Assert(modules[0].Role.Write).IsFalse()
			g.Assert(modules[0].Role.Delete).IsFalse()
			g.Assert(modules[0].Role.Update).IsFalse()
			g.Assert(modules[0].Role.Admin).IsFalse()
			g.Assert(modules[0].Year != "").IsTrue()
		})

		g.It("Should get one module for user admin", func() {
			status, data := GetModulesForUser("admin")
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(data != nil).IsTrue()

			modules := data["modules"].([]models.OutputModule)
			g.Assert(len(modules) > 0).IsTrue()
			g.Assert(modules[0].Role.Id == 0).IsTrue()
			g.Assert(modules[0].Role.Name == "Admin").IsTrue()
			g.Assert(modules[0].Role.Description == "Admin of a module / course.").IsTrue()
			g.Assert(modules[0].Role.Read).IsFalse()
			g.Assert(modules[0].Role.Write).IsFalse()
			g.Assert(modules[0].Role.Delete).IsFalse()
			g.Assert(modules[0].Role.Update).IsFalse()
			g.Assert(modules[0].Role.Admin).IsTrue()
		})

		g.It("Should get no modules for user guest", func() {
			status, data := GetModulesForUser("guest")
			modules := data["modules"].([]models.OutputModule)

			g.Assert(status).Equal(http.StatusOK)
			g.Assert(data != nil).IsTrue()
			g.Assert(len(modules) == 0).IsTrue()
		})

		g.It("Should get an error for user missing", func() {
			status, data := GetModulesForUser("missing")
			modules := data["modules"].([]models.OutputModule)

			g.Assert(status).Equal(http.StatusOK)
			g.Assert(data != nil).IsTrue()
			g.Assert(len(modules) == 0).IsTrue()
		})
	})
}
