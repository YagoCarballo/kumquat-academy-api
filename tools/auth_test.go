package tools
import (
	"testing"
	"net/http"

	. "github.com/franela/goblin"
)

func Test_Auth_Tools(t *testing.T) {
	g := Goblin(t)

	g.Describe("When using the auth tools, ", func () {
		g.It("should give an error when parsing an invalid ID", func() {
			id, status, err := ParseID("-Invalid-Id-")
			g.Assert(id).Equal(uint32(0))
			g.Assert(status).Equal(http.StatusConflict)
			g.Assert(err != nil).IsTrue()
		})

		g.It("should parse a valid ID", func() {
			id, status, err := ParseID("5")
			g.Assert(id).Equal(uint32(5))
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()
		})
	})
}