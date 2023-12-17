package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"jira_notifier/models"
	"math/rand"
	"net/http"
)

const (
	UserKey = "AUTH"
)

func LoginAction(c *gin.Context) {
	var loginModel models.Login

	if err := c.BindJSON(&loginModel); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := models.FindUserByUsername(loginModel.Username)

	//Not found or not admin
	if err != nil || user.Role != models.AdminRole {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username",
		})
		return
	}

	code := models.UserCode{
		Code:   rand.Intn(9999-1000) + 1000,
		UserID: user.ID,
	}

	models.DB.Create(&code)

	//telegram.SendTelegramPlainMessage(
	//	fmt.Sprintf("Привет, %v! "+
	//		"В целях вашей же безопастности никому не сообщайте содержимое этого письма. \n\n "+
	//		"Ваш код %v", user.Name, code.Code),
	//	user.ChatID)

	c.JSON(http.StatusCreated, gin.H{
		"name":      user.Name,
		"expiredAt": code.ExpiredAt,
	})

}

func CodeAction(c *gin.Context) {

	var code models.LoginCode

	if err := c.BindJSON(&code); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	authorizationCode, err := models.FindValidCode(code.Code)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid code",
		})
		return
	}

	user, err := models.FindUserById(authorizationCode.UserID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid code",
		})
		return
	}

	session := sessions.Default(c)
	session.Set(UserKey, user.Username)
	if err := session.Save(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Error when save session. Try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
	})

}

func ProfileAction(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
	})
}

func LogoutAction(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(UserKey)
	err := session.Save()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Error save session!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
