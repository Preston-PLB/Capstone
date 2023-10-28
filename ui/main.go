package main

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"preston-baxter.com/capstone/frontend-service/templates"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		raw := []byte{}
		buf := bytes.NewBuffer(raw)
		templates.Hello("str").Render(c.Request.Context(), buf)


		c.Data(200, "text/html", []byte(buf.String()))
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
