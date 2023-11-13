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
	FindUser(fv FilterValuesUser) []UserState
}

type ExamsService interface {
	Create(exam Exam)
	DoExam(result StudentExam) error
	GetScoreForExam(userEmail, courseId, classId string) (Score, error)
	GetScoreForExams(userEmail, courseId string) ([]Score, error)
	RemoveExam(courseId, classId string)
	GetExam(courseId, classId string) (Exam, error)
	GetAllScoreForExams(courseI string) ([]Score, error)
}

type CommentService interface {
	AddComment(courseId string, comment string, userId string) (Comments, error)
	GetComments(courseId string) (Comments, error)
}

type Group interface {
	AddToGroup(teacherEmail, studentEmail string) error
	GetGroup(teacherEmail string) []string
}

type RateInterface interface {
	AddRate(courseId, userId string, rate int) error
	GetRating(courseId string) []Rate
}
