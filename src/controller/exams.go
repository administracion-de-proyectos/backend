package controller

import (
	"backend-admin-proyect/src/middleware"
	"backend-admin-proyect/src/services"
	"backend-admin-proyect/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ExamsController struct {
	es services.ExamsService
	tv middleware.TokenValidator[UserRequest]
}

// CreateExam godoc
//
//	@Summary		Create an exam for a given class
//	@Description	Given a course id and a class id, creates an exam for the class
//	@Tags			Exams request
//	@Accept			json
//	@Produce		json
//	@Param			course	body		CreateExamRequest	true	"At least one point is required"
//	@Param			classId	path		string	true	"class id you want to remove"
//	@Param			courseId	path		string	true	"course id which you look for"
//	@Success		204
//	@Failure        400     {object}    ErrorMsg
//	@Router			/exams/{courseId}/{classId} [post]
func (e ExamsController) CreateExam(c *gin.Context) {
	var cer CreateExamRequest
	courseId := c.Param("courseId")
	classId := c.Param("classId")
	if err := c.BindJSON(&cer); err != nil {
		c.JSON(400, gin.H{
			"reason": err.Error(),
		})
		return
	} else if cer.Points == nil || len(cer.Points) == 0 {
		c.JSON(400, gin.H{
			"reason": "invalid amount of points or null",
		})
		return
	}
	examPoints := make([]services.Point, 0)
	for _, p := range cer.Points {
		if !utils.Contains(p.Possibilities, p.Answer) {
			c.JSON(400, gin.H{
				"reason": fmt.Sprintf("point: %s does not have %s between the possible answers", p.Question, p.Answer),
			})
			return
		}
		examPoints = append(examPoints, services.Point{
			Question: p.Question,
			Options:  p.Possibilities,
			Answer:   p.Answer,
		})
	}
	exam := services.Exam{
		Points: examPoints,
		Class:  classId,
		Course: courseId,
	}
	e.es.Create(exam)
	c.JSON(204, nil)
}

// GetExam godoc
//
//	@Summary		Get an exam
//	@Description	Given a course id and a class id, gets the specific exam
//	@Tags			Exams request
//	@Accept			json
//	@Produce		json
//	@Param			classId	path		string	true	"class id you want to remove"
//	@Param			courseId	path		string	true	"course id which you look for"
//	@Success		200 {object} Exam
//	@Failure        404     {object}    ErrorMsg
//	@Router			/exams/{courseId}/{classId} [get]
func (e ExamsController) GetExam(c *gin.Context) {
	courseId := c.Param("courseId")
	classId := c.Param("classId")
	if exam, err := e.es.GetExam(courseId, classId); err != nil {
		c.JSON(400, gin.H{
			"reason": err.Error(),
		})
	} else {
		c.JSON(200, exam)
	}
}

// RemoveExam godoc
//
//	@Summary		Remove exam created
//	@Description	Removes an exam already created
//	@Tags			Exams request
//	@Accept			json
//	@Produce		json
//	@Param			classId	path		string	true	"class id you want to remove"
//	@Param			courseId	path		string	true	"course id which you look for"
//	@Success		204
//	@Failure        400     {object}    ErrorMsg
//	@Router			/exams/{courseId}/{classId} [delete]
func (e ExamsController) RemoveExam(c *gin.Context) {
	courseId := c.Param("id")
	classId := c.Param("classId")
	if err := utils.FailIfZeroValue([]string{classId, courseId}); err != nil {
		c.JSON(400, gin.H{
			"reason": "one of the required fields of course at params is empty",
		})
		return
	}
	e.es.RemoveExam(courseId, classId)
	c.JSON(204, gin.H{})
}

