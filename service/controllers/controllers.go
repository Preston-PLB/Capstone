package controllers

import (
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db"
	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/youtube/v3"
)

var (
	log          *logrus.Logger
	mongo        *db.DB
	ytClientMap  map[primitive.ObjectID]*youtube.Service
	pcoClientMap map[primitive.ObjectID]*pco.PcoApiClient
)

func BuildRouter(r *gin.Engine) {
	conf := config.Config()

	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	var err error
	mongo, err = db.NewClient(conf.Mongo.Uri)
	if err != nil {
		panic(err)
	}

	ytClientMap = make(map[primitive.ObjectID]*youtube.Service)
	pcoClientMap = make(map[primitive.ObjectID]*pco.PcoApiClient)

	pco := r.Group("/pco")
	pco.Use(ValidatePcoWebhook)
	pco.POST("/:userid", ConsumePcoWebhook)
}
