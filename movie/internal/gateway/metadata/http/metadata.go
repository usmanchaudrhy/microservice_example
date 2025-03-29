package http

import (
	model "awesomeProject/metadata/pkg"
	"awesomeProject/movie/internal/gateway"
	"awesomeProject/pkg/discovery"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

// gateway defines a movie metadata HTTP gateway
type Gateway struct {
	registry discovery.Registry
}

// New creates a new HTTP gateway for a mvovie metadata service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// get gets the movie metadata by the movie ID
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {

	// changing the get method to get the service address
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/metadata"

	log.Printf("calling metadata service. Request GET: %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()

	values.Add("id", id)
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}

	var v *model.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}
