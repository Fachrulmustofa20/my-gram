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

func CreatePhoto(c *gin.Context) {
	db := infra.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := utils.GetContentType(c)

	Photo := entity.Photo{}
	userId := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	Photo.UserId = userId
	err := db.Debug().Create(&Photo).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Photo.ID,
		"title":      Photo.Title,
		"caption":    Photo.Caption,
		"photo_url":  Photo.PhotoUrl,
		"user_id":    Photo.UserId,
		"created_at": Photo.CreatedAt,
	})
}

func GetAllPhotos(c *gin.Context) {
	db := infra.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := utils.GetContentType(c)
	UserId := uint(userData["id"].(float64))
	_, _ = UserId, contentType

	photos := []entity.Photo{}
	err := db.Debug().Order("id asc").Preload("User").Find(&photos).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error,
		})
		return
	}

	response := []map[string]interface{}{}
	for _, val := range photos {
		response = append(response, map[string]interface{}{
			"id":         val.ID,
			"title":      val.Title,
			"caption":    val.Caption,
			"photo_url":  val.PhotoUrl,
			"user_id":    val.UserId,
			"created_at": val.CreatedAt,
			"updated_at": val.UpdatedAt,
			"User": map[string]interface{}{
				"email":    val.User.Email,
				"username": val.User.Username,
			},
		})
	}

	c.JSON(http.StatusOK, response)
}

func UpdatePhoto(c *gin.Context) {
	db := infra.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := utils.GetContentType(c)
	Photo := entity.Photo{}

	photoId, _ := strconv.Atoi(c.Param("photoId"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	Photo.UserId = userID
	Photo.ID = uint(photoId)

	err := db.Model(&Photo).Where("id = ?", photoId).Updates(entity.Photo{
		Title:    Photo.Title,
		Caption:  Photo.Caption,
		PhotoUrl: Photo.PhotoUrl,
	}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         Photo.ID,
		"title":      Photo.Title,
		"caption":    Photo.Caption,
		"photo_url":  Photo.PhotoUrl,
		"user_id":    Photo.UserId,
		"updated_at": Photo.UpdatedAt,
	})
}

func DeletePhoto(c *gin.Context) {
	db := infra.GetDB()
	contentType := utils.GetContentType(c)
	_ = contentType
	Photo := entity.Photo{}

	photoId, _ := strconv.Atoi(c.Param("photoId"))

	err := db.Debug().First(&Photo, photoId).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Data Not Found",
			"message": "Data doesn't exist",
		})
		return
	}

	err = db.Debug().Delete(&Photo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Delete Failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})
}
