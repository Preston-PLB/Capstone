package controllers

import (
	"bytes"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

//Responds with 200ok and the rendered template
func renderTempl(c *gin.Context, tmpl templ.Component) {
	buf := bytes.NewBuffer([]byte{})
	tmpl.Render(c.Request.Context(), buf)

	c.Data(200, "text/html", buf.Bytes())
}
