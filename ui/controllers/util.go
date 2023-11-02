package controllers

import (
	"bytes"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

// Responds with 200ok and the rendered template
func renderTempl(c *gin.Context, tmpl templ.Component) {
	buf := bytes.NewBuffer([]byte{})
	tmpl.Render(c.Request.Context(), buf)

	c.Data(200, "text/html", buf.Bytes())
}

func badRequest(c *gin.Context, reason string) {
	c.JSON(400, map[string]string{"error": reason})
}

func serverError(c *gin.Context, reason string) {
	c.JSON(504, map[string]string{"error": reason})
}

func notFound(c *gin.Context, reason string) {
	c.JSON(404, map[string]string{"error": reason})
}

func unauthorized(c *gin.Context, reason string) {
	c.JSON(403, map[string]string{"error": reason})
}
