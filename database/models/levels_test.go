package models

import (
	"testing"

	. "github.com/franela/goblin"
	"time"
)

func Test_Database_Levels(t *testing.T) {
	g := Goblin(t)

	g.Describe("When managing the class levels", func() {
		g.It("Should get an error when reading a missing level", func() {
			level, err := DBLevel.ReadLevel(1, 1, 77)

			g.Assert(err == nil).IsTrue()
			g.Assert(level == nil).IsTrue()
		})

		g.It("Should fail when creating a level for a missing course", func() {
			level, err := DBLevel.CreateLevel(1, 99, 77, time.Now(), time.Now().AddDate(1, 0, 0))

			g.Assert(err != nil).IsTrue()
			g.Assert(level == nil).IsTrue()
		})

		g.It("Should be able to create a level", func() {
			level, err := DBLevel.CreateLevel(1, 1, 77, time.Now(), time.Now().AddDate(1, 0, 0))

			g.Assert(err == nil).IsTrue()
			g.Assert(level != nil).IsTrue()
		})

		g.It("Should be able to access a level", func() {
			level, err := DBLevel.ReadLevel(1, 1, 77)

			g.Assert(err == nil).IsTrue()
			g.Assert(level != nil).IsTrue()
		})

		g.It("Should be able to update a level", func() {
			level, err := DBLevel.UpdateLevel(1, 1, 77, time.Now().AddDate(1, 0, 0), time.Now().AddDate(2, 0, 0))

			g.Assert(err == nil).IsTrue()
			g.Assert(level != nil).IsTrue()
		})

		g.It("Should be able to remove a level", func() {
			count, err := DBLevel.DeleteLevel(1, 1, 77)

			g.Assert(err == nil).IsTrue()
			g.Assert(count != 0).IsTrue()
		})

		g.It("Should fail removing a missing level", func() {
			count, err := DBLevel.DeleteLevel(1, 1, 77)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 0).IsTrue()
		})

		g.It("Should fail to update a missing level", func() {
			level, err := DBLevel.UpdateLevel(1, 1, 77, time.Now().AddDate(1, 0, 0), time.Now().AddDate(2, 0, 0))

			g.Assert(err == nil).IsTrue()
			g.Assert(level == nil).IsTrue()
		})
	})

	g.Describe("When managing the module levels", func() {
		g.It("Should succed when adding a module to a level", func() {
			level, err := DBLevel.AddModule("-TEST-", 1, 1, 1, time.Now())

			g.Assert(err == nil).IsTrue()
			g.Assert(level != nil).IsTrue()
		})

		g.It("Should fail when adding a module to a level that already exists", func() {
			level, err := DBLevel.AddModule("-TEST-", 1, 1, 1, time.Now())

			g.Assert(err != nil).IsTrue()
			g.Assert(level == nil).IsTrue()
		})

		g.It("Should succed when deleting a module from a level", func() {
			count, err := DBLevel.RemoveModule("-TEST-", 1)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 1).IsTrue()
		})

		g.It("Should fail when deleting a module from a level that doesn't exist", func() {
			count, err := DBLevel.RemoveModule("-TEST-", 1)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 0).IsTrue()
		})
	})
}
