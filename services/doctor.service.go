package services

import (
	"context"
	"errors"

	db "health/models/db"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetDoctors(ctx context.Context, page int, limit int, nameFilter string) ([]db.Doctor, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit
	filter := bson.M{}
	if nameFilter != "" {
		filter["name"] = bson.M{"$regex": nameFilter, "$options": "i"} // Case-insensitive search
	}
	var doctors []db.Doctor
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(skip))
	err := mgm.Coll(&db.Doctor{}).SimpleFind(&doctors, filter, opts)

	if err != nil {
		return nil, 0, errors.New("cannot find doctors")
	}
	total, _ := mgm.Coll(&db.Doctor{}).CountDocuments(ctx, filter)
	return doctors, total, nil
}
