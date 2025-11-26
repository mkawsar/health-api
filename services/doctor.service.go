package services

import (
	"context"
	"errors"

	db "health/models/db"
)

func GetDoctors(ctx context.Context, page int, limit int, nameFilter string) ([]db.Doctor, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := DB.Model(&db.Doctor{})
	if nameFilter != "" {
		query = query.Where("name ILIKE ?", "%"+nameFilter+"%") // Case-insensitive search
	}

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, errors.New("cannot count doctors")
	}

	var doctors []db.Doctor
	err = query.Offset(offset).Limit(limit).Find(&doctors).Error
	if err != nil {
		return nil, 0, errors.New("cannot find doctors")
	}

	return doctors, total, nil
}
