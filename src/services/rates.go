package services

import (
	"backend-admin-proyect/src/db"
	log "github.com/sirupsen/logrus"
)

type RateService struct {
	db db.WithIndex[Rate]
}

func (r *RateService) AddRate(courseId, userId string, rate int) error {
	rta := Rate{
		CourseId:  courseId,
		Score:     rate,
		UserEmail: userId,
	}
	if _, err := r.db.GetBoth(rta.GetPrimaryKey(), rta.GetSecondaryKey()); err != nil {
		r.db.Insert(rta)
	} else {
		r.db.Update(rta)
	}
	return nil
}

func (r *RateService) GetRating(courseId string) []Rate {
	response, err := r.db.GetSecondary(courseId)
	if err != nil {
		log.Errorf("error while getting course: %s", err.Error())
	}
	return response
}

func CreateRateService(db db.WithIndex[Rate]) *RateService {
	return &RateService{
		db,
	}
}
