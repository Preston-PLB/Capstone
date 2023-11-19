package main

import (
	"fmt"
	"os"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	r := gin.Default()

	controllers.BuildRouter(r)

	var addr string
	if port := os.Getenv("PORT"); port != "" {
		addr = fmt.Sprintf("0.0.0.0:%s", port)
	} else {
		addr = "0.0.0.0:8008"
	}

	err := r.Run(addr)
	if err != nil {
		panic(err)
	}
}
