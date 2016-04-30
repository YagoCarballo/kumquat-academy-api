package models

import (
	"testing"

	. "github.com/franela/goblin"
)

func Test_Database_Attachments(t *testing.T) {
	g := Goblin(t)
	var attachmentId uint32

	g.Describe("When managing the attachments", func() {
		g.It("Should get an error when reading a missing attachment", func() {
			attachment, err := DBAttachment.ReadAttachment(99999)

			g.Assert(err == nil).IsTrue()
			g.Assert(attachment == nil).IsTrue()
		})

		g.It("Should successfully create an attachment", func() {
			attachment, err := DBAttachment.CreateAttachment("test-file.jpg", "image/JPG", "49860456-4ea5-42ff-8aac-476148e6f422")

			g.Assert(err == nil).IsTrue()
			g.Assert(attachment != nil).IsTrue()

			attachmentId = attachment.ID
		})

		g.It("Should be able to get an attachment", func() {
			attachment, err := DBAttachment.ReadAttachment(attachmentId)

			g.Assert(err == nil).IsTrue()
			g.Assert(attachment != nil).IsTrue()
		})

		g.It("Should be able to get an attachment with it's name", func() {
			attachment, err := DBAttachment.FindAttachment("49860456-4ea5-42ff-8aac-476148e6f422")

			g.Assert(err == nil).IsTrue()
			g.Assert(attachment != nil).IsTrue()
		})

		g.It("Should not be able to get a missing attachment with it's name", func() {
			attachment, err := DBAttachment.FindAttachment("-null-")

			g.Assert(err == nil).IsTrue()
			g.Assert(attachment == nil).IsTrue()
		})

		g.It("Should be able to delete an attachment", func() {
			count, err := DBAttachment.DeleteAttachment(attachmentId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 1).IsTrue()
		})

		g.It("Should not delete a missing attachment", func() {
			count, err := DBAttachment.DeleteAttachment(attachmentId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 0).IsTrue()
		})

		g.It("Should get an error when reading a missing attachment", func() {
			module, err := DBAttachment.ReadAttachment(attachmentId)

			g.Assert(err == nil).IsTrue()
			g.Assert(module == nil).IsTrue()
		})
	})
}
