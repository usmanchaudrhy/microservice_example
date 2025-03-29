package model

import model "awesomeProject/metadata/pkg"

// MovieDetails includes the movie metadata and its agregated ratings
type MovieDetails struct {
	Rating   *float64
	Metadata model.Metadata
}
