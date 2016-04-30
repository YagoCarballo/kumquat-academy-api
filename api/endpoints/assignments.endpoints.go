package endpoints

import (
	"net/http"
	"time"
	"github.com/YagoCarballo/kumquat.academy.api/database/models"
	"fmt"
	"mime/multipart"
	"os"
	"github.com/YagoCarballo/kumquat.academy.api/tools"
	"github.com/wayn3h0/go-uuid"
)

func CreateAssignment(
	title, description string,
	status models.AssignmentStatus,
	weight float64,
	start, end time.Time,
	moduleCode string,
) (int, map[string]interface{}) {
	assignment := models.Assignment{
		Title: title,
		Description: description,
		Status: status,
		Weight: weight,
		Start: start,
		End: end,
		ModuleCode: moduleCode,
	}

	dbAssignment, err := models.DBAssignments.CreateAssignment(assignment)
	if err != nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Error creating the assignment.",
		}
	}

	return http.StatusCreated, map[string]interface{}{
		"message": "Assignment created successfully",
		"assignment": dbAssignment,
	}
}

func GetAssignment(assignmentId uint32) (int, map[string]interface{}) {
	assignment, err := models.DBAssignments.ReadAssignment(assignmentId)
	if err != nil || assignment == nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "NotFound",
			"message": "Assignment not found.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"assignment": assignment,
	}
}

func UpdateAssignment(
	assignmentId uint32,
	title, description string,
	status models.AssignmentStatus,
	weight float64,
	start, end time.Time,
	moduleCode string,
) (int, map[string]interface{}) {

	assignment := models.Assignment{
		Title: title,
		Description: description,
		Status: status,
		Weight: weight,
		Start: start,
		End: end,
		ModuleCode: moduleCode,
	}

	dbAssignment, err := models.DBAssignments.UpdateAssignment(assignmentId, assignment)
	if err != nil || dbAssignment == nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Assignment not updated.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"assignment": dbAssignment,
	}
}

func DeleteAssignment(assignmentId uint32) (int, map[string]interface{}) {
	rows, err := models.DBAssignments.DeleteAssignment(assignmentId)
	if err != nil || rows <= 0 {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error deleting the assignment",
		}
	}

	return http.StatusAccepted, map[string]interface{}{
		"message": fmt.Sprintf("Assignment %d removed", assignmentId),
	}
}


func FindAssignmentsForModule(username, code string, readOnly bool) (int, map[string]interface{}) {
	assignments, err := models.DBAssignments.FindAssignmentsForModule(code, readOnly)
	if err != nil {
		return http.StatusNotFound, map[string]interface{}{
			"error": "NotFound",
			"message": "Assignments not found.",
		}
	}

	// Add all the weights to get a total weight
	var totalWeight float64
	for index, assignment := range assignments {
		totalWeight += assignment.Weight
		assignments[index].CanSubmit = models.DBPermissions.CanUserSubmitAssignment(username, assignment.ID)
	}

	return http.StatusOK, map[string]interface{}{
		"total_weight": totalWeight,
		"assignments": assignments,
	}
}

func UploadAssignmentAttachments(assignmentId uint32, file multipart.File, header *multipart.FileHeader) (int, FileResponseMessage) {
	status, response := UploadFile(file, header)
	if status != http.StatusOK {
		return http.StatusExpectationFailed, FileResponseMessage{
			Error: "Unknown",
			Message: "Error uploading the file.",
		}
	}

	count, err := models.DBAssignments.AddAttachmentToAssignment(assignmentId, response.Attachment.ID)
	if err != nil || count <= 0 {
		return http.StatusExpectationFailed, FileResponseMessage{
			Error: "Unknown",
			Message: "Error adding the attachment to the assignment.",
		}
	}

	return status, response
}

func RemoveAssignmentAttachments(assignmentId, attachmentId uint32) (int, map[string]interface{}) {
	var err error

	attachment, _ := models.DBAttachment.ReadAttachment(attachmentId)
	if attachment != nil {
		serverSettings := tools.GetSettings().Server
		path := fmt.Sprintf("%s/%s", serverSettings.UploadsPath, attachment.Url)
		os.Remove(path)
	}

	count, err := models.DBAssignments.RemoveAttachmentFromAssignment(assignmentId, attachmentId)
	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error deleting the attachment from the assignment.",
		}
	}

	if count <= 0 {
		return http.StatusNotFound, map[string]interface{}{
			"error": "NotFound",
			"message": "An attachment with that Id was not found inside the assignment.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"message": "Attachment removed.",
	}
}

func SubmitAssignment(username, description string, assignmentId uint32, fileHeaders []*multipart.FileHeader) (int, map[string]interface{}) {
	serverSettings := tools.GetSettings().Server

	user, err := models.DBUser.FindUser(username)
	if err != nil {
		return http.StatusUnauthorized, map[string]interface{}{
			"error": "Unautorized",
			"message": "This user was not found",
		}
	}

	token, err := uuid.NewV4()
	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "ExpectationFailed",
			"message": "Error when handling the uploaded file",
		}
	}

	zipPath := fmt.Sprintf("%s/%s", serverSettings.UploadsPath, token.String())
	_, err = tools.ZipFiles(zipPath, fileHeaders)
	if err != nil {
		fmt.Println(err)
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "ExpectationFailed",
			"message": "Unable to package the uploaded files.",
		}
	}

	// Upload Info to DB
	zipName := fmt.Sprintf("Assignment_%d_Student_%s.zip", assignmentId, user.MatricNumber)
	attachment, err := models.DBAttachment.CreateAttachment(zipName, "application/zip", token.String())
	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "ExpectationFailed",
			"message": "Error creating the attachment",
		}
	}

	// Register the submission into the Database
	_, err = models.DBAssignments.SubmitAssignment(user.ID, assignmentId, attachment.ID, description)
	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "ExpectationFailed",
			"message": "Error registering the submission",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"message": "Assignment Submitted",
		"attachment": attachment,
	}
}


func GradeAssignment(submissionId, grade uint32) (int, map[string]interface{}) {
	submission, err := models.DBAssignments.GradeAssignment(submissionId, grade)
	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "ExpectationFailed",
			"message": "Error grading the assignment",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"message": "Assignment Graded",
		"submission": submission,
	}
}
