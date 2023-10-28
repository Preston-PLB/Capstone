package controllers

import (
	"fmt"
	"time"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/templates"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginPostBody struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func SignUpHandler (c *gin.Context) {
	//get uname and password.
	conf := config.Config()
	reqBody := &LoginPostBody{}
	c.Request.ParseForm()
	reqBody.Email = c.Request.FormValue("email")
	reqBody.Password = c.Request.FormValue("password")

	if reqBody.Email == "" {
		renderTempl(c, templates.SignupPage("Please provide an email"))
		return
	}

	if reqBody.Password == "" {
		renderTempl(c, templates.SignupPage("Please provide a password"))
		return
	}

	//Verify username and password
	user, err := mongo.FindUserByEmail(reqBody.Email)
	if err != nil {
		renderTempl(c, templates.SignupPage("Error occured. Please try again later"))
		return
	}

	if user != nil {
		renderTempl(c, templates.SignupPage(fmt.Sprintf("user already exists for %s", reqBody.Email)))
		return
	}

	user = &models.User{}

	passHash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)
	if err != nil {
		renderTempl(c, templates.SignupPage("Signup failed. Please try again later"))
		return
	}

	user.PassowrdHash = string(passHash)
	user.Email = reqBody.Email

	err = mongo.SaveModel(user)
	if err != nil {
		renderTempl(c, templates.SignupPage("Signup failed. Please try again later"))
		return
	}

	//build jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": user.UserId,
			"exp": time.Now().Add(12 * time.Hour).Unix(),
		},
	)

	jwtStr, err := token.SignedString(conf.JwtSecret)
	if err != nil {
		renderTempl(c, templates.SignupPage("Signup failed. Please try again later"))
		return
	}

	//store jwt as cookie
	//TODO: Make sure set secure for prd deployment
	c.SetCookie("authorization", jwtStr, 3600 * 24, "", "", false, true)

	c.Redirect(302, "/dashboard")
}

func LoginHandler(c *gin.Context) {
	//get uname and password.
	conf := config.Config()
	reqBody := &LoginPostBody{}
	c.Request.ParseForm()
	reqBody.Email = c.Request.FormValue("email")
	reqBody.Password = c.Request.FormValue("password")

	if reqBody.Email == "" {
		renderTempl(c, templates.LoginPage("Please provide an email"))
		return
	}

	if reqBody.Password == "" {
		renderTempl(c, templates.LoginPage("Please provide a password"))
		return
	}

	//Verify username and password
	user, err := mongo.FindUserByEmail(reqBody.Email)
	if err != nil {
		renderTempl(c, templates.LoginPage(err.Error()))
		return
	}

	if user == nil {
		renderTempl(c, templates.LoginPage(fmt.Sprintf("No user found for %s", reqBody.Email)))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassowrdHash), []byte(reqBody.Password)); err != nil {
		renderTempl(c, templates.LoginPage("Email and password are incorrect"))
		return
	}

	//build jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": user.UserId,
			"exp": time.Now().Add(12 * time.Hour).Unix(),
		},
	)

	jwtStr, err := token.SignedString(conf.JwtSecret)
	if err != nil {
		renderTempl(c, templates.LoginPage("An error occured. Please try again later"))
	}

	//store jwt as cookie
	var secure bool
	if conf.Env == "dev" {
		secure = false
	} else {
		secure = true
	}
	c.SetCookie("authorization", jwtStr, 3600 * 24, "", "", secure, true)

	c.Redirect(302, "/dashboard")
}

func LogoutHandler(c *gin.Context) {
	conf := config.Config()

	var secure bool
	if conf.Env == "dev" {
		secure = false
	} else {
		secure = true
	}
	c.SetCookie("authorization", "", 3600 * 24, "", "", secure, true)

	c.Redirect(302, "/login")
}
