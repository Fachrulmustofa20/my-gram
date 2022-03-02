package handler

import (
	"mygram/entity"
	"mygram/infra"
	"mygram/utils"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	db := infra.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := utils.GetContentType(c)

	Comment := entity.Comment{}
	UserId := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	err := db.Debug().First(&entity.Photo{}, Comment.PhotoId).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Data Not Found",
			"message": "Data doesn't exist",
		})
		return
	}

	Comment.UserId = UserId

	err = db.Debug().Create(&Comment).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Comment.ID,
		"message":    Comment.Message,
		"photo_id":   Comment.PhotoId,
		"user_id":    Comment.UserId,
		"created_at": Comment.CreatedAt,
	})

}
