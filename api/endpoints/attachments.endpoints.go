package endpoints

import (
	"os"
	"io"
	"net/http"
	"mime/multipart"
	"fmt"
	"github.com/wayn3h0/go-uuid"
	"github.com/YagoCarballo/kumquat.academy.api/database/models"
	"github.com/YagoCarballo/kumquat.academy.api/tools"
	"bufio"
)

type (
	MultipleFilesResponse struct  {
		Status int `json:"status"`
		Message string `json:"message"`
		Attachments []*models.Attachment `json:"files"`
		Messages []FileResponse `json:"files"`
	}

	FileResponse struct  {
		Status int `json:"status"`
		Message FileResponseMessage `json:"message"`
	}

	FileResponseMessage struct {
		Message string `json:"message"`
		Error string `json:"error,omitempty"`
		Attachment *models.Attachment `json:"attachment,omitempty"`
	}
)

func UploadFile(file multipart.File, header *multipart.FileHeader) (int, FileResponseMessage) {
	serverSettings := tools.GetSettings().Server
	mimeType := header.Header.Get("Content-Type")

	token, err := uuid.NewV4()
	if err != nil {
		return http.StatusExpectationFailed, FileResponseMessage{
			Error: "ExpectationFailed",
			Message: "Error when handling the uploaded file",
		}
	}

	path := fmt.Sprintf("%s/%s", serverSettings.UploadsPath, token.String())
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		fmt.Println(err)
		return http.StatusUnauthorized, FileResponseMessage{
			Error: "Unauthorized",
			Message: "Unable to create the file for writing. Check your write access privilege",
		}
	}

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		return http.StatusExpectationFailed, FileResponseMessage{
			Error: "ExpectationFailed",
			Message: "Error when handling the uploaded file",
		}
	}

	// Upload Info to DB
	attachment, err := models.DBAttachment.CreateAttachment(header.Filename, mimeType, token.String())
	if err != nil {
		return http.StatusExpectationFailed, FileResponseMessage{
			Error: "ExpectationFailed",
			Message: "Error creating the attachment",
		}
	}

	return http.StatusOK, FileResponseMessage{
		Message: "File uploaded successfully",
		Attachment: attachment,
	}
}

func UploadFiles(fileHeaders []*multipart.FileHeader) (int, map[string]interface{}) {
	output := MultipleFilesResponse{
		Status: http.StatusOK,
		Message: "Files uploaded Successfully.",
	}

	for index, _ := range fileHeaders { // loop through the files one by one
		file, err := fileHeaders[index].Open()
		defer file.Close()
		if err != nil {
			fileResponse := FileResponse{ Status: http.StatusExpectationFailed, Message: FileResponseMessage{
				Error: "ExpectationFailed",
				Message:  "Error when handling the uploaded file",
			} }
			output.Messages = append(output.Messages, fileResponse)
			continue;
		}

		status, fileResponse := UploadFile(file, fileHeaders[index])
		output.Messages = append(output.Messages, FileResponse{ Status: status, Message: fileResponse })

		if status != http.StatusOK {
			output.Status = status
			output.Message = fileResponse.Message
		} else {
			output.Attachments = append(output.Attachments, fileResponse.Attachment)
		}
	}

	return output.Status, map[string]interface{}{
		"message": output.Message,
		"attachments": output.Attachments,
		"messages": output.Messages,
	}
}

func ServeFile(name string) (*string, *[]byte, error) {
	serverSettings := tools.GetSettings().Server
	attachment, err := models.DBAttachment.FindAttachment(name)
	if err != nil || attachment == nil {
		return nil, nil, fmt.Errorf("404 -> Attachment not found.")
	}

	path := fmt.Sprintf("%s/%s", serverSettings.UploadsPath, attachment.Url)
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("404 -> Attachment not found.")
	}

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	// read file into bytes
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	return &attachment.Type, &bytes, nil
}

func DeleteAttachment(attachmentId uint32) (int, map[string]interface{}) {
	rows, err := models.DBAttachment.DeleteAttachment(attachmentId)
	if err != nil || rows <= 0 {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error": "Unknown",
			"message": "Error deleting the attachment",
		}
	}

	return http.StatusAccepted, map[string]interface{}{
		"message": fmt.Sprintf("Attachment %d removed", attachmentId),
	}
}
