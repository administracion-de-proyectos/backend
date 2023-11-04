package services

import (
	"backend-admin-proyect/src/db"
	log "github.com/sirupsen/logrus"
)

func (p Point) isCorrect(answers []Point) bool {
	for _, r := range answers {
		if p.Question == r.Question {
			return p.Answer == r.Answer
		}
	}
	return false
}

type ExamService struct {
	submissionDB db.WithIndex[StudentExam]
	examDb       db.DB[Exam]
}

func (e *ExamService) Create(exam Exam) {
	e.examDb.Insert(exam)
}

func (e *ExamService) DoExam(result StudentExam) error {
	e.submissionDB.Insert(result)
	return nil
}

func (e *ExamService) GetScoreForExam(userEmail, courseId, classId string) (Score, error) {
	submission, err := e.submissionDB.GetBoth(userEmail, getClassId(courseId, classId))
	if err != nil {
		return Score{}, err
	}
	return e.getScore(submission), nil
}

func (e *ExamService) GetAllScoreForExams(courseId string) ([]Score, error) {
	submissions, err := e.submissionDB.GetAll()
	if err != nil {
		return []Score{}, err
	}
	myCourseSubmissions := make([]Score, 0)
	for _, sub := range submissions {
		if sub.Course == courseId {
			myCourseSubmissions = append(myCourseSubmissions, e.getScore(sub))
		}
	}
	return myCourseSubmissions, nil
}

func (e *ExamService) getScore(submission StudentExam) Score {
	answers, err := e.examDb.Get(submission.GetSecondaryKey())
	if err != nil {
		log.Errorf("something fuckery has happened 2: %s", err.Error())
	}
	correct := 0
	for _, p := range submission.Points {
		if p.isCorrect(answers.Points) {
			correct += 1
		}
	}
	return Score{
		TotalAmount:   len(answers.Points),
		CorrectAmount: correct,
		Email:         submission.StudentEmail,
	}
}

func (e *ExamService) GetScoreForExams(userEmail, courseId string) ([]Score, error) {
	allSubmissions, err := e.submissionDB.GetPrimary(userEmail)
	if err != nil {
		return nil, err
	}
	scores := make([]Score, 0)
	for _, s := range allSubmissions {
		if s.Course == courseId {
			scores = append(scores, e.getScore(s))
		}
	}
	return scores, nil
}

func (e *ExamService) RemoveExam(courseId, classId string) {
	e.examDb.Delete(getClassId(courseId, classId))
}

func (e *ExamService) GetExam(courseId, classId string) (Exam, error) {
	return e.examDb.Get(getClassId(courseId, classId))
}

func CreateExamsService(examDb db.DB[Exam], submission db.WithIndex[StudentExam]) *ExamService {
	return &ExamService{
		examDb:       examDb,
		submissionDB: submission,
	}
}
