package entity

import (
	"mygram/utils"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	GormModel
	Username    string        `gorm:"not null;uniqueIndex" json:"username" form:"username" valid:"required~Your username is required"`
	Email       string        `gorm:"not null;uniqueIndex" json:"email" form:"email" valid:"required~Your email is required,email~Invalid email format"`
	Password    string        `gorm:"not null" json:"password" form:"password" valid:"required~Your password is required,minstringlength(6)~Password has to have a minimum length of 6 characters"`
	Age         uint          `gorm:"not null" json:"age" form:"age" valid:"required~Your age is required,range(8|150)~Age cannot be under 8"`
	Photo       []Photo       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photo"`
	Comment     []Comment     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comment"`
	SocialMedia []SocialMedia `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"social_media"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = utils.HashPass(u.Password)
	err = nil
	return
}
