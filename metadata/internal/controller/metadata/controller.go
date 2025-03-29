package metadata

import (
	model "awesomeProject/metadata/pkg"
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
}

// Controller defines a metadate service controller
type Controller struct {
	repo metadataRepository
}

// new function creates metadata service controller
func New(repo metadataRepository) *Controller {
	return &Controller{repo: repo}
}

// Get returns movie metadata by id
func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}
