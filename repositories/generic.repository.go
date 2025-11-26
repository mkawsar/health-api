package repositories

type GenericRepository[T Model] interface {
	FindAll() ([]T, error)
	FindAllPaginated(page int, limit int) ([]T, int64, error)
	FindByID(id uint) (T, error)
	Create(entity T) error
	Update(entity T) error
	Delete(id uint) error
}
