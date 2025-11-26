package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// StringArray is a custom type for handling JSON arrays in MySQL
type StringArray []string

// Value implements the driver.Valuer interface
func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "[]", nil
	}
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("cannot scan non-string value into StringArray")
	}

	return json.Unmarshal(bytes, a)
}

type Doctor struct {
	gorm.Model
	Name           string      `gorm:"not null"`
	Specialization string      `gorm:"not null"`
	Phone          string      `gorm:"not null"`
	Experience     string      `gorm:"type:text"`
	Location       string      `gorm:"not null"`
	License        string      `gorm:"uniqueIndex;not null"`
	WorkHours      string      `gorm:"type:varchar(100)"`
	Availability   bool        `gorm:"default:true"`
	WorkDays       StringArray `gorm:"type:json"`
	WorkTime       StringArray `gorm:"type:json"`
	WorkTimeEnd    StringArray `gorm:"type:json"`
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
		WorkDays:       StringArray(workDays),
		WorkTime:       StringArray(workTime),
		WorkTimeEnd:    StringArray(workTimeEnd),
	}
}
