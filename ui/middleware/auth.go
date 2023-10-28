package middleware

import (
	//"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(strict bool) gin.HandlerFunc {
	//conf := config.Config();
	return func(c *gin.Context) {
		//check for cookie

		/*
			no cookie:
			if strict
				route login
			else:
				continue with nil auth_user
		*/

		/*
			cookie:
			validate jwt:
			if valid:
				set context with auth_user
			else:
				if strict
					route login
				else:
					continue with nil auth_user

		*/
	}
}
