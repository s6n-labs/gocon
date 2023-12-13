package gocon

import (
	"fmt"
	"sync"
)

type ServiceNotFoundError struct {
	key string
}

func (e ServiceNotFoundError) Error() string {
	return fmt.Sprintf("service '%s' does not exist, or cannot be resolved", e.key)
}

type Container interface {
	Get(key string) (*Definition, error)
	GetTagged(tag string) ([]*Definition, error)
	Set(def *Definition) error
	All() []*Definition
	DisposeAll()
}

type unsafeContainer struct {
	services map[string]*Definition
	tags     map[string]map[string]struct{}
	inherit  Container
}

func newUnsafeContainer(inherit Container) *unsafeContainer {
	return &unsafeContainer{
		services: make(map[string]*Definition),
		tags:     make(map[string]map[string]struct{}),
		inherit:  inherit,
	}
}

func (c *unsafeContainer) Get(key string) (*Definition, error) {
	service, ok := c.services[key]
	if !ok {
		return nil, ServiceNotFoundError{
			key: key,
		}
	}

	return service, nil
}

func (c *unsafeContainer) GetTagged(tag string) ([]*Definition, error) {
	tagged, ok := c.tags[tag]
	if !ok {
		return make([]*Definition, 0), nil
	}

	defs := make([]*Definition, 0, len(tagged))
	for t := range tagged {
		def, ok := c.services[t]
		if !ok {
			continue
		}

		defs = append(defs, def)
	}

	return defs, nil
}

func (c *unsafeContainer) Set(def *Definition) error {
	c.services[def.Key] = def

	if def.configureFunc != nil {
		if err := def.configureFunc(c); err != nil {
			return err
		}
	}

	for _, t := range def.Tags {
		tagged, ok := c.tags[t]
		if !ok {
			tagged = make(map[string]struct{})
			c.tags[t] = tagged
		}

		tagged[def.Key] = struct{}{}
	}

	return nil
}

func (c *unsafeContainer) All() []*Definition {
	defs := make([]*Definition, 0, len(c.services))
	for _, def := range c.services {
		defs = append(defs, def)
	}

	return defs
}

func (c *unsafeContainer) DisposeAll() {
	for _, def := range c.All() {
		if def.Value == nil {
			continue
		}

		dispose(*def.Value)
	}
}

type container struct {
	unsafeContainer
	mutex sync.RWMutex
}

func NewContainer(inherit Container) Container {
	return &container{
		unsafeContainer: *newUnsafeContainer(inherit),
		mutex:           sync.RWMutex{},
	}
}

func (c *container) Get(key string) (*Definition, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.unsafeContainer.Get(key)
}

func (c *container) GetTagged(tag string) ([]*Definition, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.unsafeContainer.GetTagged(tag)
}

func (c *container) Set(def *Definition) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.unsafeContainer.Set(def)
}

func (c *container) All() []*Definition {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.unsafeContainer.All()
}

func (c *container) DisposeAll() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.unsafeContainer.DisposeAll()
}
