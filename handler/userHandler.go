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

var (
	appJSON = "application/json"
)

func UserRegistration(c *gin.Context) {
	db := infra.GetDB()
	contentType := utils.GetContentType(c)
	_, _ = db, contentType
	user := entity.User{}

	if contentType == appJSON {
		c.ShouldBindJSON(&user)
	} else {
		c.ShouldBind(&user)
	}

	err := db.Debug().Create(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"age":      user.Age,
	})
}

func UserLogin(c *gin.Context) {
	db := infra.GetDB()
	contentType := utils.GetContentType(c)
	_, _ = db, contentType

	user := entity.User{}
	password := ""

	if contentType == appJSON {
		c.ShouldBindJSON(&user)
	} else {
		c.ShouldBind(&user)
	}

	password = user.Password
	err := db.Debug().Where("email = ?", user.Email).Take(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email/password",
		})
		return
	}

	comparePass := utils.ComparePass([]byte(user.Password), []byte(password))
	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email/password",
		})
		return
	}

	token := utils.GenerateToken(user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func UpdateUser(c *gin.Context) {
	db := infra.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := utils.GetContentType(c)
	user := entity.User{}

	userId, _ := strconv.Atoi(c.Param("userId"))
	UserID := uint(userData["id"].(float64))

	if userId != int(UserID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "You are not allowed to access this data",
		})
		return
	}

	if contentType == appJSON {
		c.ShouldBindJSON(&user)
	} else {
		c.ShouldBind(&user)
	}

	user.ID = UserID
	user.ID = uint(userId)
	err := db.Model(&user).Where("id = ?", userId).Updates(entity.User{
		Email:    user.Email,
		Username: user.Username,
		Age:      user.Age,
	}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"age":        user.Age,
		"updated_at": user.UpdatedAt,
	})
}

func DeleteUser(c *gin.Context) {
	db := infra.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := utils.GetContentType(c)
	_ = contentType
	User := entity.User{}

	userId, _ := strconv.Atoi(c.Param("userId"))
	UserID := uint(userData["id"].(float64))

	// jika token id pada userdata tidak sama dengan param
	if userId != int(UserID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "you are not allowed to access this data",
		})
		return
	}

	err := db.Debug().First(&User, userId).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Data not found",
			"message": "Data doesn't exist",
		})
		return
	}

	err = db.Debug().Delete(&User).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Delete failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your account has been successfully deleted",
	})
}
