package repositories

import "github.com/kamva/mgm/v3"

type doctorRepository[T Model] struct {
	collection *mgm.Collection
}