// CreateSubmission godoc
//
//	@Summary		Add submission
//	@Description	Given a user identified by its token, submit a resolution
//	@Tags			Exams request
//	@Accept			json
//	@Produce		json
//	@Param          Authorization   header string      true "token required for request"
//	@Param			course	body	Submission	true	"At least one point is required"
//	@Success		200 {object} Score
//	@Failure        400     {object}    ErrorMsg
//	@Router			/exams/submission [post]
func (e ExamsController) CreateSubmission(c *gin.Context) {
	var s Submission
	tokenData, err := e.tv.GetTokenData(c)
	if err != nil {
		c.JSON(401, gin.H{
			"reason": "invalid token",
		})
		return
	}
	if err := c.BindJSON(&s); err != nil {
		c.JSON(400, gin.H{
			"reason": err.Error(),
		})
		return
	} else if s.Points == nil {
		c.JSON(400, gin.H{
			"reason": "something on points must be provided",
		})
	}

	examPoints := make([]services.Point, 0)
	for _, p := range s.Points {
		examPoints = append(examPoints, services.Point{
			Question: p.Question,
			Answer:   p.Answer,
		})
	}
	if err := e.es.DoExam(services.StudentExam{
		StudentEmail: tokenData.Email,
		Course:       s.Course,
		Class:        s.Class,
		Points:       examPoints,
	}); err != nil {
		c.JSON(400, gin.H{
			"reason": err.Error(),
		})
	}
	if score, err := e.es.GetScoreForExam(tokenData.Email, s.Course, s.Class); err != nil {
		c.JSON(404, gin.H{
			"reason": err.Error(),
		})
	} else {
		c.JSON(200, score)
	}
}

// GetScore godoc
//
//	@Summary		Get a score
//	@Description	Given a course id and a class id and a user, gets the specific score
//	@Tags			Exams request
//	@Accept			json
//	@Produce		json
//	@Param			classId	path		string	true	"class id you want to look for"
//	@Param			courseId	path		string	true	"course id which you look for"
//	@Param			userEmail	path		string	true	"email you look for, is an exact match"
//	@Success		200 {object} Score
//	@Failure        400     {object}    ErrorMsg
//	@Router			/scores/{courseId}/class/{classId}/{userEmail} [get]
func (e ExamsController) GetScore(c *gin.Context) {
	courseId := c.Param("courseId")
	classId := c.Param("classId")
	userEmail := c.Param("userEmail")
	if scores, err := e.es.GetScoreForExam(courseId, classId, userEmail); err != nil {
		c.JSON(404, gin.H{
			"reason": "could not find user email",
		})
	} else {
		c.JSON(200, scores)
	}
}

// GetScoreAuth godoc
//
//	@Summary		Get a score
//	@Description	Given a course id and a class id and a user, gets the specific score
//	@Tags			Exams request
//	@Accept			json
//	@Produce		json
//	@Param			classId	path		string	true	"class id you want to look for"
//	@Param			courseId	path		string	true	"course id which you look for"
//	@Param          Authorization   header string      true "token required for request"
//	@Success		200 {object} Score
//	@Failure        400     {object}    ErrorMsg
//	@Router			/scores/{courseId}/class/{classId} [get]
func (e ExamsController) GetScoreAuth(c *gin.Context) {
	courseId := c.Param("courseId")
	classId := c.Param("classId")
	tokenData, err := e.tv.GetTokenData(c)
	if err != nil {
		c.JSON(401, gin.H{
			"reason": "invalid token",
		})
		return
	}
	if scores, err := e.es.GetScoreForExam(courseId, classId, tokenData.Email); err != nil {
		c.JSON(404, gin.H{
			"reason": "could not find user email",
		})
	} else {
		c.JSON(200, scores)
	}
}

// GetScores godoc
//
//	@Summary		Get an exam
//	@Description	Given a course id and a user email, gets all the scores from that user in given course
//	@Tags			Exams request
//	@Accept			json
//	@Produce		json
//	@Param			courseId	path		string	true	"course id which you look for"
//	@Param			userEmail	path		string	true	"email you look for, is an exact match"
//	@Success		200 {array} Score
//	@Failure        400     {object}    ErrorMsg
//	@Router			/scores/{courseId}/user/{userEmail} [get]
func (e ExamsController) GetScores(c *gin.Context) {
	courseId := c.Param("courseId")
	userEmail := c.Param("userEmail")
	if score, err := e.es.GetScoreForExams(userEmail, courseId); err != nil {
		c.JSON(404, gin.H{
			"reason": err.Error(),
		})
	} else {
		c.JSON(200, score)
	}
}

func CreateControllerExams(s services.ExamsService, validator middleware.TokenValidator[UserRequest]) ExamsController {
	return ExamsController{
		es: s,
		tv: validator,
	}
}
