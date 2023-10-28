package services

type SubscriptionService interface {
	Subscribe(userId, courseId string) Subscription
	GetAllUserSubscriptions(userId string) []Subscription
	GetAllCoursesSubscriptions(courseId string) []Subscription
	GetSubscription(userId, courseId string) (Subscription, error)
	RemoveSubscription(userId, courseId string) (Subscription, error)
}

type CourseService interface {
	SetClassInPlaceN(n int, classTitle string, courseTitle string) CourseState
	AddClass(class Class, shouldEditCourse bool) (CourseState, error)
	AddCourse(course CourseState) CourseState
	GetCourse(courseId string) (CourseState, error)
	RemoveClass(courseId, classId string) error
	GetClass(courseId, classId string) (Class, error)
	GetCourses(values FilterValues) []CourseState
}

type UserService interface {
	CreateUser(u UserState) error
	CheckCredentials(u UserState) (UserState, error)
	GetUser(userId string) (UserState, error)
	UpdateUser(u UserState) (UserState, error)
}

type ExamsService interface {
	Create(exam Exam)
	DoExam(result StudentExam) error
	GetScoreForExam(userEmail, courseId, classId string) (Score, error)
	GetScoreForExams(userEmail, courseId string) ([]Score, error)
	RemoveExam(courseId, classId string)
	GetExam(courseId, classId string) (Exam, error)
}
