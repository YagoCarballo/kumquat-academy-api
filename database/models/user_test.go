package models

import (
	"testing"
	. "github.com/franela/goblin"
)

func Test_Database_Users(t *testing.T) {
	g := Goblin(t)

	var savedToken string

	g.Describe("When managing Users", func() {
		g.It("Should be able to find a user with it's username", func() {
			user, err := DBUser.FindUser("admin")

			g.Assert(err == nil).IsTrue()
			g.Assert(user != nil).IsTrue()
		})

		g.It("Should fail to find a missing user with it's username", func() {
			user, err := DBUser.FindUser("-missing-")

			g.Assert(err == nil).IsTrue()
			g.Assert(user == nil).IsTrue()
		})

		g.It("Should be able to find a user with it's email", func() {
			user, err := DBUser.FindUserWithEmail("jane.johnston68@example.com")

			g.Assert(err == nil).IsTrue()
			g.Assert(user != nil).IsTrue()
		})

		g.It("Should fail to find a missing user with it's email", func() {
			user, err := DBUser.FindUserWithEmail("null@gmail.com")

			g.Assert(err == nil).IsTrue()
			g.Assert(user == nil).IsTrue()
		})

		g.It("Should be able to search a users", func() {
			users, err := DBUser.SearchUsers("1", "AC31007")

			g.Assert(err == nil).IsTrue()
			g.Assert(users != nil).IsTrue()
			g.Assert(len(users) > 0).IsTrue()
		})

		g.It("Should not be able to search a missing users", func() {
			users, err := DBUser.SearchUsers("dhskjfahsg dkjfhgaskjhdfjkahs bdjh", "AC31007")

			g.Assert(err == nil).IsTrue()
			g.Assert(users != nil).IsTrue()
			g.Assert(len(users) == 0).IsTrue()
		})

		g.It("Should be able to find a user", func() {
			user, err := DBUser.FindUserWithId(1)

			g.Assert(err == nil).IsTrue()
			g.Assert(user != nil).IsTrue()
		})

		g.It("Should fail to find a missing user", func() {
			user, err := DBUser.FindUserWithId(999)

			g.Assert(err == nil).IsTrue()
			g.Assert(user == nil).IsTrue()
		})

		g.It("Should fail to trigger a reset password token on a missing user", func() {
			token, err := DBUser.AddResetPasswordToken(999)

			g.Assert(err != nil).IsTrue()
			g.Assert(token == nil).IsTrue()
		})

		g.It("Should be able to trigger a reset password token", func() {
			token, err := DBUser.AddResetPasswordToken(4)

			g.Assert(err == nil).IsTrue()
			g.Assert(token != nil).IsTrue()

			savedToken = *token
		})

		g.It("Should be able to change the password of a user", func() {
			success, err := DBUser.ChangePassword(savedToken, "$2a$10$ouCsus6K//.Xr04sNS0M9O1s8BXEDHdC9pFupCCup.leWdSlPn9hm")

			g.Assert(err == nil).IsTrue()
			g.Assert(success).IsTrue()
		})

		g.It("Should not be able to reuse an expired token", func() {
			success, err := DBUser.ChangePassword(savedToken, "$2a$10$ouCsus6K//.Xr04sNS0M9O1s8BXEDHdC9pFupCCup.leWdSlPn9hm")

			g.Assert(err != nil).IsTrue()
			g.Assert(!success).IsTrue()
		})
	})
}
