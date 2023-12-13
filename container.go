package gocon

import (
	"sync"
)

type Container interface {
	get(key string) (*Definition, error)
	getTagged(tag string) ([]*Definition, error)
	set(def *Definition) error
	lock() error
	unlock() error
	parent() Container
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

func (c *unsafeContainer) get(key string) (*Definition, error) {
	service, ok := c.services[key]
	if !ok {
		return nil, ErrServiceNotFound
	}

	return service, nil
}

func (c *unsafeContainer) getTagged(tag string) ([]*Definition, error) {
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

func (c *unsafeContainer) set(def *Definition) error {
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

func (c *unsafeContainer) lock() error {
	return nil
}

func (c *unsafeContainer) unlock() error {
	return nil
}

func (c *unsafeContainer) parent() Container {
	return c.inherit
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

func (c *container) get(key string) (*Definition, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.unsafeContainer.get(key)
}

func (c *container) getTagged(tag string) ([]*Definition, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.unsafeContainer.getTagged(tag)
}

func (c *container) set(def *Definition) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.unsafeContainer.set(def)
}

func (c *container) lock() error {
	c.mutex.Lock()

	return nil
}

func (c *container) unlock() error {
	c.mutex.Unlock()

	return nil
}
