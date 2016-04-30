package models

import (
	"time"
)

// Models Declaration

type User struct {
	ID           uint32 	`json:"id" sql:"size:11; AUTO_INCREMENT" gorm:"primary_key"`
	FirstName    string 	`json:"first_name" sql:"not null; size:50"`
	LastName     string 	`json:"last_name" sql:"not null; size:50"`
	Username     string 	`json:"username" sql:"not null; size:30; unique"`
	Email        string 	`json:"email" sql:"not null; size:100; unique"`
	Password     string 	`json:"-" sql:"not null; size:100"`
	Active       bool		`json:"active" sql:"not null"`
	Admin        bool		`json:"admin" sql:"not null; default:false"`
	DateOfBirth  time.Time	`json:"date_of_birth" sql:"not null"`
	MatricNumber string		`json:"matric_number" sql:"not null; size:30"`
	MatricDate   time.Time	`json:"matric_date"`
	CreatedAt    time.Time	`json:"created_at"`
	UpdatedAt    time.Time	`json:"updated_at"`
	Sessions     []Session	`json:"sessions,omitempty"`
	AvatarId	 uint32		`json:"avatar_id"`
	Avatar	 	 *Attachment`json:"avatar,omitempty"`
}

type Session struct {
	Token     string	`json:"token" gorm:"primary_key"`
	DeviceID  string	`json:"device_id" sql:"not null"`
	ExpiresIn time.Time	`json:"expires_in"`
	CreatedOn time.Time	`json:"created_on"`
	UserID    uint32    	`json:"user_id" sql:"not null"`
	User	  *User		`json:"user,omitempty"`
}

type Course struct {
	ID     		uint32	`json:"id" gorm:"primary_key"`
	Title  		string	`json:"title" sql:"not null;unique_index"`
	Description  	string	`json:"description"`
	LevelModules 	[]LevelModule `json:"-"`
	Modules		[]*OutputModule `json:"modules,omitempty"`
	Role		*PermissionsTable `json:"role,omitempty"`
}

type Module struct {
	ID     		uint32	`json:"id" gorm:"primary_key"`
	Title  		string	`json:"title" sql:"not null"`
	Description string	`json:"description"`
	Color 		string	`json:"color"`
	Icon		string	`json:"icon" sql:"default:\"icon-book\""`
	Duration	uint32 	`json:"duration"`
	Assignments []Assignment `json:"assignments,omitempty"`
	Permissions *PermissionsTable `json:"role,omitempty"`
}

type Role struct {
	ID          uint32	`json:"id" gorm:"primary_key"`
	Name        string  `json:"name" sql:"not null"`
	Description string	`json:"description"`
	CanRead     bool	`json:"read" sql:"not null;default:false"`
	CanWrite    bool	`json:"write" sql:"not null;default:false"`
	CanDelete   bool	`json:"delete" sql:"not null;default:false"`
	CanUpdate   bool	`json:"update" sql:"not null;default:false"`
}

type Class struct {
	ID    		uint32	`json:"id" gorm:"primary_key"`
	CourseID    	uint32	`json:"course_id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	Title  		string	`json:"title" sql:"not null"`
	Start  		time.Time`json:"start"`
	End  		time.Time`json:"end"`
	Levels		[]*CourseLevel `json:"levels,omitempty"`
}

type CourseLevel struct {
	Level		uint32	`json:"level" gorm:"primary_key" sql:"type:int(10) unsigned"`
	ClassID    	uint32	`json:"class_id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	Class		*Class	`json:"class,omitempty"`
	CourseID    	uint32	`json:"course_id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	Course		*Course	`json:"course,omitempty"`

	Start 		time.Time`json:"start"`
	End 		time.Time`json:"end"`
}

type LevelModule struct {
	Code  		string	`json:"code" gorm:"primary_key"`

	ClassID    	uint32  `json:"class_id,omitempty" gorm:"primary_key" sql:"type:int(10) unsigned"`
	Class		*Class	`json:"class,omitempty"`

	Level		uint32  `json:"level,omitempty" sql:"not null"`
	CourseLevel *CourseLevel `json:"course_level,omitempty"`

	ModuleID 	uint32	`json:"module_id,omitempty" sql:"not null"`
	Module	 	*Module `json:"module,omitempty"`

	Status		ModuleStatus `json:"status,omitempty" sql:"not null"`
	Start		time.Time `json:"start" sql:"not null"`
}

type UserModule struct {
	UserID   uint32	`json:"user_id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	User	 *User	`json:"user,omitempty"`

	ModuleCode string `json:"module_code" gorm:"primary_key"`
	Module	 *Module `json:"module,omitempty"`

	RoleID   uint32	`json:"role_id" sql:"not null"`
	Role	 *Role	`json:"role,omitempty"`

	ClassID  uint32	`json:"class_id" sql:"not null"`
	Class	 *Class	`json:"class,omitempty"`
}

