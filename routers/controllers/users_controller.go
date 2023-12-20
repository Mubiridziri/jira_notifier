package controllers

import (
	"github.com/gin-gonic/gin"
	"jira_notifier/models"
	"net/http"
	"strconv"
)

func GetUsersList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	users, err := models.FindPaginatedUsers(page, limit)
	count, err := models.GetUsersTotal()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	apiUsers := models.ConvertUserToAPIUser(users)

	c.JSON(http.StatusOK, gin.H{
		"total":   count,
		"entries": apiUsers,
	})

}

func GrantAdminRole(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	user, err := models.FindUserById(uint(id))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
	}

	user.Role = models.AdminRole
	models.DB.Save(&user)

	apiUser := user.ConvertAPIUser()
	c.JSON(http.StatusOK, apiUser)

}

func RevokeAdminRole(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	user, err := models.FindUserById(uint(id))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.Role = models.UserRole
	models.DB.Save(&user)

	apiUser := user.ConvertAPIUser()
	c.JSON(http.StatusOK, apiUser)

}

func RemoveUser(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	user, err := models.FindUserById(uint(id))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	models.DB.Delete(&user)

	apiUser := user.ConvertAPIUser()
	c.JSON(http.StatusOK, apiUser)
}
