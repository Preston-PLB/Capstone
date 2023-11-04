package main

import (
	"log"
	"net/http"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	r := gin.Default()

	controllers.BuildRouter(r)

	err := http.ListenAndServeTLS(":8080", "tls.crt", "tls.key", r)
	if err != nil {
		log.Fatal(err)
	}
}
