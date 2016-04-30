package endpoints

import (
	"testing"

	. "github.com/franela/goblin"
)

func Test_API(t *testing.T) {
	g := Goblin(t)
	g.Describe("Tests the Say Hello Endpoint", func() {

		g.It("Should say Hello, World!!", func() {
			status, data := SayHello("World")
			g.Assert(status).Equal(200)
			g.Assert(data["message"]).Equal("Hello, World!!")
		})

		g.It("Should say Hello, John!!", func() {
			status, data := SayHello("John")
			g.Assert(status).Equal(200)
			g.Assert(data["message"]).Equal("Hello, John!!")
		})
	})
}
