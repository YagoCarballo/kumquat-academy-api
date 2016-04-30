package models

import (
	"testing"

	. "github.com/franela/goblin"
	"time"
)

func Test_Database_LectureSlots(t *testing.T) {
	g := Goblin(t)
	var lectureSlotId uint32

	g.Describe("When managing the class lectureSlots", func() {
		g.It("Should get an error when reading a missing lectureSlot", func() {
			lectureSlot, err := DBLectureSlot.ReadLectureSlot(77)

			g.Assert(err == nil).IsTrue()
			g.Assert(lectureSlot == nil).IsTrue()
		})

		g.It("Should fail when creating a lectureSlot for a missing module", func() {
			lectureSlot, err := DBLectureSlot.CreateLectureSlot(
				999,
				"Seminar Room 2",
				"Lab",
				time.Now().AddDate(0, 0, 9),
				time.Now().AddDate(0, 0, 12),
			)

			g.Assert(err != nil).IsTrue()
			g.Assert(lectureSlot == nil).IsTrue()
		})

		g.It("Should be able to create a lectureSlot", func() {
			lectureSlot, err := DBLectureSlot.CreateLectureSlot(
				1,
				"Seminar Room 2",
				"Lecture",
				time.Now().AddDate(0, 0, 9),
				time.Now().AddDate(0, 0, 12),
			)

			g.Assert(err == nil).IsTrue()
			g.Assert(lectureSlot != nil).IsTrue()

			lectureSlotId = lectureSlot.ID
		})

		g.It("Should be able to access a lectureSlot", func() {
			lectureSlot, err := DBLectureSlot.ReadLectureSlot(lectureSlotId)

			g.Assert(err == nil).IsTrue()
			g.Assert(lectureSlot != nil).IsTrue()
		})

		g.It("Should be able to get a list of lectureSlots", func() {
			lectureSlots, err := DBLectureSlot.FindLectureSlotsForModule("AC31007")

			g.Assert(err == nil).IsTrue()
			g.Assert(lectureSlots != nil).IsTrue()
			g.Assert(len(lectureSlots) >= 1).IsTrue()
		})

		g.It("Should be able to update a lectureSlot", func() {
			lectureSlot, err := DBLectureSlot.UpdateLectureSlot(
				lectureSlotId,
				1,
				"Seminar Room 3",
				"Lecture",
				time.Now().AddDate(0, 0, 9),
				time.Now().AddDate(0, 0, 12),
			)

			g.Assert(err == nil).IsTrue()
			g.Assert(lectureSlot != nil).IsTrue()
		})

		g.It("Should be able to remove a lectureSlot", func() {
			count, err := DBLectureSlot.DeleteLectureSlot(lectureSlotId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count != 0).IsTrue()
		})

		g.It("Should fail removing a missing lectureSlot", func() {
			count, err := DBLectureSlot.DeleteLectureSlot(lectureSlotId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 0).IsTrue()
		})

		g.It("Should fail to update a missing lectureSlot", func() {
			lectureSlot, err := DBLectureSlot.UpdateLectureSlot(
				lectureSlotId,
				1,
				"Seminar Room 2",
				"Lecture",
				time.Now().AddDate(0, 0, 9),
				time.Now().AddDate(0, 0, 12),
			)

			g.Assert(err == nil).IsTrue()
			g.Assert(lectureSlot == nil).IsTrue()
		})
	})
}
