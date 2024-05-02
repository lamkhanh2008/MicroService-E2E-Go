package memory

import (
	"context"
	"errors"
	"microservice/pkg/discovery"
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

func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{}}
}

func (r *Registry) Register(ctx context.Context, instanceid string, servicename string, hostPort string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName(servicename)]; !ok {
		r.serviceAddrs[serviceName(servicename)] = map[instanceID]*serviceInstance{}
	}
	r.serviceAddrs[serviceName(servicename)][instanceID(instanceid)] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

func (r *Registry) Deregister(ctx context.Context, insinstanceid string, servicename string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName(servicename)]; !ok {
		return nil
	}
	delete(r.serviceAddrs[serviceName(servicename)], instanceID(insinstanceid))
	return nil
}

func (r *Registry) ReportHealthyState(instanceid string, servicename string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName(servicename)]; !ok {
		return errors.New("Service is not registered")
	}
	if _, ok := r.serviceAddrs[serviceName(servicename)][instanceID(instanceid)]; !ok {
		return errors.New("service instance is not registered yet")
	}
	r.serviceAddrs[serviceName(servicename)][instanceID(instanceid)].lastActive = time.Now()
	return nil

}

func (r *Registry) ServiceAddresses(ctx context.Context, servicename string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.serviceAddrs[serviceName(servicename)]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, i := range r.serviceAddrs[serviceName(servicename)] {
		if i.lastActive.Before(time.Now().Add(-5 * time.
			Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}
