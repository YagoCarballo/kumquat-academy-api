package endpoints

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"net/textproto"

	"golang.org/x/crypto/bcrypt"

	"github.com/YagoCarballo/kumquat.academy.api/database/models"
	"github.com/YagoCarballo/kumquat.academy.api/tools"

	emailHandler "github.com/jordan-wright/email"
)

//
// Sign In
// Checks the given credentials and returns an access token
//
// @param username (string) the username
// @param password (string) the password (hashed)
// @param deviceId (string) the a generated ID for the device being used
// @return status (int) The Response Status
// @return message (map[string]) The response object
//
func SignIn(username, password, deviceId string) (int, map[string]interface{}, *tools.JWTSession) {
	var jwtSession tools.JWTSession

	// Gets the User from the Database
	dbUser, userError := models.DBUser.FindUser(username)
	if userError != nil || dbUser == nil {
		return http.StatusForbidden, map[string]interface{}{
			"error":   "InvalidCredentials",
			"message": "The provided credentials are invalid.",
		}, &jwtSession
	}

	// Checks if the password is correct
	bcriptError := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if bcriptError != nil {
		return http.StatusForbidden, map[string]interface{}{
			"error":   "InvalidPassword",
			"message": "The provided password does not match.",
		}, &jwtSession
	}

	// Generates a Session
	session, sessionErr := models.DBSession.Create(dbUser.ID, deviceId)
	if sessionErr != nil {
		log.Println(sessionErr)
		return http.StatusExpectationFailed, map[string]interface{}{
			"error":   "SessionError",
			"message": "Error creating the session",
		}, &jwtSession
	}

	// Creates a struct with session info
	jwtSession = tools.JWTSession{
		AccessToken: 	session.Token,
		UserId: 		session.UserID,
		ExpiresIn: 		session.ExpiresIn,
		Device: 		session.DeviceID,
		Username: 		dbUser.Username,
		Admin:			dbUser.Admin,
	}

	modulesWithAccess, _ := models.DBUser.ListOfAreasWithWritePermissionForUser(dbUser.ID)
	outputUser := map[string]interface{}{
		"id": 				dbUser.ID,
		"first_name":		dbUser.FirstName,
		"last_name":		dbUser.LastName,
		"username":			dbUser.Username,
		"email":			dbUser.Email,
		"matric_number":	dbUser.MatricNumber,
		"matric_date":		dbUser.MatricDate,
		"date_of_birth":	dbUser.DateOfBirth,
		"admin":			dbUser.Admin,
		"avatar_id":		dbUser.AvatarId,
		"avatar":			nil,
		"write_access":		modulesWithAccess,
	}

	if dbUser.Avatar != nil {
		outputUser["avatar"] = dbUser.Avatar.Url
	}

	// Returns the access token
	return http.StatusAccepted,
		map[string]interface{}{
			"message": "Successfully logged in.",
			"user": outputUser,
		},
		&jwtSession
}

//
// Sign Up
// Registers a new User.
//
// @param user (*models.User) the user
// @return status (int) The Response Status
// @return message (map[string]) The response object
//
func SignUp(user *models.User) (int, map[string]interface{}) {
	if user == nil || user.Username == "" {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error":   "InvalidUser",
			"message": "The provided User is Invalid.",
		}
	}

	// Salts and Hashes the password
	salt, bcriptError := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if bcriptError != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error":   "InvalidPassword",
			"message": "Password rejected.",
		}
	}

	// Updates the password with the salt
	user.Password = string(salt)

	// Forces admin flag to be disabled for new users
	user.Admin = false

	// Inserts the User into the Database
	result, registerError := models.DBUser.CreateUser(user)
	if registerError != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error":   "SignUpError",
			"message": "Error registering the user",
		}
	}

	if result <= 0 {
		return http.StatusConflict, map[string]interface{}{
			"error":   "DuplicatedUser",
			"message": "A user with that username or email already exists.",
		}
	}

	// Returns the success message
	return http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully, an email will arrive soon with instructions.",
	}
}

