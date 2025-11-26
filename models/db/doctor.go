package models

import (
	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	Name           string   `gorm:"not null"`
	Specialization string   `gorm:"not null"`
	Phone          string   `gorm:"not null"`
	Experience     string   `gorm:"type:text"`
	Location       string   `gorm:"not null"`
	License        string   `gorm:"uniqueIndex;not null"`
	WorkHours      string   `gorm:"type:varchar(100)"`
	Availability   bool     `gorm:"default:true"`
	WorkDays       []string `gorm:"type:json"`
	WorkTime       []string `gorm:"type:json"`
	WorkTimeEnd    []string `gorm:"type:json"`
}

func NewDoctor(name string, specialization string, phone string, experience string, location string, license string, workHours string, availability bool, workDays []string, workTime []string, workTimeEnd []string) *Doctor {
	return &Doctor{
		Name:           name,
		Specialization: specialization,
		Phone:          phone,
		Experience:     experience,
		Location:       location,
		License:        license,
		WorkHours:      workHours,
		Availability:   availability,
		WorkDays:       workDays,
		WorkTime:       workTime,
		WorkTimeEnd:    workTimeEnd,
	}
}
