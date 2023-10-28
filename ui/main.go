package main

import (
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	controllers.BuildRouter(r)

	r.Run()
}