type UserCourse struct {
	UserID   uint32	`json:"user_id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	User	 *User	`json:"user,omitempty"`

	CourseID uint32	`json:"course_id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	Course	 *Course`json:"course,omitempty"`

	RoleID   uint32	`json:"role_id" sql:"not null"`
	Role	 *Role	`json:"role,omitempty"`
}

type Assignment struct {
	ID     			uint32	`json:"id" 'gorm:"primary_key"`
	Title    		string	`json:"title" sql:"not null"`
	Description    	string	`json:"description" sql:"type:varchar(4096); not null"`
	Status			AssignmentStatus `json:"status"`
	Weight			float64	`json:"weight"`
	Start   		time.Time `json:"start"`
	End   			time.Time `json:"end"`

	ModuleCode    	string	`json:"module_code" sql:"not null"`
	Module			*Module	`json:"module,omitempty"`

	Attachments		[]Attachment `json:"attachments,omitempty"gorm:"many2many:assignment_attachments;"`

	CanSubmit		bool `json:"submission_open,omitempty" sql:"-"`
	Students		[]map[string]interface{} `json:"students,omitempty" sql:"-"`
}

type Exam struct {
	ID     			uint32	`json:"id' gorm:"primary_key"`
	Topic    		string	`json:"topic" sql:"not null"`
	Location    	string	`json:"location" sql:"not null"`
	Weight			float64	`json:"weight"`
	Date   			time.Time `json:"date"`

	ModuleCode    	string	`json:"module_code" sql:"not null"`
	Module			*Module	`json:"module,omitempty"`

	AttachmentID    uint32	`json:"attachment_id"`
	Attachment		*Attachment `json:"attachment,omitempty"`
}

type Attachment struct {
	ID     uint32	`json:"id" gorm:"primary_key"`
	Name   string	`json:"name" sql:"not null"`
	Type   string	`json:"type" sql:"not null"`
	Url    string	`json:"url" sql:"not null"`
}

type AssignmentAttachments struct {
	AttachmentID    uint32	`json:"attachment_id"`
	Attachment		*Attachment `json:"attachment,omitempty"`

	AssignmentID    uint32	`json:"assignment_id"`
	Assignment		*Assignment `json:"assignment,omitempty"`
}

