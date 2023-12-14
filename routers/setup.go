package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"jira_notifier/models"
	"jira_notifier/routers/controllers"
	"net/http"
)

// TODO utils.GenerateToken(length) instead hardcode
var secret = []byte("RHYaxoa6iqb1VTCsFtdM2PAAu8i8CYhU")

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions("AUTH", cookie.NewStore(secret)))

	v1 := r.Group("/api/v1")

	v1.POST("/security/login", controllers.LoginAction)
	v1.POST("/security/code", controllers.CodeAction)

	v1.Use(AuthRequired)
	{
		v1.GET("/profile", controllers.ProfileAction)
		v1.POST("/logout", controllers.LogoutAction)

	}

	return r
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	userKey := session.Get(controllers.UserKey)

	if userKey == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user, err := models.FindUserByUsername(userKey.(string))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid cookie",
		})
		return
	}

	c.Set("user", user)
	c.Next()
}
