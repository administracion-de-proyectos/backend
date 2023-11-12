package services

import "backend-admin-proyect/src/db"

type Groups struct {
	db db.DB[GroupDto]
}

func (g *Groups) AddToGroup(teacherEmail, studentEmail string) error {
	data := g.GetGroup(teacherEmail)
	data = append(data, studentEmail)
	if _, err := g.db.Get(teacherEmail); err != nil {
		g.db.Insert(GroupDto{teacherEmail, data})
	} else {
		g.db.Update(GroupDto{teacherEmail, data})
	}
	return nil
}

func (g *Groups) GetGroup(teacherEmail string) []string {
	if d, err := g.db.Get(teacherEmail); err != nil {
		return make([]string, 0)
	} else {
		return d.StudentsGroup
	}
}

func CreateGroupsService(db db.DB[GroupDto]) *Groups {
	return &Groups{
		db,
	}
}
