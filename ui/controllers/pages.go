package controllers

import (
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/templates"
	"github.com/gin-gonic/gin"
)

func LandingPage(c *gin.Context) {
	renderTempl(c, templates.LandingPage(false))
}

func LoginPage(c *gin.Context) {
	renderTempl(c, templates.LoginPage(""))
}

func SignUpPage(c *gin.Context) {
	renderTempl(c, templates.SignupPage(""))
}

func DashboardPage(c *gin.Context) {
	c.JSON(200, gin.H{"response": "dashboard"})
}
