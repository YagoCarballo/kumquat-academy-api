package models

import (
	"testing"

	. "github.com/franela/goblin"
	"time"
)

func Test_Database_Classes(t *testing.T) {
	g := Goblin(t)

	var classId uint32 = 77

	g.Describe("When managing the classes", func() {
		g.It("Should get an error when reading a missing class", func() {
			class, err := DBClass.ReadClass(classId)

			g.Assert(err == nil).IsTrue()
			g.Assert(class == nil).IsTrue()
		})

		g.It("Should not be able to access a class with a wrong title", func() {
			class, err := DBClass.GetClassWithTitle(1, "-null-")

			g.Assert(err == nil).IsTrue()
			g.Assert(class == nil).IsTrue()
		})

		g.It("Should success when creating a class", func() {
			class, err := DBClass.CreateClass(
				1,
				"-test-class-",
				time.Now(),
				time.Now(),
				[]*CourseLevel{},
			)

			g.Assert(err == nil).IsTrue()
			g.Assert(class != nil).IsTrue()

			classId = class.ID
		})

		g.It("Should be able to access a class", func() {
			class, err := DBClass.ReadClass(classId)

			g.Assert(err == nil).IsTrue()
			g.Assert(class != nil).IsTrue()
		})

		g.It("Should be able to access a class with it's title", func() {
			class, err := DBClass.GetClassWithTitle(1, "-test-class-")

			g.Assert(err == nil).IsTrue()
			g.Assert(class != nil).IsTrue()
		})

		g.It("Should be able to update a class", func() {
			class, err := DBClass.UpdateClass(
				classId,
				1,
				"-test-class-updated-",
				time.Now(),
				time.Now(),
			)

			g.Assert(err == nil).IsTrue()
			g.Assert(class != nil).IsTrue()
		})


		g.It("Should be able to access a class with it's title", func() {
			class, err := DBClass.ReadClass(classId)

			g.Assert(err == nil).IsTrue()
			g.Assert(class != nil).IsTrue()
		})

		g.It("Should be able to remove a class", func() {
			count, err := DBClass.DeleteClass(classId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count != 0).IsTrue()
		})

		g.It("Should fail removing a missing class", func() {
			count, err := DBClass.DeleteClass(classId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 0).IsTrue()
		})

		g.It("Should fail to update a missing class", func() {
			class, err := DBClass.UpdateClass(
				classId,
				1,
				"-test-class-",
				time.Now(),
				time.Now(),
			)

			g.Assert(err == nil).IsTrue()
			g.Assert(class == nil).IsTrue()
		})

		g.It("Should be able to access all the classes for a course", func() {
			classes, err := DBClass.GetClassesForCourse(1)

			g.Assert(err == nil).IsTrue()
			g.Assert(classes != nil).IsTrue()
			g.Assert(len(classes) > 0).IsTrue()
		})

		g.It("Should not be able to access all the classes for a missing course", func() {
			classes, err := DBClass.GetClassesForCourse(999)

			g.Assert(err == nil).IsTrue()
			g.Assert(classes != nil).IsTrue()
			g.Assert(len(classes) == 0).IsTrue()
		})
	})
}
