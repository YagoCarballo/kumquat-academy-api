package models

import (
	"github.com/jinzhu/gorm"
	"github.com/YagoCarballo/kumquat-academy-api/database"
)

type AttachmentsModel struct{}
var DBAttachment AttachmentsModel

func (model AttachmentsModel) DB() *gorm.DB {
	return database.DB
}

func (model AttachmentsModel) CreateAttachment(name, mimeType, token string) (*Attachment, error) {
	attachment := Attachment{
		Name: name,
		Type: mimeType,
		Url: token,
	}

	query := model.DB().Create(&attachment)
	if query.Error != nil {
		return nil, query.Error
	}

	return &attachment, nil
}

func (model AttachmentsModel) ReadAttachment(id uint32) (*Attachment, error) {
	var attachment Attachment

	query := model.DB().First(&attachment, "id = ?", id)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	return &attachment, nil
}

func (model AttachmentsModel) DeleteAttachment(id uint32) (int64, error) {
	query := model.DB().
		Table("attachments").
		Where("id = ?", id).
		Delete(Attachment{})
	if query.Error != nil {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}

func (model AttachmentsModel) FindAttachment(name string) (*Attachment, error) {
	var attachment Attachment

	query := model.DB().First(&attachment, "url = ?", name)
	if query.Error != nil {
		// If no Records found, return NIL otherwise return the error
		switch query.Error {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, query.Error
		}
	}

	return &attachment, nil
}