//
// UserInfo
// Gets information about a user.
//
// @param userId (uint32) the user id
// @return status (int) The Response Status
// @return message (map[string]) The response object
//
func UserInfo(userId uint32) (int, map[string]interface{}) {
	// Gets the User from the Database
	dbUser, userError := models.DBUser.FindUserWithId(userId)
	if userError != nil || dbUser == nil {
		return http.StatusForbidden, map[string]interface{}{
			"error":   "InvalidCredentials",
			"message": "The provided credentials3 are invalid.",
		}
	}

	modulesWithAccess, _ := models.DBUser.ListOfAreasWithWritePermissionForUser(userId)
	outputUser := map[string]interface{}{
		"id": 				dbUser.ID,
		"first_name":		dbUser.FirstName,
		"last_name":		dbUser.LastName,
		"username":			dbUser.Username,
		"email":			dbUser.Email,
		"matric_number":	dbUser.MatricNumber,
		"matric_date":		dbUser.MatricDate,
		"date_of_birth":	dbUser.DateOfBirth,
		"admin":			dbUser.Admin,
		"avatar_id":		dbUser.AvatarId,
		"avatar":			nil,
		"write_access":		modulesWithAccess,
	}

	if dbUser.Avatar != nil {
		outputUser["avatar"] = dbUser.Avatar.Url
	}

	// Returns the access token
	return http.StatusOK, map[string]interface{}{
		"user": outputUser,
	}
}

//
// Log Out
// Gets information about a user.
//
// @param userId (uint32) the user id
// @return status (int) The Response Status
// @return message (map[string]) The response object
//
func LogOut(accessToken string) (int, map[string]interface{}) {
	// Gets the User from the Database
	rowsAffected, err := models.DBSession.RemoveSession(accessToken)
	if err != nil || rowsAffected == 0 {
		return http.StatusForbidden, map[string]interface{}{
			"error":   "InvalidCredentials",
			"message": "The provided credentials are invalid.",
		}
	}

	// Returns the access token
	return http.StatusAccepted,
	map[string]interface{}{
		"message": "Successfully logged out",
	}
}

func ChangePassword(token, newPassword string) (int, map[string]interface{}) {
	// Salts and Hashes the password
	salt, bcriptError := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if bcriptError != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error":   "InvalidPassword",
			"message": "Password rejected.",
		}
	}

	changed, err := models.DBUser.ChangePassword(token, string(salt))
	if err != nil || !changed {
		return http.StatusForbidden, map[string]interface{}{
			"error":   "InvalidToken",
			"message": "The provided token is invalid.",
		}
	}

	return http.StatusAccepted, map[string]interface{}{
		"message": "Successfully changed password",
	}
}

func ForgotPassword(email string) (int, map[string]interface{}) {
	emailSettings := tools.GetSettings().Email

	user, err := models.DBUser.FindUserWithEmail(email)
	if err != nil || user == nil {
		return http.StatusNotFound, map[string]interface{}{
			"error":   "UserNotFound",
			"message": "A user with that email does not exist.",
		}
	}

	token, err := models.DBUser.AddResetPasswordToken(user.ID)
	if err != nil || token == nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error":   "UnknownError",
			"message": "Error generation reset token.",
		}
	}

	resetUrl := fmt.Sprintf("%s%s", "https://kumquat.academy/password/reset/", *token)

	plainText := fmt.Sprintf(`
		A password reset has been triggered for your account. If you haven’t triggered this please ignore this email.\n
		To reset your password open the following link and follow the instructions:\n\n
		%s\n\n
		Thanks,\n
		Kumquat Academy Team\n
	`, resetUrl);

	htmlText := fmt.Sprintf(`
    	<div style="word-wrap: break-word; -webkit-nbsp-mode: space; -webkit-line-break: after-white-space;">
    		A password reset has been triggered for your account. If you haven’t triggered this please ignore this email.
        	<br />
        	<div>To reset your password open the following link and follow the instructions:</div>
       		<br />
			<div>
				<a href="%s">
					<span style="color: rgb(0, 0, 0);">Click Here to reset your password</span>
				</a>
			</div>
			<br />
			<br />
			<p>Thanks,</p>
			<p>Kumquat Academy Team</p>
		</div>
	`, resetUrl);


	emailTemplate := &emailHandler.Email {
		To: []string{ user.Email },
		From: emailSettings.Sender,
		Subject: "Password Reset",
		Text: []byte(plainText),
		HTML: []byte(htmlText),
		Headers: textproto.MIMEHeader{},
	}

	smtpServer := fmt.Sprintf("%s:%d", emailSettings.Server, emailSettings.Port)
	err = emailTemplate.Send(
		smtpServer,
		smtp.PlainAuth("", emailSettings.User, emailSettings.Password, emailSettings.Server),
	)

	if err != nil {
		return http.StatusExpectationFailed, map[string]interface{}{
			"error":   "EmailError",
			"message": "Error sending the email with the instructions.",
		}
	}

	return http.StatusOK, map[string]interface{}{
		"message": "Email with instructions to reset the password sent.",
	}
}