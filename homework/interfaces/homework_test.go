package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}
type Singleton struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	dependencies map[string]interface{}
}

func NewContainer() *Container {
	return &Container{dependencies: make(map[string]interface{})}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	c.dependencies[name] = constructor
}

func (c *Container) RegisterSingletonType(name string, constructor interface{}) {
	if _, ok := c.dependencies[name]; ok {
		return
	}
	if fn, ok := constructor.(func() interface{}); ok {
		c.dependencies[name] = fn()
	}
}

func (c *Container) Resolve(name string) (interface{}, error) {
	constructor, ok := c.dependencies[name]
	if !ok {
		return nil, errors.New("no constructor registered")
	}
	switch fn := constructor.(type) {
	case func() interface{}:
		return fn(), nil
	case interface{}:
		return fn, nil
	}
	return nil, errors.New("unknown type of constructor")
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})
	container.RegisterSingletonType("Singleton", func() interface{} {
		return &Singleton{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)

	service3, err := container.Resolve("Singleton")
	assert.NoError(t, err)
	service4, err := container.Resolve("Singleton")
	assert.NoError(t, err)
	s3 := service3.(*Singleton)
	s4 := service4.(*Singleton)
	assert.True(t, s3 == s4)
}
