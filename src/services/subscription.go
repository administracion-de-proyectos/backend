package services

import (
	"backend-admin-proyect/src/db"
	log "github.com/sirupsen/logrus"
)

type subscriptionService struct {
	db db.WithIndex[Subscription]
}

func (s *subscriptionService) Subscribe(userId, courseId string) Subscription {
	sub := Subscription{
		UserId:   userId,
		CourseId: courseId,
		Metadata: nil,
	}
	s.db.Insert(sub)
	return sub
}

func (s *subscriptionService) GetAllUserSubscriptions(userId string) []Subscription {
	subs, err := s.db.GetPrimary(userId)
	if err != nil {
		log.Errorf("some weird error happened: %s", err.Error())
	}
	return subs
}

func (s *subscriptionService) GetAllCoursesSubscriptions(courseId string) []Subscription {
	subs, err := s.db.GetSecondary(courseId)
	if err != nil {
		log.Errorf("some weird error happened: %s", err.Error())
	}
	return subs
}

func (s *subscriptionService) GetSubscription(userId, courseId string) (Subscription, error) {
	return s.db.GetBoth(userId, courseId)
}

func (s *subscriptionService) RemoveSubscription(userId, courseId string) (Subscription, error) {
	return s.db.DeleteSpecific(userId, courseId), nil
}

func CreateSubscriptionService(db db.WithIndex[Subscription]) SubscriptionService {
	return &subscriptionService{db: db}
}
