package models

import (
	"time"

	"github.com/wayn3h0/go-uuid"
	"github.com/jinzhu/gorm"
	"github.com/YagoCarballo/kumquat-academy-api/database"
)

type SessionModel struct {}
var DBSession SessionModel

func (model SessionModel) DB () *gorm.DB {
	return database.DB
}

func (model SessionModel) Create(userId uint32, deviceId string) (*Session, error) {
	// tries to find an active session for this User & Device Id
	dbSession, findError := model.findSessionForUserOnDevice(userId, deviceId)
	if findError != nil {
		return nil, findError
	} else if dbSession != nil {
		return dbSession, nil
	}

	// Generates an Access Token
	accessToken, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	// Creates the Session
	session := Session{
		Token:		accessToken.String(),
		UserID:		userId,
		DeviceID:	deviceId,
		ExpiresIn:	time.Now().AddDate(0, 0, 7),
		CreatedOn:	time.Now(),
	}

	// Create the Session
	query := model.DB().Create(session)
	if query.Error != nil {
		return nil, query.Error
	}

	return &session, nil
}

func (model SessionModel) FindSession(accessToken string) (*Session, error) {
	// Creates empty Session
	var session Session

	// Query the Session
	query := model.DB().Find(&session, "token = ?", accessToken)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	// Returns the Session
	return &session, nil
}

func (model SessionModel) findSessionForUserOnDevice(userId uint32, deviceId string) (*Session, error) {
	// Creates empty Session
	var session Session

	// Query the Session
	query := model.DB().Find(&session, "user_id = ? and device_id = ?", userId, deviceId)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	// Returns the Session
	return &session, nil
}

func (model SessionModel) RemoveSession(accessToken string) (int64, error) {
	// Query the Session
	query := model.DB().Where("token = ?", accessToken).Delete(Session{})
	if query.Error != nil {
		return 0, query.Error
	}

	// Returns the Session
	return query.RowsAffected, nil
}
