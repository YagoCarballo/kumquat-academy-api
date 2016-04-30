package constants

const (
	ReadPermission		AccessRight = "READ"
	WritePermission		AccessRight = "WRITE"
	DeletePermission	AccessRight = "DELETE"
	UpdatePermission	AccessRight = "UPDATE"

	ModuleResource		ResourceType = "module"
	CourseResource		ResourceType = "course"
)

type AccessRight string
type ResourceType string