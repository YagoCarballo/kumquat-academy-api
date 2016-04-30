package models

import (
	"io"
	"fmt"
	"time"
	"regexp"
	"strings"
	"encoding/json"

	"github.com/albrow/forms"
	"github.com/jinzhu/gorm"
	"github.com/wayn3h0/go-uuid"
	"github.com/YagoCarballo/kumquat.academy.api/database"
)

type (
	UserModel struct {}

	SignInUser struct {
		Username string
		Password string
		Device string
	}
)
var DBUser UserModel

func (model UserModel) DB () *gorm.DB {
	return database.DB
}

func (model UserModel) FindUser(username string) (*User, error) {
	// Creates empty User
	var user User

	// Query the User
	query := model.DB().Preload("Avatar").Find(&user, "users.username = ?", username)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.RecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	// Returns the User
	return &user, nil
}

func (model UserModel) FindUserWithId(userId uint32) (*User, error) {
	// Creates empty User
	var user User

	// Query the User
	query := model.DB().Preload("Avatar").Find(&user, "users.id = ?", userId)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.RecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	// Returns the User
	return &user, nil
}

func (model UserModel) FindUserWithEmail(email string) (*User, error) {
	var user User

	query := model.DB().Preload("Avatar").Find(&user, "users.email = ?", email)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.RecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	return &user, nil
}

