package repositories

import (
	"health/services"

	"gorm.io/gorm"
)

type newBaseRepository[T Model] struct {
	db *gorm.DB
}

func BaseRepository[T Model]() GenericRepository[T] {
	return &newBaseRepository[T]{
		db: services.DB,
	}
}

func (r *newBaseRepository[T]) FindAll() ([]T, error) {
	var results []T
	err := r.db.Find(&results).Error
	return results, err
}

func (r *newBaseRepository[T]) FindAllPaginated(page int, limit int) ([]T, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var results []T
	var total int64

	err := r.db.Model(&results).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *newBaseRepository[T]) FindByID(id uint) (T, error) {
	var result T
	err := r.db.First(&result, id).Error
	return result, err
}

func (r *newBaseRepository[T]) Create(entity T) error {
	return r.db.Create(&entity).Error
}

func (r *newBaseRepository[T]) Update(entity T) error {
	return r.db.Save(&entity).Error
}

func (r *newBaseRepository[T]) Delete(id uint) error {
	var result T
	err := r.db.First(&result, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(&result).Error
}
