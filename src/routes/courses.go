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
	validator := middleware.CreateValidator[controller.UserRequest](os.Getenv(secretValidator))
	c := controller.CreateControllerCourse(services.CreateCourseService(courseDb, classDb), validator)
	group := r.Router.Group("/course")
	group.POST("/", validator.SetTokenDataInContext, c.CreateCourse)
	group.POST("/:id", c.AddClass)
	group.GET("/:id", c.GetCourse)
	group.GET("/:id/:classId", c.GetClass)
	group.DELETE("/:id/:classId", c.RemoveClass)
	group.GET("/", c.GetCourses)
	return nil
}
