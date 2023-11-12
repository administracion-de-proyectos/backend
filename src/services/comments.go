package services

import (
	"backend-admin-proyect/src/db"
	"time"
)

type CommentsService struct {
	db db.DB[Comments]
}

func (c *CommentsService) AddComment(courseId string, comment string, userId string) (Comments, error) {
	newComment := Comment{
		CreatedAt:  time.Now().Unix(),
		UserId:     userId,
		Commentary: comment,
	}
	comments, _ := c.GetComments(courseId)
	pastComments := comments.Data
	comments.Data = append(pastComments, newComment)
	if _, err := c.db.Get(courseId); err != nil {
		c.db.Insert(comments)
	} else {
		c.db.Update(comments)
	}
	return comments, nil
}

func (c *CommentsService) GetComments(courseId string) (Comments, error) {
	var comments Comments
	var err error
	if comments, err = c.db.Get(courseId); err != nil {
		pastComments := make([]Comment, 0)
		comments = Comments{
			CourseId: courseId,
			Data:     pastComments,
		}
	}
	return comments, nil
}

func CreateCommentsService(commentsDb db.DB[Comments]) *CommentsService {
	return &CommentsService{
		commentsDb,
	}
}
