package controller

import (
	"backend-admin-proyect/src/middleware"
	"backend-admin-proyect/src/services"
	"backend-admin-proyect/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Course struct {
	service services.CourseService
	tv      middleware.TokenValidator[UserRequest]
}

func (ce Course) CreateCourse(c *gin.Context) {
	var tokenData UserRequest
	var err error
	var cr CourseRequest
	if tokenData, err = ce.tv.GetTokenData(c); err != nil {
		c.JSON(401, gin.H{
			"reason": "invalid token",
		})
		return
	}
	if err := c.BindJSON(&cr); err != nil {
		c.JSON(400, gin.H{
			"reason": err.Error(),
		})
		return
	}
	if err := utils.FailIfZeroValue([]string{cr.Title}); err != nil {
		c.JSON(400, gin.H{
			"reason": "one of the required fields of title is empty",
		})
		return
	}
	classToCreate := make([]string, 0)
	for _, class := range cr.Classes {
		cs := services.Class{
			Id:          class.Title,
			CourseTitle: cr.Title,
			Metadata:    class.Metadata,
		}
		ce.service.AddClass(cs, false)
		classToCreate = append(classToCreate, cs.Id)
	}
	course := services.CourseState{
		CreatorEmail: tokenData.Email,
		CourseTitle:  cr.Title,
		Classes:      classToCreate,
	}
	courseCreated := ce.service.AddCourse(course)

	c.JSON(200, courseCreated)
}

func (ce Course) GetCourse(c *gin.Context) {
	courseId := c.Param("id")
	if course, err := ce.service.GetCourse(courseId); err != nil {
		log.Errorf("error while getting course: %s", err.Error())
		c.JSON(404, gin.H{
			"reason": "course not found",
		})
	} else {
		c.JSON(200, course)
	}
}

func (ce Course) AddClass(c *gin.Context) {
	courseId := c.Param("id")
	var cr Class
	if err := c.BindJSON(&cr); err != nil {
		c.JSON(400, gin.H{
			"reason": err.Error(),
		})
		return
	}
	if err := utils.FailIfZeroValue([]string{cr.Title, courseId}); err != nil {
		c.JSON(400, gin.H{
			"reason": "one of the required fields of title or course at params is empty",
		})
		return
	}
	if course, err := ce.service.AddClass(services.Class{
		Id:          cr.Title,
		CourseTitle: courseId,
		Metadata:    cr.Metadata,
	}, true); err != nil {
		c.JSON(400, gin.H{
			"reason": fmt.Sprintf("could not create class, error: %s", err.Error()),
		})
	} else {
		c.JSON(200, course)
	}
}

func (ce Course) RemoveClass(c *gin.Context) {
	courseId := c.Param("id")
	classId := c.Param("classId")
	if err := utils.FailIfZeroValue([]string{classId, courseId}); err != nil {
		c.JSON(400, gin.H{
			"reason": "one of the required fields of course at params is empty",
		})
		return
	}
	ce.service.RemoveClass(courseId, classId)
	c.JSON(204, gin.H{})
}

func (ce Course) GetClass(c *gin.Context) {
	courseId := c.Param("id")
	classId := c.Param("classId")
	if err := utils.FailIfZeroValue([]string{classId, courseId}); err != nil {
		c.JSON(400, gin.H{
			"reason": "one of the required fields of course at params is empty",
		})
		return
	}
	if class, err := ce.service.GetClass(courseId, classId); err != nil {
		log.Errorf("error while getting class: %s", err.Error())
		c.JSON(404, gin.H{
			"reason": "class not found",
		})
		return
	} else {
		c.JSON(200, class)
	}
}

func (ce Course) GetCourses(c *gin.Context) {
	title := c.Query("title")
	ownerEmail := c.Query("ownerEmail")
	courses := ce.service.GetCourses(title, ownerEmail)
	c.JSON(200, gin.H{
		"courses": courses,
		"amount":  len(courses),
	})
}

func CreateControllerCourse(s services.CourseService, validator middleware.TokenValidator[UserRequest]) Course {
	return Course{
		service: s,
		tv:      validator,
	}
}
