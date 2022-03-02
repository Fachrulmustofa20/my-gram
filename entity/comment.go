package entity

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comment struct {
	GormModel
	UserId  uint   `json:"user_id" form:"user_id"`
	PhotoId uint   `json:"photo_id" form:"photo_id"`
	Message string `json:"message" form:"message" valid:"required~Message of your comment is required"`
	User    *User
	Photo   *Photo
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(c)
	if errCreate != nil {
		err = errCreate
		return
	}
	err = nil
	return
}

func (c *Comment) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(c)
	if errCreate != nil {
		err = errCreate
		return
	}
	err = nil
	return
}
