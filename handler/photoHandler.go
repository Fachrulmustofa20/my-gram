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

	c.JSON(http.StatusCreated, Photo)
}

func GetPhotos(c *gin.Context) {
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

	c.JSON(http.StatusOK, Photo)
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
