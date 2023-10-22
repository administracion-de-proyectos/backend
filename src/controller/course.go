package controller

import (
	"backend-admin-proyect/src/middleware"
	"backend-admin-proyect/src/services"
	"backend-admin-proyect/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Course struct {
	service services.CourseService
	tv      middleware.TokenValidator[UserRequest]
	ss      services.SubscriptionService
}

func transform(course services.CourseState) CourseState {
	return CourseState{
		CreatorEmail:     course.CreatorEmail,
		CourseTitle:      course.CourseTitle,
		Classes:          course.Classes,
		Category:         course.Category,
		Metadata:         course.Metadata,
		AgeFiltered:      course.AgeFiltered,
		MinAge:           course.MinAge,
		MaxAge:           course.MaxAge,
		IsSchoolOriented: course.IsSchoolOriented,
	}
}

// CreateCourse godoc
//
//		@Summary		Create course
//		@Description	Create course using the token as a way to add account to course owner
//		@Tags			Course request
//		@Accept			json
//		@Produce		json
//		@Param			course	body		CourseRequest	true	"Title and Category are required"
//	    @Param          Authorization   header string      true "token required for request"
//		@Success		200		{object}	CourseState
//		@Failure		400		{object}	ErrorMsg
//		@Router			/course/ [post]
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
	if err := utils.FailIfZeroValue([]string{cr.Title, cr.Category}); err != nil {
		c.JSON(400, gin.H{
			"reason": "one of the required fields of title is empty",
		})
		return
	}
	if (cr.MaxAge != nil && cr.MinAge == nil) || (cr.MinAge != nil && cr.MaxAge == nil) || (cr.MinAge != nil && cr.MaxAge != nil && (*cr.MinAge == 0 || *cr.MaxAge == 0)) {
		c.JSON(400, gin.H{
			"reason": "if defined min or max age, then the other should be defined to (and also 0 is invalid)",
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
		CreatorEmail:     tokenData.Email,
		CourseTitle:      cr.Title,
		Classes:          classToCreate,
		Category:         cr.Category,
		Metadata:         cr.Metadata,
		IsSchoolOriented: cr.IsSchoolOriented,
	}
	if cr.MinAge != nil {
		course.MaxAge = *cr.MaxAge
		course.MinAge = *cr.MinAge
		course.AgeFiltered = true
	}
	courseCreated := ce.service.AddCourse(course)

	c.JSON(200, transform(courseCreated))
}

// GetCourse godoc
//
//		@Summary		Fetch a course
//		@Description	Fetch a course with a given id
//		@Tags			Course request
//		@Accept			json
//		@Produce		json
//		@Param			id      path		string	true	"course id which you look for"
//	    @Param          Authorization   header string      true "token required for request"
//		@Success		200		{object}	CourseState
//		@Failure		404		{object}	ErrorMsg
//		@Router			/course/{id} [get]
func (ce Course) GetCourse(c *gin.Context) {
	tokenData, err := ce.tv.GetTokenData(c)
	if err != nil {
		c.JSON(401, gin.H{
			"reason": "invalid token",
		})
		return
	}
	courseId := c.Param("id")
	if cs, err := ce.getCourse(courseId); err != nil {
		c.JSON(404, gin.H{
			"reason": "course not found",
		})
	} else {
		_, err = ce.ss.GetSubscription(tokenData.Email, courseId)
		cs.IsSubscribed = err == nil
		c.JSON(200, cs)
	}
}

func (ce Course) getCourse(courseId string) (CourseState, error) {
	var cs CourseState
	if course, err := ce.service.GetCourse(courseId); err != nil {
		log.Errorf("error while getting course: %s", err.Error())
		return cs, err
	} else {
		cs = transform(course)
		return cs, err
	}
}

// AddClass godoc
//
//	@Summary		Create class for created course
//	@Description	Create class for a previously created course, if course does not exist this endpoint will fail
//	@Tags			Course request
//	@Accept			json
//	@Produce		json
//	@Param			class	body		Class	true	"Title is required"
//	@Param			id	path		string	true	"course id which you want to add a course"
//	@Success		200		{object}	CourseState
//	@Failure        400     {object}    ErrorMsg
//	@Router			/course/{id} [post]
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
		c.JSON(200, transform(course))
	}
}

// RemoveClass godoc
//
//	@Summary		Remove class created
//	@Description	Removes a class
//	@Tags			Course request
//	@Accept			json
//	@Produce		json
//	@Param			classId	path		string	true	"class id you want to remove"
//	@Param			id	path		string	true	"course id which you look for"
//	@Success		204
//	@Failure        400     {object}    ErrorMsg
//	@Router			/course/{id}/{classId} [delete]
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

