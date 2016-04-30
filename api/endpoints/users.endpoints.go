package endpoints

import (
	"fmt"
	"time"
	"net/http"
	"net/smtp"
	"net/textproto"

	"github.com/YagoCarballo/kumquat.academy.api/tools"
	"github.com/YagoCarballo/kumquat.academy.api/database/models"

	emailHandler "github.com/jordan-wright/email"
	"mime/multipart"
)

func emailInstructionsToSetPassword(email, token string) (error) {
	emailSettings := tools.GetSettings().Email

	resetUrl := fmt.Sprintf("%s%s", "https://kumquat.academy/password/reset/", token)

	plainText := fmt.Sprintf(`
		Welcome to Kumquat Academy!!\n\n
		You are now registered as a Student, to access your learning follow the link and set your password.\n\n
		%s\n\n
		Thanks,\n
		Kumquat Academy Team\n
	`, resetUrl);

	htmlText := fmt.Sprintf(`
		<b>Welcome to Kumquat Academy!!</b>
		<br />
		<br />
		You are now registered as a Student, to access your learning follow the link and set your password.
		<br />
		<br />
		<p><a href="%s"><span style="color: rgb(0, 0, 0);">Click Here to set your password</span></a></p>
		<br />
		<p>Thanks,</p>
		<p>Kumquat Academy Team</p>
	`, resetUrl);


	emailTemplate := &emailHandler.Email {
		To: []string{ email },
		From: emailSettings.Sender,
		Subject: "Welcome to Kumquat Academy",
		Text: []byte(plainText),
		HTML: []byte(htmlText),
		Headers: textproto.MIMEHeader{},
	}

	smtpServer := fmt.Sprintf("%s:%d", emailSettings.Server, emailSettings.Port)
	err := emailTemplate.Send(
		smtpServer,
		smtp.PlainAuth("", emailSettings.User, emailSettings.Password, emailSettings.Server),
	)

	return err
}

func CreateStudentAndAddToModule(moduleCode, firstName, lastName, username, email, matricNumber string, matricDate, dateOfBirth time.Time, avatarId *uint32) (int, map[string]interface{}) {
	user := models.User{
		FirstName: firstName,
		LastName: lastName,
		Username: username,
		Email: email,
		MatricNumber: matricNumber,
		MatricDate: matricDate,
		DateOfBirth: dateOfBirth,
		Admin: false,
		Active: false,
	}

	if avatarId != nil {
		user.AvatarId = *avatarId
	}

	count, err := models.DBUser.CreateUser(&user)
	if err != nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Error creating student.",
		}
	}

	if count <= 0 {
		return http.StatusConflict, map[string]interface{}{
			"error": "Student already exists.",
		}
	}

	dbUser, err := models.DBUser.FindUser(username)
	if err != nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Error processing the student.",
		}
	}

	token, err := models.DBUser.AddResetPasswordToken(dbUser.ID)
	if err != nil || token == nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error":   "UnknownError",
			"message": "Error generating the password reset token.",
		}
	}

	err = emailInstructionsToSetPassword(email, *token)
	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error":   "EmailError",
			"message": "Error sending welcome email.",
		}
	}

	return AddStudentToModule(moduleCode, dbUser.ID)
}


func UpdateStudent(userId uint32, moduleCode, firstName, lastName, username, email, matricNumber string, matricDate, dateOfBirth time.Time, avatarId *uint32) (int, map[string]interface{}) {
	dbStudent, err := models.DBModule.GetModuleStudent(userId, moduleCode)
	if err != nil || dbStudent == nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Error no student with that ID was found on that module.",
		}
	}

	user := models.User{
		ID: userId,
		FirstName: firstName,
		LastName: lastName,
		Username: username,
		Email: email,
		MatricNumber: matricNumber,
		MatricDate: matricDate,
		DateOfBirth: dateOfBirth,
		Admin: false,
		Active: false,
	}

	if avatarId != nil {
		user.AvatarId = *avatarId
	}

	student, err := models.DBUser.UpdateStudent(&user)
	if err != nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Error updating the student.",
		}
	}

	if student == nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Student has duplicated data",
		}
	}

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
	}

	if student.Avatar != nil {
		studentMap["avatar"] = student.Avatar.Url
	}

	return http.StatusOK,  map[string]interface{}{
		"message": "Updated student successfully",
		"student": studentMap,
	}
}


func SearchUsers(query, moduleCode string) (int, map[string]interface{}) {
	students, err := models.DBUser.SearchUsers(query, moduleCode)
	if err != nil {
		return http.StatusConflict, map[string]interface{}{
			"error": "Error searching students.",
		}
	}

	studentMaps := []map[string]interface{}{}
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
		}

		if student.Avatar != nil {
			studentMap["avatar"] = student.Avatar.Url
		}

		studentMaps = append(studentMaps, studentMap)
	}

	return http.StatusOK,  map[string]interface{}{
		"students": studentMaps,
	}
}

func UploadAvatar(userId uint32, file multipart.File, header *multipart.FileHeader) (int, FileResponseMessage) {
	status, response := UploadFile(file, header)
	if status != http.StatusOK {
		return http.StatusExpectationFailed, FileResponseMessage{
			Error: "Unknown",
			Message: "Error uploading the file.",
		}
	}

	count, err := models.DBUser.AddAvatarToUser(userId, &response.Attachment.ID)
	if err != nil || count <= 0 {
		return http.StatusExpectationFailed, FileResponseMessage{
			Error: "Unknown",
			Message: "Error updating the avatar.",
		}
	}

	return status, response
}
