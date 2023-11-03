package controllers

import (
	"net/http"
	"time"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const USER_OBJ_KEY = "userObj"

type AuthClaims struct {
	Subject   string `json:"sub"`
	Expires   int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	NotBefore int64  `json:"nbf"`
	Issuer    string `json:"iss"`
	Audience  string `json:"aud"`
}

func (claims AuthClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	time := time.Unix(claims.Expires, 0)
	return jwt.NewNumericDate(time), nil
}

func (claims AuthClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	time := time.Unix(claims.IssuedAt, 0)
	return jwt.NewNumericDate(time), nil
}

func (claims AuthClaims) GetNotBefore() (*jwt.NumericDate, error) {
	time := time.Unix(claims.NotBefore, 0)
	return jwt.NewNumericDate(time), nil
}

func (claims AuthClaims) GetIssuer() (string, error) {
	return claims.Issuer, nil
}

func (claims AuthClaims) GetSubject() (string, error) {
	return claims.Subject, nil
}

func (claims AuthClaims) GetAudience() (jwt.ClaimStrings, error) {
	return []string{claims.Subject}, nil
}

func AuthMiddleware(strict bool) gin.HandlerFunc {
	conf := config.Config()
	return func(c *gin.Context) {
		//check for cookie
		token, err := c.Cookie("authorization")
		if err != nil {
			if err == http.ErrNoCookie {
				if strict {
					c.Redirect(301, "/login")
					return
				} else {
					return
				}
			} else {
				log.WithError(err).Error("Unable to get cookie from browser")
				c.AbortWithError(504, err)
				return
			}
		}

		claims := &AuthClaims{}

		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
			return []byte(conf.JwtSecret), nil
		})
		if err != nil {
			if err == jwt.ErrTokenExpired {
				log.Warn("Redirecting, jwt expired")
				c.Redirect(301, "/login")
				return
			} else {
				if strict {
					log.Warnf("Redirecting, jwt issue: %s", err)
					c.Redirect(301, "/login")
					return
				} else {
					log.Warnf("Jwt is invalid, but auth is not strict. Reason: %s", err)
					return
				}
			}
		}

		if !parsedToken.Valid {
			if strict {
				log.Warn("Redirecting, jwt invalid")
				c.Redirect(301, "/login")
				return
			} else {
				log.Warn("Jwt is invalid, but auth is not strict")
				return
			}
		}

		user, err := mongo.FindUserById(claims.Subject)
		if err != nil {
			log.WithError(err).Errorf("Unable to get user: %s from DB", claims.Subject)
			c.AbortWithError(502, err)
		}

		if user == nil {
			log.Errorf("Unable to find user: %s in DB", claims.Subject)
			c.AbortWithError(502, nil)
		}

		//store user object reference in session.
		c.Set(USER_OBJ_KEY, user)
	}
}
