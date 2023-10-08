package controller

import (
	"backend-admin-proyect/src/middleware"
	"backend-admin-proyect/src/services"
	"backend-admin-proyect/src/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserController struct {
	service   services.UserService
	validator middleware.TokenValidator[UserRequest]
}

func (uc UserController) CreateUser(c *gin.Context) {
	var ur UserRequest
	if err := c.BindJSON(&ur); err != nil {
		c.JSON(400, gin.H{
			"reason": err.Error(),
		})
		return
	}
	if err := utils.FailIfZeroValue([]string{ur.Profile, ur.Name, ur.Email, ur.Password}); err != nil {
		c.JSON(400, gin.H{
			"reason": "one of the required fields of email, profile, password or name is empty",
		})
		return
	}
	if len(ur.Password) < 6 {
		c.JSON(400, gin.H{
			"reason": "password too short",
		})
		return
	}
	if err := uc.service.CreateUser(services.CreateUserState(ur.Email, ur.Password, ur.Name, ur.Profile)); err != nil {
		log.Errorf("could not write into db")
		c.JSON(400, gin.H{
			"reason": err.Error(),
		})
		return
	}
	ur.Password = ""
	s, err := uc.validator.CreateToken(ur)
	if err != nil {
		log.Errorf("some fuckery happened, err: %s", err.Error())
	}
	c.JSON(200, gin.H{
		"token": s,
	})
}

func (uc UserController) SignInUser(c *gin.Context) {
	var ur UserRequest
	if err := c.BindJSON(&ur); err != nil {
		c.JSON(400, gin.H{
			"reason": err.Error(),
		})
		return
	}
	if err := utils.FailIfZeroValue([]string{ur.Email, ur.Password}); err != nil {
		c.JSON(400, gin.H{
			"reason": "one of the required fields of email or password is empty",
		})
		return
	}
	if state, err := uc.service.CheckCredentials(services.CreateBasicUserState(ur.Email, ur.Password)); err != nil {
		c.JSON(400, gin.H{
			"reason": err.Error(),
		})
		return
	} else {
		ur.Profile = state.Profile
		ur.Password = ""
		ur.Name = state.Name
		s, err := uc.validator.CreateToken(ur)
		if err != nil {
			log.Errorf("some fuckery happened, err: %s", err.Error())
		}
		c.JSON(200, gin.H{
			"token": s,
		})
	}
}

func CreateUserController(s services.UserService, validator middleware.TokenValidator[UserRequest]) UserController {
	return UserController{service: s, validator: validator}
}
