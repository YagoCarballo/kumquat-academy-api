package models

import (
	"errors"
	"database/sql/driver"
)

const (
	ModuleDraft		ModuleStatus = "draft"
	ModuleOngoing	ModuleStatus = "ongoing"
	ModuleFuture	ModuleStatus = "future"
	ModuleEnded		ModuleStatus = "ended"

	AssignmentDraft		AssignmentStatus = "draft"
	AssignmentCreated	AssignmentStatus = "created"
	AssignmentAvailable	AssignmentStatus = "available"
	AssignmentSent		AssignmentStatus = "sent"
	AssignmentGraded	AssignmentStatus = "graded"
	AssignmentReturned	AssignmentStatus = "returned"

	SubmissionSent		SubmissionStatus = "sent"
	SubmissionReview	SubmissionStatus = "review"
	SubmissionGraded	SubmissionStatus = "graded"
	SubmissionCanceled	SubmissionStatus = "canceled"

	ExamComplete	SubmissionStatus = "complete"
	ExamReview		SubmissionStatus = "review"
	ExamGraded		SubmissionStatus = "graded"
)

type ModuleStatus string
type AssignmentStatus string
type SubmissionStatus string
type ExamStatus string

func (status *ModuleStatus) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	*status = ModuleStatus(string(asBytes))
	return nil
}

func (status ModuleStatus) Value() (driver.Value, error)  {
	return string(status), nil
}

func (status *AssignmentStatus) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	*status = AssignmentStatus(string(asBytes))
	return nil
}

func (status AssignmentStatus) Value() (driver.Value, error)  {
	return string(status), nil
}

func (status *SubmissionStatus) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	*status = SubmissionStatus(string(asBytes))
	return nil
}

func (status SubmissionStatus) Value() (driver.Value, error)  {
	return string(status), nil
}

func (status *ExamStatus) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	*status = ExamStatus(string(asBytes))
	return nil
}

func (status ExamStatus) Value() (driver.Value, error)  {
	return string(status), nil
}