// GetClass godoc
//
//	@Summary		Fetch a class
//	@Description	Get class with id and class id
//	@Tags			Course request
//	@Accept			json
//	@Produce		json
//	@Param			classId	path		string	true	"class id you want to fetch"
//	@Param			id	path		string	true	"course id which you look for"
//	@Success		200		{object}	Class
//	@Failure        400     {object}    ErrorMsg
//	@Failure        404     {object}    ErrorMsg
//	@Router			/course/{id}/{classId} [get]
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

// GetCourses godoc
//
//	@Summary		Get all courses
//	@Description	Get all courses that follows a criteria
//	@Tags			Course request
//	@Accept			json
//	@Produce		json
//	@Param title query string false "Title string for which you want to look"
//	@Param ownerEmail query string false "ownerEmail string for which you want to look"
//	@Param category query string false "category string for which you want to look"
//	@Param desiredAge query int false "Age of the course you want to retrieve"
//	@Param isSchoolOriented query bool false "true if school oriented, any other value otherwise"
//	@Success		200		{object}	CourseStateResponse
//	@Failure        404     {object}    ErrorMsg
//	@Router			/course/ [get]
func (ce Course) GetCourses(c *gin.Context) {
	title := c.Query("title")
	ownerEmail := c.Query("ownerEmail")
	category := c.Query("category")
	age := c.Query("desiredAge")
	isSchoolOrientedS := c.Query("isSchoolOriented")
	var numberAge *int = nil
	if age != "" {
		if iage, err := strconv.Atoi(age); err != nil {
			c.JSON(404, gin.H{
				"reason": "age should be a number",
			})
			return
		} else {
			numberAge = &iage
		}
	}
	var isSchoolOriented *bool = nil
	if isSchoolOrientedS != "" {
		b := isSchoolOrientedS == "true"
		isSchoolOriented = &b
	}
	v := services.FilterValues{
		Title:            title,
		OwnerEmail:       ownerEmail,
		Category:         category,
		DesiredAge:       numberAge,
		IsSchoolOriented: isSchoolOriented,
	}
	courses := ce.service.GetCourses(v)
	c.JSON(200, gin.H{
		"courses": courses,
		"amount":  len(courses),
	})
}

// Subscribe godoc
//
//		@Summary		Subscribe
//		@Description	Subscribe a user given by its token to a course
//		@Tags			Subscription
//		@Accept			json
//		@Produce		json
//	    @Param          id   path string      true "course in which the current user wants to subscribe"
//	    @Param          Authorization   header string      true "token required for request"
//		@Success		200		{Object}	SubscriptionRequest
//		@Failure        401     {object}    ErrorMsg
//		@Failure        404     {object}    ErrorMsg
//		@Router			/course/subscribe/{id} [post]
func (ce Course) Subscribe(c *gin.Context) {
	courseId := c.Param("id")
	tokenData, err := ce.tv.GetTokenData(c)
	if err != nil {
		c.JSON(401, gin.H{
			"reason": "invalid token",
		})
		return
	}
	s := ce.ss.Subscribe(tokenData.Email, courseId)
	sr := SubscriptionRequest{
		UserId:      s.UserId,
		CourseTitle: s.CourseId,
		Metadata:    s.Metadata,
	}
	c.JSON(200, sr)
}

// GetSubscribed godoc
//
//		@Summary		Get subscribed courses
//		@Description	Get all courses in which the user has subscribed
//		@Tags			Subscription
//		@Accept			json
//		@Produce		json
//	    @Param          Authorization   header string      true "token required for request"
//		@Success		200		{object}	CourseStateResponse
//		@Failure        401     {object}    ErrorMsg
//		@Failure        404     {object}    ErrorMsg
//		@Router			/course/subscribe/ [get]
func (ce Course) GetSubscribed(c *gin.Context) {
	tokenData, err := ce.tv.GetTokenData(c)
	if err != nil {
		c.JSON(401, gin.H{
			"reason": "invalid token",
		})
		return
	}
	s := ce.ss.GetAllUserSubscriptions(tokenData.Email)
	courses := make([]CourseState, 0)
	for _, course := range s {
		if state, err := ce.getCourse(course.CourseId); err == nil {
			state.IsSubscribed = true
			courses = append(courses, state)
		}
	}
	c.JSON(200, gin.H{
		"courses": courses,
		"amount":  len(courses),
	})
}

func CreateControllerCourse(s services.CourseService, validator middleware.TokenValidator[UserRequest], ss services.SubscriptionService) Course {
	return Course{
		service: s,
		tv:      validator,
		ss:      ss,
	}
}
