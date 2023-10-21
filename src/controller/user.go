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

// CreateUser godoc
//
//	@Summary		Sign Up User
//	@Description	Create User Account
//	@Tags			User request
//	@Accept			json
//	@Produce		json
//	@Param			user	body		UserRequest	true	"User required Data to SignUp"
//	@Success		200		{object}	Token
//	@Failure		400		{object}	ErrorMsg
//	@Router			/user/signUp/ [post]
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
	c.JSON(200, Token{
		Token: s,
	})
}

// SignInUser godoc
//
//	@Summary		SignIn User
//	@Description	SignInUser
//	@Tags			User request
//	@Accept			json
//	@Produce		json
//	@Param			user	body		UserRequest	true	"Email and Password are required"
//	@Success		200		{object}	Token
//	@Failure		400		{object}	ErrorMsg
//	@Router			/user/login/ [post]
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
		c.JSON(200, Token{
			Token: s,
		})
	}
}

// GetUser godoc
//
//	@Summary		Get User Profile
//	@Description	Get User Profile
//	@Tags			User request
//	@Param			id	path	string	true	"User ID"
//	@Produce		json
//	@Success		200	{object}	UserResponse
//	@Failure		400	{object}	ErrorMsg
//	@Router			/user/profile/{id} [get]
func (uc UserController) GetUser(c *gin.Context) {
	userId := c.Param("id")
	uc.sendUserWithId(c, userId)
}

func (uc UserController) sendUserWithId(c *gin.Context, userId string) {
	if user, err := uc.service.GetUser(userId); err != nil {
		log.Errorf("error while getting user: %s", err.Error())
		c.JSON(404, gin.H{
			"reason": "user not found",
		})
	} else {
		c.JSON(200, UserResponse{
			Email:   user.Email,
			Name:    user.Name,
			Profile: user.Profile,
		})
	}
}
func (uc UserController) GetUserWithToken(c *gin.Context) {
	userTokenData, err := uc.validator.GetTokenData(c)
	if err != nil {
		log.Errorf("error while checking token: %s", err.Error())
		c.JSON(403, gin.H{
			"reason": "invalid token",
		})
		return
	}
	userId := userTokenData.Email
	uc.sendUserWithId(c, userId)
}

// UpdateUser godoc
//
//	@Summary		Update User Profile
//	@Description	Update User Profile
//	@Tags			User request
//	@Param			id		path	string		true	"User ID"
//	@Param			user	body	UserRequest	true	"Profile and Name are updatable"
//	@Produce		json
//	@Success		200	{object}	UserResponse
//	@Failure		400	{object}	ErrorMsg
//	@Router			/user/profile/{id} [patch]
func (uc UserController) UpdateUser(c *gin.Context) {
	userId := c.Param("id")
	var ur UserRequest
	if err := c.BindJSON(&ur); err != nil {
		c.JSON(400, gin.H{
			"reason": err.Error(),
		})
		return
	}
	if user, err := uc.service.GetUser(userId); err != nil {
		log.Errorf("error while getting user: %s", err.Error())
		c.JSON(404, gin.H{
			"reason": "user not found",
		})
		return
	} else {
		if ur.Profile == "" {
			log.Infof("No se intenta modificar profile")
		} else {
			user.Profile = ur.Profile
		}
		if ur.Name == "" {
			log.Infof("No se intenta modificar Name")
		} else {
			user.Name = ur.Name
		}
		if updatedUser, err := uc.service.UpdateUser(user); err != nil {
			log.Errorf("error while updating user: %s", err.Error())
			c.JSON(404, gin.H{
				"reason": err.Error(),
			})
			return
		} else {
			c.JSON(200, UserResponse{
				Email:   updatedUser.Email,
				Name:    updatedUser.Name,
				Profile: updatedUser.Profile,
			})
		}
	}
}

func CreateUserController(s services.UserService, validator middleware.TokenValidator[UserRequest]) UserController {
	return UserController{service: s, validator: validator}
}
