package models

import (
	"time"
	"reflect"
	"testing"

	. "github.com/franela/goblin"
)

func Test_Database_Lectures(t *testing.T) {
	g := Goblin(t)
	var lectureId uint32
	lectureSlotId := uint32(1)

	g.Describe("When managing the class lectures", func() {
		g.It("Should get an error when reading a missing lecture", func() {
			lecture, err := DBLecture.ReadLecture(77)

			g.Assert(err == nil).IsTrue()
			g.Assert(lecture == nil).IsTrue()
		})

		g.It("Should fail when creating a lecture for a missing module", func() {
			lecture, err := DBLecture.CreateLecture(
				999,
				"Seminar Room 2",
				"Introduction",
				"Introduction",
				time.Now().AddDate(0, 0, 9),
				time.Now().AddDate(0, 0, 12),
				false,
				nil,
			)

			g.Assert(err != nil).IsTrue()
			g.Assert(lecture == nil).IsTrue()
		})

		g.It("Should be able to create a lecture", func() {
			lecture, err := DBLecture.CreateLecture(
				1,
				"Seminar Room 2",
				"Introduction",
				"Introduction",
				time.Now().AddDate(0, 0, 9),
				time.Now().AddDate(0, 0, 12),
				false,
				&lectureSlotId,
			)

			g.Assert(err == nil).IsTrue()
			g.Assert(lecture != nil).IsTrue()

			lectureId = lecture.ID
		})

		g.It("Should be able to access a lecture", func() {
			lecture, err := DBLecture.ReadLecture(lectureId)

			g.Assert(err == nil).IsTrue()
			g.Assert(lecture != nil).IsTrue()
		})

		g.It("Should be able to get a list of lectures", func() {
			lectures, err := DBLecture.FindLecturesForModule("AC31007")

			g.Assert(err == nil).IsTrue()
			g.Assert(lectures != nil).IsTrue()
			g.Assert(len(lectures) >= 1).IsTrue()
		})

		g.It("Should be able to get a list of lectures in a date-range", func() {
			start := time.Now().AddDate(0, 0, 2)
			end := time.Now().AddDate(0, 0, 20)
			lectures, err := DBLecture.FindLecturesForModuleInRange("AC31007", &start, &end, -1)

			g.Assert(err == nil).IsTrue()
			g.Assert(lectures != nil).IsTrue()
			g.Assert(len(lectures) >= 1).IsTrue()
		})

		g.It("Should be able to get an empty list of lectures in a date out of range", func() {
			start := time.Now().AddDate(0, 0, 900)
			end := time.Now().AddDate(0, 0, 1000)
			lectures, err := DBLecture.FindLecturesForModuleInRange("AC31007", &start, &end, -1)

			g.Assert(err == nil).IsTrue()
			g.Assert(lectures != nil).IsTrue()
			g.Assert(len(lectures) == 0).IsTrue()
		})

		g.It("Should be able to group a list of lectures into weeks", func() {
			lectures, err := DBLecture.FindLecturesForModule("AC31007")

			g.Assert(err == nil).IsTrue()
			g.Assert(lectures != nil).IsTrue()
			g.Assert(len(lectures) >= 1).IsTrue()

			weeks := GroupLecturesInWeeks(&lectures)

			g.Assert(weeks != nil).IsTrue()
			g.Assert(reflect.TypeOf(weeks).String()).Equal("map[string][]models.Lecture")
		})

		g.It("Should be able to update a lecture", func() {
			lecture, err := DBLecture.UpdateLecture(
				lectureId,
				1,
				"Seminar Room 3",
				"Introduction",
				"Introduction",
				time.Now().AddDate(0, 0, 9),
				time.Now().AddDate(0, 0, 12),
				false,
				&lectureSlotId,
			)

			g.Assert(err == nil).IsTrue()
			g.Assert(lecture != nil).IsTrue()
		})

		g.It("Should be able to remove a lecture", func() {
			count, err := DBLecture.DeleteLecture(lectureId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count != 0).IsTrue()
		})

		g.It("Should fail removing a missing lecture", func() {
			count, err := DBLecture.DeleteLecture(lectureId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 0).IsTrue()
		})

		g.It("Should fail to update a missing lecture", func() {
			lecture, err := DBLecture.UpdateLecture(
				lectureId,
				1,
				"Seminar Room 2",
				"Introduction",
				"Introduction",
				time.Now().AddDate(0, 0, 9),
				time.Now().AddDate(0, 0, 12),
				false,
				&lectureSlotId,
			)

			g.Assert(err == nil).IsTrue()
			g.Assert(lecture == nil).IsTrue()
		})
	})

	g.Describe("When accessing the list of lectures", func () {
		g.It("Should be able to find the lectures grouped by weeks", func () {
			weeks, err := DBLecture.FindLecturesWeeksForModule("AC31007")

			g.Assert(err == nil).IsTrue()
			g.Assert(weeks != nil).IsTrue()
			g.Assert(reflect.TypeOf(weeks).String()).Equal("map[string][]models.Lecture")
		})
	})
}
