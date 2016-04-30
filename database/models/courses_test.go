package models

import (
	"os"
	"testing"

	. "github.com/franela/goblin"
	"github.com/jinzhu/gorm"

	"github.com/YagoCarballo/kumquat.academy.api/tools"
	"github.com/YagoCarballo/kumquat.academy.api/database"
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

func Test_Database_Courses(t *testing.T) {
	g := Goblin(t)

	var courseId uint32 = 77

	g.Describe("When asking the courses a user belongs", func() {
		g.It("should a list of courses for that user", func() {
			courses, err := DBCourse.FindCoursesForUser(3)

			g.Assert(err == nil).IsTrue()
			g.Assert(len(courses) > 0).IsTrue()
			g.Assert(len(courses[0].Modules) > 0).IsTrue()
		})

		g.It("should a list of courses for an admin", func() {
			courses, err := DBCourse.FindCoursesForUser(1)

			g.Assert(err == nil).IsTrue()
			g.Assert(len(courses) > 0).IsTrue()
			g.Assert(len(courses[0].Modules) > 0).IsTrue()
		})

		g.It("should get no courses for a missing user", func() {
			courses, err := DBCourse.FindCoursesForUser(999)

			g.Assert(err == nil).IsTrue()
			g.Assert(len(courses) == 0).IsTrue()
		})

		g.It("Should no modules for a user that has no modules", func() {
			modules, err := DBCourse.FindLevelModulesForCourseAndUser(1, 1, false)

			g.Assert(err == nil).IsTrue()
			g.Assert(modules != nil).IsTrue()
			g.Assert(len(modules) == 0).IsTrue()
		})

		g.It("Should be able to get a list of all the level modules for a user", func() {
			modules, err := DBCourse.FindLevelModulesForCourseAndUser(1, 2, false)

			g.Assert(err == nil).IsTrue()
			g.Assert(modules != nil).IsTrue()
			g.Assert(len(modules) >= 2).IsTrue()
		})

		g.It("Admins should be able to get a list of all the level modules", func() {
			modules, err := DBCourse.FindLevelModulesForCourseAndUser(1, 1, true)

			g.Assert(err == nil).IsTrue()
			g.Assert(modules != nil).IsTrue()
			g.Assert(len(modules) >= 3).IsTrue()
		})
	})

	g.Describe("When managing the courses", func() {
		g.It("Should get an error when reading a missing course", func() {
			course, err := DBCourse.ReadCourse(courseId)

			g.Assert(err == nil).IsTrue()
			g.Assert(course == nil).IsTrue()
		})

		g.It("Should success when creating a course", func() {
			course, err := DBCourse.CreateCourse("test-course", "course created through tests")

			g.Assert(err == nil).IsTrue()
			g.Assert(course != nil).IsTrue()

			courseId = course.ID
		})

		g.It("Should be able to access a course", func() {
			course, err := DBCourse.ReadCourse(courseId)

			g.Assert(err == nil).IsTrue()
			g.Assert(course != nil).IsTrue()
		})

		g.It("Should be able to update a course", func() {
			course, err := DBCourse.UpdateCourse(courseId, "updated-test-course", "a test course")

			g.Assert(err == nil).IsTrue()
			g.Assert(course != nil).IsTrue()
		})


		g.It("Should be able to access a course with it's title", func() {
			course, err := DBCourse.GetCourseWithTitle("updated-test-course")

			g.Assert(err == nil).IsTrue()
			g.Assert(course != nil).IsTrue()
		})

		g.It("Should be able to remove a course", func() {
			count, err := DBCourse.DeleteCourse(courseId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count != 0).IsTrue()
		})

		g.It("Should fail removing a missing course", func() {
			count, err := DBCourse.DeleteCourse(courseId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 0).IsTrue()
		})

		g.It("Should fail to update a missing course", func() {
			course, err := DBCourse.UpdateCourse(courseId, "should-fail", "")

			g.Assert(err == nil).IsTrue()
			g.Assert(course == nil).IsTrue()
		})

		g.It("Should not be able to access a missing course with it's title", func() {
			course, err := DBCourse.GetCourseWithTitle("-null-")

			g.Assert(err == nil).IsTrue()
			g.Assert(course == nil).IsTrue()
		})
	})

	// Remove the Temporary Settings
	os.Remove(SETTINGS_PATH)
}
