package controllers

import (
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var mongo *db.DB
var log *logrus.Logger

func BuildRouter(r *gin.Engine) {
	conf := config.Config()

	var err error
	mongo, err = db.NewClient(conf.Mongo.Uri)
	if err != nil {
		panic(err)
	}

	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	r.Use(cors.Default())
	//mainpage
	r.GET("/", AuthMiddleware(false), LandingPage)

	//Auth
	r.GET("/login", AuthMiddleware(false), LoginPage)
	r.GET("/signup", AuthMiddleware(false), SignUpPage)
	r.POST("/login", LoginHandler)
	r.POST("/signup", SignUpHandler)
	r.POST("/logout", LogoutHandler)

	dashboard := r.Group("/dashboard")
	dashboard.Use(AuthMiddleware(true))
	dashboard.GET("", DashboardPage)

	//Vendor stuff
	vendor := r.Group("/vendor")
	vendor.Use(AuthMiddleware(true))

	youtube := vendor.Group("/youtube")
	youtube.POST("/initiate", InitiateYoutubeOuath)
	youtube.GET("/callback", ReceiveYoutubeOauth)

	pco := vendor.Group("/pco")
	pco.POST("/initiate", InitiatePCOOuath)
}