func (data *SignInUser) ParseUser(body *io.ReadCloser) error {
	decoder := json.NewDecoder(*body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func (model UserModel) CreateUser(user *User) (int64, error) {
	// Create the User
	query := DBUser.DB().Create(user)

	if query.Error != nil {
		isDuplicated := strings.HasPrefix(query.Error.Error(), "Error 1062")
		if isDuplicated {
			return -1, nil
		}

		return 0, query.Error
	}

	return query.RowsAffected, nil
}

func (model UserModel) CleanResetPasswordTokens(userId *uint32) (int64, error) {
	queryText := "expires <= now()"

	if userId != nil {
		queryText = fmt.Sprintf("%s or user_id = %d", queryText, *userId)
	}

	query := model.DB().Table("reset_passwords").Where(queryText).Delete(ResetPassword{})
	if query.Error != nil {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}

func (model UserModel) AddResetPasswordToken(userId uint32) (*string, error) {
	// Clean expired and previous tokens
	DBUser.CleanResetPasswordTokens(&userId)

	// Generate Unique token
	token, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// Populate the structure
	stringToken := token.String()
	resetPassword := ResetPassword{
		Token: stringToken,
		UserID: userId,
		Expires: time.Now().AddDate(0, 0, 5),
	}

	// Insert into DB
	query := DBUser.DB().Create(resetPassword)
	if query.Error != nil {
		return nil, query.Error
	}

	return &stringToken, nil
}

func (model UserModel) ChangePassword(token, newPassword string) (bool, error) {
	resetPassword := ResetPassword{}
	query := model.DB().Where("token = ?", token).First(&resetPassword)
	if query.Error != nil {
		return false, query.Error
	}

	// Clean expired and previous tokens
	DBUser.CleanResetPasswordTokens(&resetPassword.UserID)

	if !resetPassword.Expires.After(time.Now()) {
		return false, fmt.Errorf("The provided token has expired.")
	}

	query = DBUser.DB().Table("users").Where("id = ?", resetPassword.UserID).Updates(map[string]interface{}{
		"password": newPassword,
	})
	if query.Error != nil {
		return false, query.Error
	}

	return true, nil
}

func (model UserModel) UpdateStudent(user *User) (*User, error) {
	query := DBUser.DB().Table("users").Where("id = ?", user.ID).Updates(map[string]interface{}{
		"first_name": user.FirstName,
		"last_name": user.LastName,
		"username": user.Username,
		"email": user.Email,
		"matric_date": user.MatricDate,
		"matric_number": user.MatricNumber,
		"date_of_birth": user.DateOfBirth,
		"avatar_id": user.AvatarId,
	})

	if query.Error != nil {
		isDuplicated := strings.HasPrefix(query.Error.Error(), "Error 1062")
		if isDuplicated {
			return nil, nil
		}

		return nil, query.Error
	}

	dbUser, err := DBUser.FindUserWithId(user.ID)
	if err != nil {
		return nil, fmt.Errorf("Error reading the updated user")
	}

	return dbUser, nil
}

func (model UserModel) AddAvatarToUser(userId uint32, attachmentId *uint32) (int64, error) {
	query := DBUser.DB().Table("users").Where("id = ?", userId).Updates(map[string]interface{}{
		"avatar_id": attachmentId,
	})

	if query.Error != nil {
		return query.RowsAffected, query.Error
	}

	return query.RowsAffected, nil
}

func (user *User) Save() (*User, error) {
	query := DBUser.DB().Where("id = ?", user.ID).Save(&user)
	if query.Error != nil {
		return nil, query.Error
	}

	return user, nil
}

func (model UserModel) SearchUsers(text, moduleCode string) ([]User, error) {
	// Creates empty User
	var users []User
	queryText := fmt.Sprint("%", text, "%")

	// Query the User
	query := model.DB().Table("users").Preload("Avatar").Limit(10).Select("distinct users.*").Find(
		&users,
		"(users.username like ? or users.email like ? or users.matric_number like ? or users.first_name like ? or " +
		"users.last_name like ? or users.id = ?) and " +
		"(users.id not in (select user_id from user_modules where module_code = ?))",
		queryText, queryText, queryText, queryText, queryText, text, moduleCode,
	)
	if query.Error != nil {
		return users, query.Error
	}

	// Returns the User
	return users, nil
}

func (model UserModel) ListOfAreasWithWritePermissionForUser(userId uint32) ([]string, error) {
	var userModules []UserModule
	var userCourses []UserCourse
	modulesWithAccess := []string{}

	query := model.DB().Table("user_modules").Joins(
		"inner join roles on roles.id = user_modules.role_id",
	).Where("user_modules.user_id = ? and roles.can_write = ?", userId, true).Find(&userModules)
	if query.Error != nil {
		return modulesWithAccess, query.Error
	}

	for _, userModule := range userModules {
		modulesWithAccess = append(modulesWithAccess, userModule.ModuleCode)
	}

	query = model.DB().Table("user_courses").Joins(
		"inner join roles on roles.id = user_courses.role_id",
	).Where("user_courses.user_id = ? and roles.can_write = ?", userId, true).Find(&userCourses)
	if query.Error != nil {
		return modulesWithAccess, query.Error
	}

	for _, userCourse := range userCourses {
		var levelModules []LevelModule
		query = model.DB().Table("level_modules").Joins(
			"inner join classes on classes.id = level_modules.class_id",
		).Where("classes.course_id = ?", userCourse.CourseID).Find(&levelModules)
		if query.Error != nil {
			continue
		}
		for _, levelModule := range levelModules {
			modulesWithAccess = append(modulesWithAccess, levelModule.Code)
		}
	}

	return modulesWithAccess, nil
}

func (user *User) Parse(userData *forms.Data) (error, map[string][]string) {
	var err error

	// Checks that the given Data is Valid
	validationErrors := isValid(userData)
	if validationErrors != nil {
		return nil, validationErrors
	}

	// Saves the User info into the pointer
	user.Username = userData.Get("username")
	user.Password = userData.Get("password")
	user.Email = userData.Get("email")
	user.FirstName = userData.Get("first_name")
	user.LastName = userData.Get("last_name")
	user.MatricNumber = userData.Get("matric_number")
	user.Active = userData.GetBool("active")

	if userData.KeyExists("avatar_id") {
		user.AvatarId = uint32(userData.GetInt("avatar_id"))
	}

	// Parses the Dates into Go Time
	var dateOfBirth, matricDate time.Time

	// Parses the Date of Birth
	dateOfBirth, err = time.Parse("02-01-2006", userData.Get("date_of_birth"))
	if err != nil {
		return nil, map[string][]string{
			"date_of_birth": []string{
				"Invalid date format. [dd-MM-yyyy]",
			},
		}
	}
	user.DateOfBirth = dateOfBirth

	// Parses the Matriculation Date
	if userData.KeyExists("matric_date") {
		matricDate, err = time.Parse("02-01-2006 15:04:05", userData.Get("matric_date"))
		if err != nil {
			return nil, map[string][]string{
				"matric_date": []string{
					"Invalid date format. [dd-MM-yyyy hh:mm:ss]",
				},
			}
		}
		user.MatricDate = matricDate
	}

	return nil, nil
}

func isValid(userData *forms.Data) map[string][]string {
	// Validate
	val := userData.Validator()

	if userData.KeyExists("id") {
		val.TypeInt("id")
	}

	val.Require("username")
	val.LengthRange("username", 4, 20)

	val.Require("password")
	val.MinLength("password", 8)

	val.Require("email")
	val.MatchEmail("email")

	val.Require("first_name")
	val.LengthRange("first_name", 2, 30)

	val.Require("last_name")
	val.LengthRange("last_name", 2, 30)

	pattern := regexp.MustCompile("[0-9]{2}-[0-9]{2}-[0-9]{4}")
	val.Require("date_of_birth")
	val.Match("date_of_birth", pattern)

	val.Require("matric_number")
	val.LengthRange("matric_number", 4, 30)

	if userData.KeyExists("matric_date") {
		pattern = regexp.MustCompile("[0-9]{2}-[0-9]{2}-[0-9]{4}\\s[0-9]{2}:[0-9]{2}:[0-9]{2}")
		val.Match("matric_date", pattern)
	}

	if val.HasErrors() {
		return val.ErrorMap()
	}

	return nil
}
