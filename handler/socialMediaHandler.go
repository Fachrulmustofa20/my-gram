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

func CreateSocialMedia(c *gin.Context) {
	db := infra.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := utils.GetContentType(c)

	SocialMedia := entity.SocialMedia{}
	UserId := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserId = UserId
	err := db.Debug().Create(&SocialMedia).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":               SocialMedia.ID,
		"name":             SocialMedia.Name,
		"social_media_url": SocialMedia.SocialMediaUrl,
		"user_id":          SocialMedia.UserId,
		"created_at":       SocialMedia.CreatedAt,
	})
}

func GetSocialMedia(c *gin.Context) {
	db := infra.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := utils.GetContentType(c)
	UserId := uint(userData["id"].(float64))
	_, _ = UserId, contentType

	socialMedia := []entity.SocialMedia{}
	err := db.Debug().Order("id asc").Preload("User").Where("user_id = ?", UserId).Find(&socialMedia).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error,
		})
		return
	}

	response := []map[string]interface{}{}
	for _, val := range socialMedia {
		response = append(response, map[string]interface{}{
			"id":               val.ID,
			"name":             val.Name,
			"social_media_url": val.SocialMediaUrl,
			"UserId":           val.UserId,
			"createdAt":        val.CreatedAt,
			"updatedAt":        val.UpdatedAt,
			"User": map[string]interface{}{
				"id":       val.User.ID,
				"username": val.User.Username,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"social_medias": response,
	})
}

func UpdateSocialMedia(c *gin.Context) {
	db := infra.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := utils.GetContentType(c)
	SocialMedia := entity.SocialMedia{}

	socialMediaId, _ := strconv.Atoi(c.Param("socialMediaId"))
	UserID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserId = UserID
	SocialMedia.ID = uint(socialMediaId)

	err := db.Model(&SocialMedia).Where("id = ?", socialMediaId).Updates(entity.SocialMedia{
		Name:           SocialMedia.Name,
		SocialMediaUrl: SocialMedia.SocialMediaUrl,
	}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":               SocialMedia.ID,
		"name":             SocialMedia.Name,
		"social_media_url": SocialMedia.SocialMediaUrl,
		"user_id":          SocialMedia.UserId,
		"updated_at":       SocialMedia.UpdatedAt,
	})
}

func DeleteSocialMedia(c *gin.Context) {
	db := infra.GetDB()
	contentType := utils.GetContentType(c)
	_ = contentType
	SocialMedia := entity.SocialMedia{}

	socialMediaId, _ := strconv.Atoi(c.Param("socialMediaId"))

	err := db.Debug().First(&SocialMedia, socialMediaId).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Data Not Found",
			"message": "Data doesn't exist",
		})
		return
	}

	err = db.Debug().Delete(&SocialMedia).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Delete Failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})
}
