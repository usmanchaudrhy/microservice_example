package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Register defines a service registery
type Registry interface {
	// Register creates a service entry record in the registery
	Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error
	// De-register will remove the service instance record from the registery
	Deregister(ctx context.Context, instanceID string, serviceName string) error
	// Service addresses returns the list of the addresses of active instances of the given service
	ServiceAddresses(ctx context.Context, serviceID string) ([]string, error)
	// Push mechanism for reporting healthy state to the registery
	ReportHealthyState(instanceID string, serviceName string) error
}

// ErrNotFound is returned when no service addresses are found
var ErrNotFound = errors.New("no service address found")

// Generate instanceID generates a pseudo-random service identifier, using a servicve name
// sufficed by a dash and a random number
func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
