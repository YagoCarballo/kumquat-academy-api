package models

import (
	"github.com/jinzhu/gorm"
	"github.com/YagoCarballo/kumquat-academy-api/database"
	"time"
)

type AssignmentsModel struct{}
var DBAssignments AssignmentsModel

func (model AssignmentsModel) DB() *gorm.DB {
	return database.DB
}

func (model AssignmentsModel) CreateAssignment(assignment Assignment) (*Assignment, error) {
	query := model.DB().Create(&assignment)
	if query.Error != nil {
		return nil, query.Error
	}

	return &assignment, nil
}

func (model AssignmentsModel) ReadAssignment(id uint32) (*Assignment, error) {
	var assignment Assignment

	query := model.DB().Preload("Attachments").First(&assignment, "id = ?", id)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	return &assignment, nil
}

func (model AssignmentsModel) UpdateAssignment(id uint32, assignment Assignment) (*Assignment, error) {
	query := model.DB().Table("assignments").Where("id = ?", id).Update(&assignment)
	if query.Error != nil {
		return nil, query.Error
	}


	if query.RowsAffected <= 0 {
		return nil, nil
	}

	assignment.ID = id
	return &assignment, nil
}

func (model AssignmentsModel) DeleteAssignment(id uint32) (int64, error) {
	query := model.DB().
		Table("assignments").
		Where("id = ?", id).
		Delete(Assignment{})
	if query.Error != nil {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}

func (model AssignmentsModel) FindAssignmentsForModule(code string, readOnly bool) ([]Assignment, error) {
	var assignments []Assignment

	query := model.DB().Preload("Attachments").Scopes(validAssignmentStatus(readOnly)).Find(&assignments, "module_code = ?", code)
	if query.Error != nil {
		return nil, query.Error
	}

	// If this is a teacher, add a list of students pending to submit the assignment
	if !readOnly {
		for index, assignment := range assignments {
			assignments[index].Students = []map[string]interface{}{}

			students, err := DBModule.FindStudentsForModule(code, "Student")
			if err != nil {
				return assignments, err;
			}

			for _, student := range students {
				studentMap := map[string]interface{}{
					"id": 				student.ID,
					"first_name":		student.FirstName,
					"last_name":		student.LastName,
					"username":			student.Username,
					"email":			student.Email,
					"matric_number":	student.MatricNumber,
					"matric_date":		student.MatricDate,
					"date_of_birth":	student.DateOfBirth,
					"admin":			student.Admin,
					"avatar_id":		student.AvatarId,
					"avatar":			nil,
					"submission":		nil,
				};

				if student.Avatar != nil {
					studentMap["avatar"] = student.Avatar.Url
				}

				// Get the submission
				submission := &Submission{}
				query := model.DB().Preload("Attachment").Where("user_id = ? and assignment_id = ?", student.ID, assignment.ID).First(submission)
				if query.Error == nil && submission != nil {
					studentMap["submission"] = submission
				}

				assignments[index].Students = append(assignments[index].Students, studentMap)
			}
		}
	}

	return assignments, nil
}

func validAssignmentStatus(readOnly bool) func (db *gorm.DB) *gorm.DB {
	var validStatus []AssignmentStatus

	if readOnly {
		validStatus = []AssignmentStatus{
			AssignmentAvailable,
			AssignmentSent,
			AssignmentGraded,
			AssignmentReturned,
		}
	} else {
		validStatus = []AssignmentStatus{
			AssignmentDraft,
			AssignmentCreated,
			AssignmentAvailable,
			AssignmentSent,
			AssignmentGraded,
			AssignmentReturned,
		}
	}

	return func (db *gorm.DB) *gorm.DB {
		return db.Where("status in (?)", validStatus)
	}
}

func (model AssignmentsModel) AddAttachmentToAssignment(assignmentId, attachmentId uint32) (int64, error) {
	assignmentAttachment := AssignmentAttachments{
		AssignmentID: assignmentId,
		AttachmentID: attachmentId,
	}

	query := model.DB().Create(&assignmentAttachment)
	if query.Error != nil {
		return query.RowsAffected, query.Error
	}

	return query.RowsAffected, nil
}

func (model AssignmentsModel) RemoveAttachmentFromAssignment(assignmentId, attachmentId uint32) (int64, error) {
	query := model.DB().
				Table("assignment_attachments").
				Where("assignment_id = ? and attachment_id = ?", assignmentId, attachmentId).
				Delete(AssignmentAttachments{})
	if query.Error != nil {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}


func (model AssignmentsModel) SubmitAssignment(userId, assignmentId, attachmentId uint32, description string) (*Submission, error) {
	submission := Submission{
		UserID: userId,
		AssignmentID: assignmentId,
		AttachmentID: attachmentId,
		Status: SubmissionSent,
		Description: description,
		SubmittedOn: time.Now(),
	}

	query := model.DB().Create(&submission)
	if query.Error != nil {
		return nil, query.Error
	}

	return &submission, nil
}

func (model AssignmentsModel) GradeAssignment(submissionId, grade uint32) (*Submission, error) {
	var submission Submission
	gradedOn := time.Now()

	query := DBUser.DB().Table("submissions").Where("id = ?", submissionId).Updates(map[string]interface{}{
		"grade": grade,
		"graded_on": &gradedOn,
	})
	if query.Error != nil {
		return &submission, query.Error
	}

	// Get the submission
	query = model.DB().Preload("Attachment").Where("id = ?", submissionId).First(&submission)
	if query.Error == nil {
		return &submission, query.Error
	}

	return &submission, nil
}
