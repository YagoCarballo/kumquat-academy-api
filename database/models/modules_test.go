package models

import (
	"testing"

	. "github.com/franela/goblin"
)

func Test_Database_Modules(t *testing.T) {
	g := Goblin(t)
	var moduleId uint32

	g.Describe("When managing the modules", func() {
		g.It("Should get an error when reading a missing module", func() {
			module, err := DBModule.FindModule(99999)

			g.Assert(err == nil).IsTrue()
			g.Assert(module == nil).IsTrue()
		})

		g.It("Should successfully create a module", func() {
			module, err := DBModule.CreateModule(Module{
				Title: "-test-module-",
				Description: "-module-created-through-tests-",
				Color: "red",
				Icon: "fa-cloud",
				Duration: 6,
			})

			g.Assert(err == nil).IsTrue()
			g.Assert(module != nil).IsTrue()

			moduleId = module.ID
		})

		g.It("Should be able to get a module", func() {
			module, err := DBModule.FindModule(moduleId)

			g.Assert(err == nil).IsTrue()
			g.Assert(module != nil).IsTrue()
		})

		g.It("Should be able to delete a module", func() {
			count, err := DBModule.DeleteModule(moduleId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 1).IsTrue()
		})

		g.It("Should not delete a missing module", func() {
			count, err := DBModule.DeleteModule(moduleId)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 0).IsTrue()
		})

		g.It("Should get an error when reading a missing module", func() {
			module, err := DBModule.FindModule(moduleId)

			g.Assert(err == nil).IsTrue()
			g.Assert(module == nil).IsTrue()
		})
	})

	g.Describe("When accessing the modules", func() {
		g.It("Should be able to get a module with it's code", func() {
			module, err := DBModule.FindModuleWithCode("AC22001")

			g.Assert(err == nil).IsTrue()
			g.Assert(module != nil).IsTrue()
		})

		g.It("Should break when accessing a module with an invalid code", func() {
			module, err := DBModule.FindModuleWithCode("nop' or true); -- ")

			g.Assert(err == nil).IsTrue()
			g.Assert(module == nil).IsTrue()
		})

		g.It("Should be able to get a list of modules for a user", func() {
			modules, err := DBModule.FindModulesForUser("admin")

			g.Assert(err == nil).IsTrue()
			g.Assert(modules != nil).IsTrue()
			g.Assert(len(modules) >= 3).IsTrue()
		})

		g.It("Should not be able to get a list of modules for a missing user", func() {
			modules, err := DBModule.FindModulesForUser("-john-doe-")

			g.Assert(err == nil).IsTrue()
			g.Assert(modules != nil).IsTrue()
			g.Assert(len(modules) == 0).IsTrue()
		})

		g.It("Should be able to get a list of modules for a level", func() {
			modules, err := DBModule.FindModulesForLevel(1, 1)

			g.Assert(err == nil).IsTrue()
			g.Assert(modules != nil).IsTrue()
			g.Assert(len(modules) >= 1).IsTrue()
		})

		g.It("Should not be able to get a list of modules for a missing level", func() {
			modules, err := DBModule.FindModulesForLevel(1, 999)

			g.Assert(err == nil).IsTrue()
			g.Assert(modules != nil).IsTrue()
			g.Assert(len(modules) == 0).IsTrue()
		})

		g.It("Should not be able to get a list of modules for a missing class", func() {
			modules, err := DBModule.FindModulesForLevel(999, 1)

			g.Assert(err == nil).IsTrue()
			g.Assert(modules != nil).IsTrue()
			g.Assert(len(modules) == 0).IsTrue()
		})

		g.It("Should be able to get module for a level", func() {
			module, err := DBModule.GetLevelModel(1, 1, 1)

			g.Assert(err == nil).IsTrue()
			g.Assert(module != nil).IsTrue()
		})
	})

	g.Describe("When accessing the raw modules", func() {
		g.It("Should be able to get a list of all the raw modules", func() {
			modules, err := DBModule.FindRawModules("", 0)

			g.Assert(err == nil).IsTrue()
			g.Assert(modules != nil).IsTrue()
			g.Assert(len(modules) >= 3).IsTrue()
		})

		g.It("Should be able to get query the raw modules", func() {
			modules, err := DBModule.FindRawModules("Graphics", 0)

			g.Assert(err == nil).IsTrue()
			g.Assert(modules != nil).IsTrue()
			g.Assert(len(modules) == 1).IsTrue()
		})
	})

	g.Describe("When accessing the students for a module", func() {
		g.It("Should be able to get the students for a module", func() {
			students, err := DBModule.FindStudentsForModule("AC31007", "Student")

			g.Assert(err == nil).IsTrue()
			g.Assert(students != nil).IsTrue()
			g.Assert(len(students) > 0).IsTrue()
		})

		g.It("Should be able to get a student from a module", func() {
			student, err := DBModule.GetModuleStudent(3, "AC31007")

			g.Assert(err == nil).IsTrue()
			g.Assert(student != nil).IsTrue()
		})

		g.It("Should not be able to get a missing student from a module", func() {
			student, err := DBModule.GetModuleStudent(99, "AC31007")

			g.Assert(err == nil).IsTrue()
			g.Assert(student == nil).IsTrue()
		})

		g.It("Should be able to add a student to a module", func() {
			student, err := DBModule.AddStudentToModule("AC31007", 1)

			g.Assert(err == nil).IsTrue()
			g.Assert(student != nil).IsTrue()
		})

		g.It("Should not be able to add an existent student to a module", func() {
			student, err := DBModule.AddStudentToModule("AC31007", 1)

			g.Assert(err == nil).IsTrue()
			g.Assert(student == nil).IsTrue()
		})

		g.It("Should be able to remove a student from a module", func() {
			count, err := DBModule.RemoveStudentFromModule("AC31007", 1)

			g.Assert(err == nil).IsTrue()
			g.Assert(count > 0).IsTrue()
		})

		g.It("Should not be able to remove a missing student from a module", func() {
			count, err := DBModule.RemoveStudentFromModule("AC31007", 999)

			g.Assert(err == nil).IsTrue()
			g.Assert(count == 0).IsTrue()
		})
	})
}
