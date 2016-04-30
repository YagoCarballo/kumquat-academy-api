package models

import (
	"testing"
	. "github.com/franela/goblin"
)

func Test_Database_Sessions(t *testing.T) {
	g := Goblin(t)
	var token string

	g.Describe("When managing sessions", func() {
		g.It("Should faile to find a missing session", func() {
			session, err := DBSession.FindSession("-missing-")

			g.Assert(err == nil).IsTrue()
			g.Assert(session == nil).IsTrue()
		})

		g.It("Should be able to create a new session", func() {
			session, err := DBSession.Create(1, "-test-")

			g.Assert(err == nil).IsTrue()
			g.Assert(session != nil).IsTrue()
			token = session.Token
		})

		g.It("Should be able to find a session", func() {
			session, err := DBSession.FindSession(token)

			g.Assert(err == nil).IsTrue()
			g.Assert(session != nil).IsTrue()
		})

		g.It("Should be able to find a session for a user on a device", func() {
			session, err := DBSession.findSessionForUserOnDevice(1, "-test-")

			g.Assert(err == nil).IsTrue()
			g.Assert(session != nil).IsTrue()
			g.Assert(session.Token == token).IsTrue()
		})

		g.It("Should be able to remove a session", func() {
			count, err := DBSession.RemoveSession(token)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 1).IsTrue()
		})
	})
}
