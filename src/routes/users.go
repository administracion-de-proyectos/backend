package routes

import (
	"backend-admin-proyect/src/controller"
	"backend-admin-proyect/src/db"
	"backend-admin-proyect/src/middleware"
	"backend-admin-proyect/src/services"
	"github.com/gin-gonic/gin"
	"os"
)

const (
	secretValidator = "SECRET_VALIDATOR"
)

func (r Routes) AddUserRoutes(dbUrl string) error {
	uDb, err := db.CreateDB[services.UserState]("userTable", dbUrl)
	if err != nil {
		return err
	}
	validator := middleware.CreateValidator[controller.UserRequest](os.Getenv(secretValidator))
	c := controller.CreateUserController(services.CreateUserService(uDb), validator)
	group := r.Router.Group("/user")
	group.POST("/login", c.SignInUser)
	group.POST("/signUp", c.CreateUser)
	group.GET("/validateToken", validateToken)
	group.GET("/profile/:id", c.GetUser)
	group.GET("/profile/", validator.SetTokenDataInContext, c.GetUserWithToken)
	group.PATCH("/profile/:id", c.UpdateUser)
	group.GET("/find", c.FindUsers)
	return nil
}

func validateToken(c *gin.Context) {
	token := c.GetHeader("authorization")
	validator := middleware.CreateValidator[controller.UserRequest](os.Getenv(secretValidator))
	if err := validator.ValidateToken(token); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	data, _ := validator.GetData(token)
	c.JSON(200, gin.H{
		"message": "pong",
		"data":    data,
	})
}
