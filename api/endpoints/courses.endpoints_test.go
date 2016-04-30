package endpoints

import (
	"testing"
	"net/http"

	. "github.com/franela/goblin"

	"github.com/YagoCarballo/kumquat.academy.api/database/models"
)


func Test_Courses(t *testing.T) {
	g := Goblin(t)

	// Disable Verbose Logger
	db.LogMode(false)

	var crudCourse uint32

	g.Describe("Tests CRUD for Courses", func() {
		g.It("Should get a 404 not found course", func() {
			status, data := GetCourse(999)
			g.Assert(status).Equal(http.StatusNotFound)
			g.Assert(data["error"] != nil).IsTrue()
		})

		g.It("Should create a course", func() {
			status, data := CreateCourse(
				"Test Course",
				"This is a course created through unit tests.",
			)
			g.Assert(status).Equal(http.StatusCreated)
			g.Assert(data["error"] == nil).IsTrue()

			course := data["course"].(*models.Course)
			g.Assert(course != nil).IsTrue()
			g.Assert(course.ID != 0).IsTrue()
			crudCourse = course.ID;
		})

		g.It("Should break when creating a course with an existing title", func() {
			status, data := CreateCourse(
				"Test Course",
				"This is a course created through unit tests.",
			)
			g.Assert(status).Equal(http.StatusExpectationFailed)
			g.Assert(data["error"] != nil).IsTrue()
		})

		g.It("Should update a course", func() {
			status, data := UpdateCourse(
				crudCourse,
				"Updated Test Course",
				"This is a course created through unit tests and has been updated.",
			)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(data["error"] == nil).IsTrue()
			g.Assert(data["course"] != nil).IsTrue()
		})

		g.It("Should delete a course", func() {
			status, data := DeleteCourse(crudCourse)
			g.Assert(status).Equal(http.StatusAccepted)
			g.Assert(data["error"] == nil).IsTrue()
			g.Assert(data["message"] != nil).IsTrue()

			status, data = GetCourse(crudCourse)
			g.Assert(status).Equal(http.StatusNotFound)
			g.Assert(data["error"] != nil).IsTrue()
		})
	})

	g.Describe("Tests Get Courses For User endpoint", func() {
		g.It("Should get two courses for user student", func() {
			status, data := GetCoursesForUser(3)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(data != nil).IsTrue()

			courses := data["courses"].([]models.Course)
			g.Assert(len(courses) == 2).IsTrue()
		})

		g.It("Should get all courses for user admin", func() {
			status, data := GetCoursesForUser(1)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(data != nil).IsTrue()

			courses := data["courses"].([]models.Course)
			g.Assert(len(courses) > 0).IsTrue()
		})

		g.It("Should get no courses for user guest", func() {
			status, data := GetCoursesForUser(4)
			courses := data["courses"].([]models.Course)

			g.Assert(status).Equal(http.StatusOK)
			g.Assert(data != nil).IsTrue()
			g.Assert(len(courses) == 0).IsTrue()
		})

		g.It("Should get an error for user missing", func() {
			status, data := GetCoursesForUser(999)
			courses := data["courses"].([]models.Course)

			g.Assert(status).Equal(http.StatusOK)
			g.Assert(data != nil).IsTrue()
			g.Assert(len(courses) == 0).IsTrue()
		})
	})
}
