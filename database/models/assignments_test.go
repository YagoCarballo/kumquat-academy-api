package models

import (
	"testing"

	. "github.com/franela/goblin"
	"time"
)

func Test_Database_Assignments(t *testing.T) {
	g := Goblin(t)
	var assignmentId uint32

	g.Describe("When managing the assignments", func() {
		g.It("Should get an error when reading a missing attachment", func() {
			assignment, err := DBAssignments.ReadAssignment(99999)

			g.Assert(err == nil).IsTrue()
			g.Assert(assignment == nil).IsTrue()
		})

		g.It("Should get no assignments for a missing module", func() {
			assignments, err := DBAssignments.FindAssignmentsForModule("-null-", false)

			g.Assert(err == nil).IsTrue()
			g.Assert(assignments != nil).IsTrue()
			g.Assert(len(assignments) <= 0).IsTrue()
		})

		g.It("Should successfully create an assignment", func() {
			assignment, err := DBAssignments.CreateAssignment(Assignment{
				Title: "-test-asignment-",
				Description: "-test-asignment-",
				Status: AssignmentDraft,
				Weight: 0.15,
				Start: time.Now(),
				End: time.Now().AddDate(0, 2, 0),
				ModuleCode: "AC31007",
			})

			g.Assert(err == nil).IsTrue()
			g.Assert(assignment != nil).IsTrue()

			assignmentId = assignment.ID
		})

		g.It("Should be able to get an assignment", func() {
			assignment, err := DBAssignments.ReadAssignment(assignmentId)

			g.Assert(err == nil).IsTrue()
			g.Assert(assignment != nil).IsTrue()
		})

		g.It("Should get the assignments for module AC31007", func() {
			assignments, err := DBAssignments.FindAssignmentsForModule("AC31007", false)

			g.Assert(err == nil).IsTrue()
			g.Assert(assignments != nil).IsTrue()
			g.Assert(len(assignments) >= 1).IsTrue()
		})

		g.It("Should successfully update an assignment", func() {
			assignment, err := DBAssignments.UpdateAssignment(assignmentId, Assignment{
				Title: "-test-asignment-",
				Description: "-updated-asignment-",
				Status: AssignmentDraft,
				Weight: 0.15,
				Start: time.Now(),
				End: time.Now().AddDate(0, 2, 0),
				ModuleCode: "AC31007",
			})

			g.Assert(err == nil).IsTrue()
			g.Assert(assignment != nil).IsTrue()
			g.Assert(assignment.Description == "-updated-asignment-").IsTrue()
		})

		g.It("Should be able to delete an assignment", func() {
			count, err := DBAssignments.DeleteAssignment(assignmentId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 1).IsTrue()
		})

		g.It("Should not delete a missing assignment", func() {
			count, err := DBAssignments.DeleteAssignment(assignmentId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 0).IsTrue()
		})

		g.It("Should not update a deleted assignment", func() {
			assignment, err := DBAssignments.UpdateAssignment(assignmentId, Assignment{
				Title: "-test-asignment-",
				Description: "-deleted-asignment-",
				Status: AssignmentDraft,
				Weight: 0.15,
				Start: time.Now(),
				End: time.Now().AddDate(0, 2, 0),
				ModuleCode: "AC31007",
			})

			g.Assert(err == nil).IsTrue()
			g.Assert(assignment == nil).IsTrue()
		})

		g.It("Should get an error when reading a missing assignment", func() {
			module, err := DBAssignments.ReadAssignment(assignmentId)

			g.Assert(err == nil).IsTrue()
			g.Assert(module == nil).IsTrue()
		})
	})
}
