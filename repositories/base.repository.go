package repositories

import (
	"context"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type newBaseRepository[T Model] struct {
	collection *mgm.Collection
}

func BaseRepository[T Model](model T) GenericRepository[T] {
	return &newBaseRepository[T]{
		collection: mgm.Coll(model),
	}
}

func (r *newBaseRepository[T]) FindAll() ([]T, error) {
	var results []T
	err := r.collection.SimpleFind(&results, bson.M{})
	return results, err
}

func (r *newBaseRepository[T]) FindAllPaginated(page int, limit int) ([]T, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit

	var results []T
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	err := r.collection.SimpleFind(&results, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(context.TODO(), bson.M{})
	return results, total, err
}

func (r *newBaseRepository[T]) FindByID(id string) (T, error) {
	var result T
	err := r.collection.FindByID(id, result)
	return result, err
}

func (r *newBaseRepository[T]) Create(entity T) error {
	return r.collection.Create(entity)
}

func (r *newBaseRepository[T]) Update(entity T) error {
	return r.collection.Update(entity)
}

func (r *newBaseRepository[T]) Delete(id string) error {
	var result T
	err := r.collection.FindByID(id, result)
	if err != nil {
		return err
	}
	return r.collection.Delete(result)
}