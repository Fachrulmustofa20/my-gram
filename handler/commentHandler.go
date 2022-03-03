package handler

import (
	"mygram/entity"
	"mygram/infra"
	"mygram/utils"
	"net/http"
	"strconv"

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
			"message": "Data photo doesn't exist",
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

func GetAllComment(c *gin.Context) {
	db := infra.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := utils.GetContentType(c)
	UserId := uint(userData["id"].(float64))
	_ = contentType

	Comment := []entity.Comment{}
	err := db.Debug().Order("id asc").Preload("User").Preload("Photo").Where("user_id = ?", UserId).Find(&Comment).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error,
		})
		return
	}

	response := []map[string]interface{}{}
	for _, val := range Comment {
		response = append(response, map[string]interface{}{
			"id":         val.ID,
			"message":    val.Message,
			"photo_id":   val.PhotoId,
			"user_id":    val.UserId,
			"updated_at": val.UpdatedAt,
			"created_at": val.CreatedAt,
			"User": map[string]interface{}{
				"id":       val.User.ID,
				"email":    val.User.Email,
				"username": val.User.Username,
			},
			"Photo": map[string]interface{}{
				"id":        val.Photo.ID,
				"title":     val.Photo.Title,
				"caption":   val.Photo.Caption,
				"photo_url": val.Photo.PhotoUrl,
				"user_id":   val.Photo.UserId,
			},
		})
	}

	c.JSON(http.StatusOK, response)

}

func UpdateComment(c *gin.Context) {
	db := infra.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := utils.GetContentType(c)
	Comment := entity.Comment{}

	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userId := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserId = userId
	Comment.ID = uint(commentId)

	err := db.Model(&Comment).Where("id = ?", commentId).Updates(entity.Comment{
		Message: Comment.Message,
	}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         Comment.ID,
		"message":    Comment.Message,
		"user_id":    Comment.UserId,
		"photo_id":   Comment.PhotoId,
		"updated_at": Comment.UpdatedAt,
	})
}

func DeleteComment(c *gin.Context) {
	db := infra.GetDB()
	contentType := utils.GetContentType(c)
	_ = contentType
	Comment := entity.Comment{}

	commentId, _ := strconv.Atoi(c.Param("commentId"))

	// check comment
	err := db.Debug().First(&Comment, commentId).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Data Not Found",
			"message": "Data doesn't exist",
		})
		return
	}

	err = db.Debug().Delete(&Comment).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Delete Failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}
