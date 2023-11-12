package routes

import (
	"backend-admin-proyect/src/controller"
	"backend-admin-proyect/src/db"
	"backend-admin-proyect/src/middleware"
	"backend-admin-proyect/src/services"
	"os"
)

func (r Routes) AddCoursesRoutes(dbUrl string) error {
	courseDb, err := db.CreateDB[services.CourseState]("coursesTable", dbUrl)
	if err != nil {
		return err
	}
	classDb, err := db.CreateDB[services.Class]("classTable", dbUrl)
	if err != nil {
		return err
	}
	subscriptionDb, err := db.CreateDBWithIndex[services.Subscription]("subsTable", dbUrl)
	if err != nil {
		return err
	}
	commentsDb, err := db.CreateDB[services.Comments]("commentsTable", dbUrl)
	if err != nil {
		return err
	}
	groupsDb, err := db.CreateDB[services.GroupDto]("groupsTable", dbUrl)
	if err != nil {
		return err
	}
	rateDb, err := db.CreateDBWithIndex[services.Rate]("rateTable", dbUrl)
	if err != nil {
		return err
	}
	validator := middleware.CreateValidator[controller.UserRequest](os.Getenv(secretValidator))

	c := controller.CreateControllerCourse(services.CreateCourseService(courseDb, classDb), validator, services.CreateSubscriptionService(subscriptionDb), services.CreateCommentsService(commentsDb), services.CreateGroupsService(groupsDb), services.CreateRateService(rateDb))
	courseGroup := r.Router.Group("/course")
	courseGroup.POST("/", validator.SetTokenDataInContext, c.CreateCourse)
	courseGroup.POST("/:id", c.AddClass)
	courseGroup.GET("/:id", validator.SetTokenDataInContext, c.GetCourse)
	courseGroup.GET("/:id/:classId", c.GetClass)
	courseGroup.DELETE("/:id/:classId", c.RemoveClass)
	courseGroup.GET("/", c.GetCourses)
	r.Router.GET("/courses", validator.SetTokenDataInContext, c.GetOwnedCourses) // because I am fucking lazy
	subsGroup := courseGroup.Group("/subscribe")
	subsGroup.POST("/:id", validator.SetTokenDataInContext, c.Subscribe)
	subsGroup.GET("/", validator.SetTokenDataInContext, c.GetSubscribed)
	subsGroup.GET("/courses/:id", c.GetAllSubscribed)
	commentsGroup := courseGroup.Group("/comments")
	commentsGroup.GET("/:courseId", c.GetComment)
	commentsGroup.POST("/", validator.SetTokenDataInContext, c.AddComment)
	groupsGroup := courseGroup.Group("/group")
	groupsGroup.GET("/", validator.SetTokenDataInContext, c.GetGroup)
	groupsGroup.POST("/add/:userId", validator.SetTokenDataInContext, c.AddToGroup)
	rateGroup := courseGroup.Group("/rate")
	rateGroup.POST("/add", validator.SetTokenDataInContext, c.AddRate)
	rateGroup.GET("/:courseId", c.GetRate)
	return nil
}
