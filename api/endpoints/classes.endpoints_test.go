package endpoints

import (
	"testing"
	"net/http"

	. "github.com/franela/goblin"

	"github.com/YagoCarballo/kumquat.academy.api/database/models"
	"time"
)


func Test_Classes(t *testing.T) {
	g := Goblin(t)

	// Disable Verbose Logger
	db.LogMode(false)

	var crudClass uint32

	g.Describe("Tests CRUD for Classes", func() {
		g.It("Should get a 404 not found class", func() {
			status, data := GetClass(999)
			g.Assert(status).Equal(http.StatusNotFound)
			g.Assert(data["error"] != nil).IsTrue()

			status, data = GetClassForYear(1, "-missing-")
			g.Assert(status).Equal(http.StatusNotFound)
			g.Assert(data["error"] != nil).IsTrue()
		})

		g.It("Should create a class", func() {
			status, data := CreateClass(
				1,
				"1990/1991 <- Test",
				time.Now(),
				time.Now().AddDate(1, 0, 0),
				[]*models.CourseLevel{},
			)
			g.Assert(status).Equal(http.StatusCreated)
			g.Assert(data["error"] == nil).IsTrue()

			class := data["class"].(*models.Class)
			g.Assert(class != nil).IsTrue()
			g.Assert(class.ID != 0).IsTrue()
			g.Assert(class.Title == "1990/1991 <- Test").IsTrue()
			crudClass = class.ID;
		})

		g.It("Should break when creating a class with an existing year", func() {
			status, data := CreateClass(
				1,
				"1990/1991 <- Test",
				time.Now(),
				time.Now().AddDate(1, 0, 0),
				[]*models.CourseLevel{},
			)
			g.Assert(status).Equal(http.StatusExpectationFailed)
			g.Assert(data["error"] != nil).IsTrue()
		})

		g.It("Should update a class", func() {
			status, data := UpdateClass(
				crudClass,
				1,
				"1990/1991 <- Test -> updated",
				time.Now(),
				time.Now().AddDate(1, 0, 0),
			)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(data["error"] == nil).IsTrue()
			g.Assert(data["class"] != nil).IsTrue()

			class := data["class"].(*models.Class)
			g.Assert(class.ID == crudClass).IsTrue()
			g.Assert(class.Title == "1990/1991 <- Test -> updated").IsTrue()
		})

		g.It("Should delete a class", func() {
			status, data := DeleteClass(crudClass)
			g.Assert(status).Equal(http.StatusAccepted)
			g.Assert(data["error"] == nil).IsTrue()
			g.Assert(data["message"] != nil).IsTrue()

			status, data = GetClass(crudClass)
			g.Assert(status).Equal(http.StatusNotFound)
			g.Assert(data["error"] != nil).IsTrue()
		})
	})
}
