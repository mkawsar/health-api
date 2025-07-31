package models

import (
	"github.com/kamva/mgm/v3"
)

type Doctor struct {
	mgm.DefaultModel `bson:",inline"`
	Name string `bson:"name"`
	Specialization string `bson:"specialization"`
	Phone string `bson:"phone"`
	Experience string `bson:"experience"`
	Location string `bson:"location"`
	License string `bson:"license"`
	WorkHours string `bson:"work_hours"`
	Availability bool `bson:"availability"`
	WorkDays []string `bson:"work_days"`
	WorkTime []string `bson:"work_time"`
	WorkTimeEnd []string `bson:"work_time_end"`
}

func NewDoctor(name string, specialization string, phone string, experience string, location string, license string, workHours string, availability bool, workDays []string, workTime []string, workTimeEnd []string) *Doctor {
	return &Doctor{
		Name: name,
		Specialization: specialization,
		Phone: phone,
		Experience: experience,
		Location: location,
		License: license,
		WorkHours: workHours,
		Availability: availability,
		WorkDays: workDays,
		WorkTime: workTime,
		WorkTimeEnd: workTimeEnd,
	}
}

func (model *Doctor) CollectionName() string {
	return "doctors"
}
