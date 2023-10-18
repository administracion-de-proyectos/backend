package services

import (
	"backend-admin-proyect/src/db"
	"backend-admin-proyect/src/utils"
	"strings"
)

type classService struct {
	classDB  db.DB[Class]
	courseDB db.DB[CourseState]
}

func (c *classService) SetClassInPlaceN(n int, classTitle string, courseTitle string) CourseState {
	var err error
	var course CourseState
	if course, err = c.GetCourse(classTitle); err != nil || !utils.Contains(course.Classes, courseTitle) {
		return course
	}
	finalClasses := make([]string, 0)
	inserted := false
	for i, class := range course.Classes {
		if i == n {
			finalClasses = append(finalClasses, classTitle)
			inserted = true
		}
		if class != classTitle {
			finalClasses = append(finalClasses, classTitle)
		}
	}
	if !inserted {
		finalClasses = append(finalClasses, classTitle)
	}
	course.Classes = finalClasses
	c.courseDB.Update(course)
	return course
}

func (c *classService) AddClass(class Class, shouldEditCourse bool) (CourseState, error) {
	var err error
	var course CourseState
	if shouldEditCourse {
		if course, err = c.courseDB.Get(class.CourseTitle); err != nil {
			return course, err
		}
		course.Classes = append(course.Classes, class.Id)
		c.courseDB.Update(course)
	}
	c.classDB.Insert(class)
	return course, err
}

func (c *classService) AddCourse(course CourseState) CourseState {
	c.courseDB.Insert(course)
	return course
}

func (c *classService) GetCourse(courseId string) (CourseState, error) {
	return c.courseDB.Get(courseId)
}

func (c *classService) RemoveClass(courseId, classId string) error {
	var err error
	var course CourseState
	if course, err = c.GetCourse(courseId); err != nil {
		return err
	}
	finalClasses := make([]string, 0)
	for _, class := range course.Classes {
		if class != classId {
			finalClasses = append(finalClasses, class)
		}
	}
	c.classDB.Delete(getClassId(courseId, classId))
	course.Classes = finalClasses
	c.courseDB.Update(course)
	return nil
}

func (c *classService) GetClass(courseId, classId string) (Class, error) {
	return c.classDB.Get(getClassId(courseId, classId))
}

func (c *classService) GetCourses(title, ownerEmail string) []CourseState {
	courses, _ := c.courseDB.GetAll()
	filtered := make([]CourseState, 0)
	for _, course := range courses {
		if strings.Contains(course.CourseTitle, title) && strings.Contains(course.CreatorEmail, ownerEmail) {
			filtered = append(filtered, course)
		}
	}
	return filtered
}

func CreateCourseService(courseDb db.DB[CourseState], classDb db.DB[Class]) *classService {
	return &classService{
		courseDB: courseDb,
		classDB:  classDb,
	}
}
