package memory

import (
	"awesomeProject/rating/internal/repository"
	model "awesomeProject/rating/pkg"
	"context"
)

type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

func New() *Repository {
	return &Repository{
		data: map[model.RecordType]map[model.RecordID][]model.Rating{},
	}
}

// Get retrieves all ratings for a given record
func (r *Repository) Get(ctx context.Context,
	recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {

	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}

	if ratings, ok := r.data[recordType][recordID]; !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}

	return r.data[recordType][recordID], nil
}

func (r *Repository) Put(
	ctx context.Context,
	recordID model.RecordID, recordType model.RecordType,
	rating *model.Rating,
) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]model.Rating{}
	}

	r.data[recordType][recordID] = append(
		r.data[recordType][recordID], *rating)

	return nil
}
