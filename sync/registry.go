package sync

import (
	"fmt"
	"sync"
)

// Registry maintains a collection of available sync providers
type Registry struct {
	factories map[string]SyncProviderFactory
	mu        sync.RWMutex
}

// NewRegistry creates a new provider registry
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]SyncProviderFactory),
	}
}

// Register adds a new provider factory to the registry
func (r *Registry) Register(providerType string, factory SyncProviderFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[providerType] = factory
}

// CreateProvider creates a new sync provider based on its type and configuration
func (r *Registry) CreateProvider(config ProviderConfig) (SyncProvider, error) {
	r.mu.RLock()
	factory, exists := r.factories[config.Type]
	r.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("unknown sync provider type: %s", config.Type)
	}
	
	return factory(config)
}

// ListProviderTypes returns a list of all registered provider types
func (r *Registry) ListProviderTypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	types := make([]string, 0, len(r.factories))
	for t := range r.factories {
		types = append(types, t)
	}
	
	return types
}

// DefaultRegistry is the global sync provider registry
var DefaultRegistry = NewRegistry()

// Register adds a new provider factory to the default registry
func Register(providerType string, factory SyncProviderFactory) {
	DefaultRegistry.Register(providerType, factory)
}

// CreateProvider creates a new sync provider from the default registry
func CreateProvider(config ProviderConfig) (SyncProvider, error) {
	return DefaultRegistry.CreateProvider(config)
}

// ListProviderTypes returns a list of all registered provider types in the default registry
func ListProviderTypes() []string {
	return DefaultRegistry.ListProviderTypes()
}