package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	"k8s.io/client-go/tools/cache"
)

var (
	// Metrics
	attestationCacheMiss = promauto.NewCounter(prometheus.CounterOpts{
		Name: "attestation_cache_miss",
		Help: "The number of attestation data requests that aren't present in the cache.",
	})
	attestationCacheHit = promauto.NewCounter(prometheus.CounterOpts{
		Name: "attestation_cache_hit",
		Help: "The number of attestation data requests that are present in the cache.",
	})
	attestationCacheSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "attestation_cache_size",
		Help: "The number of attestation data in the attestations cache",
	})
)

type AttestationCache struct {
	cache      *cache.FIFO
	lock       sync.RWMutex
	inProgress map[string]bool
}

func NewAttestationCache() *AttestationCache {
	return &AttestationCache{
		cache:      cache.NewFIFO(wrapperToKey),
		inProgress: make(map[string]bool),
	}
}

func (c *AttestationCache) Get(ctx context.Context, req *pb.AttestationDataRequest) (*pb.AttestationDataResponse, error) {
	if req == nil {
		return nil, errors.New("nil attestation data request")
	}

	s, e := reqToKey(req)
	if e != nil {
		return nil, e
	}

	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		c.lock.RLock()
		if !c.inProgress[s] {
			c.lock.RUnlock()
			break
		}
		c.lock.RUnlock()
	}

	item, exists, err := c.cache.GetByKey(s)
	if err != nil {
		return nil, err
	}

	if exists && item != nil && item.(*attestationReqResWrapper).res != nil {
		attestationCacheHit.Inc()
		return item.(*attestationReqResWrapper).res, nil
	} else {
		attestationCacheMiss.Inc()
		return nil, nil
	}
}

func (c *AttestationCache) MarkInProgress(req *pb.AttestationDataRequest) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	s, e := reqToKey(req)
	if e != nil {
		return e
	}
	if c.inProgress[s] {
		return errors.New("already in progress")
	}
	c.inProgress[s] = true
	return nil
}

func (c *AttestationCache) MarkNotInProgress(req *pb.AttestationDataRequest) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	s, e := reqToKey(req)
	if e != nil {
		return e
	}
	delete(c.inProgress, s)
	return nil
}

func (c *AttestationCache) Put(ctx context.Context, req *pb.AttestationDataRequest, res *pb.AttestationDataResponse) error {
	attestationCacheSize.Inc()
	data := &attestationReqResWrapper{
		req,
		res,
	}
	if err := c.cache.AddIfNotPresent(data); err != nil {
		return err
	}
	trim(c.cache, maxCacheSize)
	return nil
}

func wrapperToKey(i interface{}) (string, error) {
	w := i.(*attestationReqResWrapper)
	if w == nil {
		return "", errors.New("nil wrapper")
	}
	if w.req == nil {
		return "", errors.New("nil wrapper.request")
	}
	return reqToKey(w.req)
}

func reqToKey(req *pb.AttestationDataRequest) (string, error) {
	return fmt.Sprintf("%d-%d", req.Shard, req.Slot), nil
}

type attestationReqResWrapper struct {
	req *pb.AttestationDataRequest
	res *pb.AttestationDataResponse
}
