package movie

import (
	metadatamodel "awesomeProject/metadata/pkg"
	"awesomeProject/movie/internal/gateway"
	"awesomeProject/movie/pkg/model"
	ratingmodel "awesomeProject/rating/pkg"
	"context"
	"errors"
)

// ErrNotFound is returned when the movie metadata is not found
var ErrNotFound = errors.New("movie metadata not found")

type ratingGateway interface {
	GetAggregatedRating(
		ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType,
	) (float64, error)
	PutRating(
		ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType,
		rating *ratingmodel.Rating,
	) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadatamodel.Metadata, error)
}

// Now we can define our service controller
// Controller defines a movie service controller
type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

// New creates a new movie service controller
func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{ratingGateway, metadataGateway}
}

// Finally implement the function for getting the movie details, including both rating and metadata
func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	details := &model.MovieDetails{Metadata: *metadata}
	rating, err := c.ratingGateway.GetAggregatedRating(
		ctx, ratingmodel.RecordID(id), ratingmodel.RecordTypeMovie,
	)

	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		// we can procced in this case, its okay to not have a rating
	} else if err != nil {
		return nil, err
	} else {
		details.Rating = &rating
	}

	return details, nil
}
