package provider

import (
	"fmt"
	"sync"
)

// Registry holds all initialized AI providers and resolves them by provider name.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

// NewRegistry creates a new empty provider registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry under the given name.
func (r *Registry) Register(name string, p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = p
}

// Get returns the provider registered under the given name.
func (r *Registry) Get(name string) (Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("no AI provider registered for: %q", name)
	}
	return p, nil
}

// MustGet returns the provider or panics if not found.
func (r *Registry) MustGet(name string) Provider {
	p, err := r.Get(name)
	if err != nil {
		panic(err)
	}
	return p
}
