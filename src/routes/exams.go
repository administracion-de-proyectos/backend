package routes

import (
	"backend-admin-proyect/src/controller"
	"backend-admin-proyect/src/db"
	"backend-admin-proyect/src/middleware"
	"backend-admin-proyect/src/services"
	"os"
)

func (r Routes) AddExamsRoutes(dbUrl string) error {
	exams, err := db.CreateDB[services.Exam]("examsTable", dbUrl)
	if err != nil {
		return err
	}
	submissionsDb, err := db.CreateDBWithIndex[services.StudentExam]("submissionTable", dbUrl)
	if err != nil {
		return err
	}
	validator := middleware.CreateValidator[controller.UserRequest](os.Getenv(secretValidator))
	c := controller.CreateControllerExams(services.CreateExamsService(exams, submissionsDb), validator)
	examGroup := r.Router.Group("/exams")
	examGroup.POST("/:courseId/:classId", c.CreateExam)
	examGroup.GET("/:courseId/:classId", c.GetExam)
	examGroup.DELETE("/:courseId/:classId", c.RemoveExam)
	examGroup.POST("/submission", validator.SetTokenDataInContext, c.CreateSubmission)
	subsGroup := r.Router.Group("/scores")
	subsGroup.GET("/:courseId/class/:classId/:userEmail", c.GetScore)
	subsGroup.GET("/:courseId/class/:classId", validator.SetTokenDataInContext, c.GetScoreAuth)
	subsGroup.GET("/:courseId/user/:userEmail", c.GetScores)
	subsGroup.GET("{courseId}/teacher", c.GetScoresTeacher)
	return nil
}
