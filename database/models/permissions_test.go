package models

import (
	"os"
	"testing"
	"net/http"

	"github.com/YagoCarballo/kumquat-academy-api/tools"

	. "github.com/franela/goblin"
	. "github.com/YagoCarballo/kumquat-academy-api/constants"
)

func Test_Database_Permissions(t *testing.T) {
	g := Goblin(t)

	g.Describe("When asking if a user has permissions", func() {
		g.It("should give an access denied for user guest accessing the module 1", func() {
			canRead		:= DBPermissions.IsActionPermittedOnModule(uint32(4), uint32(1), ReadPermission)
			canWrite	:= DBPermissions.IsActionPermittedOnModule(uint32(4), uint32(1), WritePermission)
			canDelete	:= DBPermissions.IsActionPermittedOnModule(uint32(4), uint32(1), DeletePermission)
			canUpdate	:= DBPermissions.IsActionPermittedOnModule(uint32(4), uint32(1), UpdatePermission)
			g.Assert(canRead).IsFalse()
			g.Assert(canWrite).IsFalse()
			g.Assert(canDelete).IsFalse()
			g.Assert(canUpdate).IsFalse()
		})

		g.It("should give an access permitted for user admin accessing any module", func() {
			canRead		:= DBPermissions.IsActionPermittedOnModule(uint32(1), uint32(1), ReadPermission)
			canWrite	:= DBPermissions.IsActionPermittedOnModule(uint32(1), uint32(1), WritePermission)
			canDelete	:= DBPermissions.IsActionPermittedOnModule(uint32(1), uint32(1), DeletePermission)
			canUpdate	:= DBPermissions.IsActionPermittedOnModule(uint32(1), uint32(1), UpdatePermission)
			g.Assert(canRead).IsTrue()
			g.Assert(canWrite).IsTrue()
			g.Assert(canDelete).IsTrue()
			g.Assert(canUpdate).IsTrue()

			canRead		= DBPermissions.IsActionPermittedOnModule(uint32(1), uint32(2), ReadPermission)
			canWrite	= DBPermissions.IsActionPermittedOnModule(uint32(1), uint32(2), WritePermission)
			canDelete	= DBPermissions.IsActionPermittedOnModule(uint32(1), uint32(2), DeletePermission)
			canUpdate	= DBPermissions.IsActionPermittedOnModule(uint32(1), uint32(2), UpdatePermission)
			g.Assert(canRead).IsTrue()
			g.Assert(canWrite).IsTrue()
			g.Assert(canDelete).IsTrue()
			g.Assert(canUpdate).IsTrue()

			canRead		= DBPermissions.IsActionPermittedOnModule(uint32(1), uint32(3), ReadPermission)
			canWrite	= DBPermissions.IsActionPermittedOnModule(uint32(1), uint32(3), WritePermission)
			canDelete	= DBPermissions.IsActionPermittedOnModule(uint32(1), uint32(3), DeletePermission)
			canUpdate	= DBPermissions.IsActionPermittedOnModule(uint32(1), uint32(3), UpdatePermission)
			g.Assert(canRead).IsTrue()
			g.Assert(canWrite).IsTrue()
			g.Assert(canDelete).IsTrue()
			g.Assert(canUpdate).IsTrue()
		})

		g.It("should have propper access rights for student", func() {
			canRead		:= DBPermissions.IsActionPermittedOnModule(uint32(3), uint32(1), ReadPermission)
			canWrite	:= DBPermissions.IsActionPermittedOnModule(uint32(3), uint32(1), WritePermission)
			canDelete	:= DBPermissions.IsActionPermittedOnModule(uint32(3), uint32(1), DeletePermission)
			canUpdate	:= DBPermissions.IsActionPermittedOnModule(uint32(3), uint32(1), UpdatePermission)
			g.Assert(canRead).IsTrue()
			g.Assert(canWrite).IsFalse()
			g.Assert(canDelete).IsFalse()
			g.Assert(canUpdate).IsFalse()
		})

		g.It("should have propper access rights for teacher", func() {
			canRead		:= DBPermissions.IsActionPermittedOnModule(uint32(2), uint32(1), ReadPermission)
			canWrite	:= DBPermissions.IsActionPermittedOnModule(uint32(2), uint32(1), WritePermission)
			canDelete	:= DBPermissions.IsActionPermittedOnModule(uint32(2), uint32(1), DeletePermission)
			canUpdate	:= DBPermissions.IsActionPermittedOnModule(uint32(2), uint32(1), UpdatePermission)
			g.Assert(canRead).IsTrue()
			g.Assert(canWrite).IsTrue()
			g.Assert(canDelete).IsTrue()
			g.Assert(canUpdate).IsTrue()
		})
	})

	g.Describe("When asking for permissions on a module, ", func () {
		g.It("Guest user should get access denied", func() {
			status, err := tools.VerifyAccess(uint32(1), uint32(4), ReadPermission, DBPermissions.IsActionPermittedOnModule)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(err != nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(1), uint32(4), WritePermission, DBPermissions.IsActionPermittedOnModule)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(err != nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(1), uint32(4), UpdatePermission, DBPermissions.IsActionPermittedOnModule)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(err != nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(1), uint32(4), DeletePermission, DBPermissions.IsActionPermittedOnModule)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(err != nil).IsTrue()
		})

		g.It("Admin user should get full access", func() {
			status, err := tools.VerifyAccess(uint32(1), uint32(1), ReadPermission, DBPermissions.IsActionPermittedOnModule)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(1), uint32(1), WritePermission, DBPermissions.IsActionPermittedOnModule)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(1), uint32(1), UpdatePermission, DBPermissions.IsActionPermittedOnModule)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(1), uint32(1), DeletePermission, DBPermissions.IsActionPermittedOnModule)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()
		})
	})

	g.Describe("When asking for permissions on a module with code, ", func () {
		g.It("Guest user should get access denied", func() {
			status, err := tools.VerifyAccess("AC31007", uint32(4), ReadPermission, DBPermissions.IsActionPermittedOnModuleWithCode)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(err != nil).IsTrue()

			status, err = tools.VerifyAccess("AC31007", uint32(4), WritePermission, DBPermissions.IsActionPermittedOnModuleWithCode)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(err != nil).IsTrue()

			status, err = tools.VerifyAccess("AC31007", uint32(4), UpdatePermission, DBPermissions.IsActionPermittedOnModuleWithCode)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(err != nil).IsTrue()

			status, err = tools.VerifyAccess("AC31007", uint32(4), DeletePermission, DBPermissions.IsActionPermittedOnModuleWithCode)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(err != nil).IsTrue()
		})

		g.It("Admin user should get full access", func() {
			status, err := tools.VerifyAccess("AC31007", uint32(1), ReadPermission, DBPermissions.IsActionPermittedOnModuleWithCode)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()

			status, err = tools.VerifyAccess("AC31007", uint32(1), WritePermission, DBPermissions.IsActionPermittedOnModuleWithCode)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()

			status, err = tools.VerifyAccess("AC31007", uint32(1), UpdatePermission, DBPermissions.IsActionPermittedOnModuleWithCode)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()

			status, err = tools.VerifyAccess("AC31007", uint32(1), DeletePermission, DBPermissions.IsActionPermittedOnModuleWithCode)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()
		})
	})

	g.Describe("When asking for permissions on a course, ", func () {
		g.It("Guest user should get access denied", func() {
			status, err := tools.VerifyAccess(uint32(1), uint32(4), ReadPermission, DBPermissions.IsActionPermittedOnCourse)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(err != nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(1), uint32(4), WritePermission, DBPermissions.IsActionPermittedOnCourse)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(err != nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(1), uint32(4), UpdatePermission, DBPermissions.IsActionPermittedOnCourse)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(err != nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(1), uint32(4), DeletePermission, DBPermissions.IsActionPermittedOnCourse)
			g.Assert(status).Equal(http.StatusForbidden)
			g.Assert(err != nil).IsTrue()
		})

		g.It("Teacher user should get partial access", func() {
			status, err := tools.VerifyAccess(uint32(2), uint32(2), ReadPermission, DBPermissions.IsActionPermittedOnCourse)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(2), uint32(2), WritePermission, DBPermissions.IsActionPermittedOnCourse)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(2), uint32(2), UpdatePermission, DBPermissions.IsActionPermittedOnCourse)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(2), uint32(2), DeletePermission, DBPermissions.IsActionPermittedOnCourse)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()
		})

		g.It("Admin user should get full access", func() {
			status, err := tools.VerifyAccess(uint32(1), uint32(1), ReadPermission, DBPermissions.IsActionPermittedOnCourse)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(1), uint32(1), WritePermission, DBPermissions.IsActionPermittedOnCourse)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(1), uint32(1), UpdatePermission, DBPermissions.IsActionPermittedOnCourse)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()

			status, err = tools.VerifyAccess(uint32(1), uint32(1), DeletePermission, DBPermissions.IsActionPermittedOnCourse)
			g.Assert(status).Equal(http.StatusOK)
			g.Assert(err == nil).IsTrue()
		})
	})

	// Remove the Temporary Settings
	os.Remove(SETTINGS_PATH)
}
