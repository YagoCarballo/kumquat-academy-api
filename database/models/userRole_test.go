package models

import (
	"testing"
	. "github.com/franela/goblin"
)

func Test_Database_UserRoles(t *testing.T) {
	g := Goblin(t)

	g.Describe("When managing User Roles", func() {
		g.It("Should be able to find a user role with it's name", func() {
			user, err := DBUserRole.FindUserRole("Student")

			g.Assert(err == nil).IsTrue()
			g.Assert(user != nil).IsTrue()
		})

		g.It("Should fail to find a missing user role with it's name", func() {
			user, err := DBUserRole.FindUserRole("-missing-")

			g.Assert(err == nil).IsTrue()
			g.Assert(user == nil).IsTrue()
		})
	})
}