type Page struct {
	ID     		uint32	`json:"id" gorm:"primary_key"`
	ModuleID    uint32	`json:"module_id" sql:"not null"`
	Title   	string	`json:"title" sql:"not null"`
	Content   	string	`json:"content" sql:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Lecture struct {
	ID            uint32 `json:"id" gorm:"primary_key"`
	Canceled      bool	 `json:"canceled"`

	ModuleID      uint32  `json:"module_id,omitempty" sql:"not null"`
	Module        *Module `json:"module,omitempty"`

	LectureSlotID *uint32 `json:"lecture_slot_id,omitempty"`
	LectureSlot   *LectureSlot `json:"slot,omitempty"`

	Location      string	`json:"location" sql:"not null"`
	Description   string	`json:"description" sql:"type:varchar(4096); not null"`
	Topic         string	`json:"topic"`
	Start         time.Time `json:"start"`
	End           time.Time `json:"end"`

	Attachments		[]Attachment `json:"attachments,omitempty"gorm:"many2many:lecture_attachments;"`
}

type LectureAttachments struct {
	AttachmentID    uint32	`json:"attachment_id"`
	Attachment		*Attachment `json:"attachment,omitempty"`

	LectureID    uint32	`json:"lecture_id"`
	Lecture		*Lecture `json:"lecture,omitempty"`
}

type LectureSlot struct {
	ID     		uint32	`json:"id" gorm:"primary_key"`

	ModuleID    uint32	`json:"module_id,omitempty" sql:"not null"`
	Module		*Module `json:"module,omitempty"`

	Location   	string	`json:"location" sql:"not null"`
	Type   		string	`json:"type"`
	Start   	time.Time `json:"start"`
	End   		time.Time `json:"end"`
}

type Materials struct {
	ID     			uint32	`json:"id" gorm:"primary_key"`
	Type   			string	`json:"type"`

	ModuleID    	uint32	`json:"module_id" sql:"not null"`
	Module			*Module	`json:"module,omitempty"`

	LectureID   	uint32	`json:"lecture_id"`
	Lecture			*Lecture `json:"lecture,omitempty"`

	AttachmentID    uint32	`json:"attachment_id"`
	Attachment		*Attachment	`json:"attachment"`
}

type Submission struct {
	ID     			uint32	`json:"id" gorm:"primary_key"`
	Grade     		float64	`json:"grade"`
	Status     		SubmissionStatus `json:"status"`
	Description   	string	`json:"description" sql:"type:varchar(4096); not null"`
	SubmittedOn	   	time.Time `json:"submitted_on" sql:"not null"`
	GradedOn	   	*time.Time `json:"graded_on"`

	UserID			uint32	`json:"user_id" sql:"not null"`
	User 			*User	`json:"user,omitempty"`

	AssignmentID	uint32	`json:"assignment_id" sql:"not null"`
	Assignment		*Assignment `json:"assignment,omitempty"`

	AttachmentID	uint32	`json:"attachment_id" sql:"not null"`
	Attachment		*Attachment `json:"attachment,omitempty"`
}

type StudentExam struct {
	ExamID     		uint32	`json:"id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	Exam			*Exam	`json:"exam,omitempty"`

	UserID			uint32	`json:"user_id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	User 			*User	`json:"user,omitempty"`

	Grade     		float64	`json:"grade"`
	Status			ExamStatus `json:"status"`
}

type Task struct {
	ID     			uint32	`json:"id" gorm:"primary_key"`

	AssignmentID	uint32	`json:"assignment_id" sql:"not null"`
	Assignment		*Assignment `json:"assignment,omitempty"`
}

type CompletedTask struct {
	UserID    uint32 `json:"user_id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	User	*User `json:"user,omitempty"`

	TaskID	uint32 `json:"task_id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	Task	*Task `json:"task,omitempty"`

	Date	time.Time `json:"date"`
}

type Team struct {
	ID     			uint32	`json:"id" gorm:"primary_key"`

	AssignmentID	uint32	`json:"assignment_id" sql:"not null"`
	Assignment		*Assignment `json:"assignment,omitempty"`
}

type TeamMember struct {
	TeamID	uint32 `json:"team_id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	Team	*Team `json:"team,omitempty"`

	UserID	uint32 `json:"user_id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	User	*User `json:"user,omitempty"`
}

type TeamCompletedTask struct {
	TeamID  uint32 `json:"team_id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	Team	*Team `json:"team,omitempty"`

	TaskID	uint32	`json:"task_id" gorm:"primary_key" sql:"type:int(10) unsigned"`
	Task	*Task `json:"task,omitempty"`

	Date	time.Time `json:"date"`
}

type Announcement struct {
	ID				uint32	`json:"id" gorm:"primary_key" sql:"type:int(10) unsigned"`

	UserID			uint32 `json:"user_id"`
	User			*User `json:"user,omitempty"`

	ModuleID		uint32 `json:"module_id"`
	Module			*Module `json:"module,omitempty"`

	AssignmentID	uint32 `json:"assignment_id"`
	Assignment		*Assignment `json:"assignment,omitempty"`

	CourseID 		uint32 `json:"course_id"`
	Course	 		*Course `json:"course,omitempty"`
}

type ResetPassword struct {
	Token			string	`json:"token" gorm:"primary_key" sql:"type:varchar(255)"`

	UserID			uint32 `json:"user_id"`
	User			*User `json:"user,omitempty"`

	Expires	 		time.Time `json:"expires"`
}
