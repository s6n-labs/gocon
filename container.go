package gocon

import (
	"sync"
)

type Container interface {
	Get(key string) (*Definition, error)
	GetTagged(tag string) ([]*Definition, error)
	Set(def *Definition) error
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
		return nil, ErrServiceNotFound
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

func (c *unsafeContainer) DisposeAll() {
	for _, def := range c.services {
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
