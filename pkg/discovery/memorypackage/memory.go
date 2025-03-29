package memory

import (
	"awesomeProject/pkg/discovery"
	"context"
	"errors"
	"sync"
	"time"
)

type serviceName string
type instanceID string

type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistery() *Registry {
	return &Registry{
		serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{},
	}
}

// Implementing the Register and Deregister functions
// Register creates a new service record in the registry
func (r *Registry) Register(ctx context.Context, iID string, sName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName(sName)]; !ok {
		r.serviceAddrs[serviceName(sName)] = map[instanceID]*serviceInstance{}
	}

	r.serviceAddrs[serviceName(sName)][instanceID(iID)] = &serviceInstance{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}

	return nil
}

// Deregister removes a service record form the registry
func (r *Registry) Deregister(ctx context.Context, instanceId string, servicename string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName(servicename)]; !ok {
		return nil
	}

	delete(r.serviceAddrs[serviceName(servicename)], instanceID(instanceId))
	return nil
}

// Report health is a push mechanism to pushg healthy state to the registry
func (r *Registry) ReportHealthyState(instance_id string, service_name string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName(service_name)]; !ok {
		return errors.New("service is not registered yet")
	}

	if _, ok := r.serviceAddrs[serviceName(service_name)][instanceID(instance_id)]; !ok {
		return errors.New("service instance is not registered yet")
	}

	r.serviceAddrs[serviceName(service_name)][instanceID(instance_id)].lastActive = time.Now()
	return nil
}

// returns the list of addresses of active instances of the given  service
func (r *Registry) ServiceAddresses(ctx context.Context, service_name string) ([]string, error) {
	r.Lock()
	defer r.Unlock()

	if len(r.serviceAddrs[serviceName(service_name)]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string
	for _, i := range r.serviceAddrs[serviceName(service_name)] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}

		res = append(res, i.hostPort)
	}

	return res, nil
}
