package http

import (
	"awesomeProject/movie/internal/gateway"
	model "awesomeProject/rating/pkg"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// defines an http gateway for the rating service
type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr}
}

// Returns the aggregated rating for a record or ErrNotFound if there are no ratings for it
func (g *Gateway) GetAggregatedRating(
	ctx context.Context, recordID model.RecordID, recordType model.RecordType,
) (float64, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/rating", nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)

	// setting the URL params for the request
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", string(recordType))

	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		// the response is not in the 2xx range
		return 0, fmt.Errorf("non-2xx response: %v", resp)
	}

	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}

	return v, nil
}

// PutRating handles the rating creation request
func (g *Gateway) PutRating(
	ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating,
) error {
	req, err := http.NewRequest(http.MethodPut, g.addr+"/rating", nil)
	if err != nil {
		return err
	}

	// add the URL params to the request
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}

	return nil
}
