package controllers

import (
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/middleware"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/templates"
	"github.com/gin-gonic/gin"
)


func BuildRouter(r *gin.Engine) {
	r.GET("/", middleware.AuthMiddleware(false) ,LandingPage)
	r.GET("/login", middleware.AuthMiddleware(false), LoginPage)
}

func LandingPage(c *gin.Context) {
	renderTempl(c, templates.LandingPage(false))
}

func LoginPage(c *gin.Context) {
	renderTempl(c, templates.LoginPage())
}
