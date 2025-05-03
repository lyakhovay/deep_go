package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node[T cmp.Ordered] struct {
	left, right *Node[T]
	key         T
	value       any
}

type OrderedMap[T cmp.Ordered] struct {
	root *Node[T]
	size int
}

func NewOrderedMap[T cmp.Ordered]() OrderedMap[T] {
	return OrderedMap[T]{}
}

func (m *OrderedMap[T]) Insert(key T, value any) {
	m.root = m.insert(key, value, m.root)
}

func (m *OrderedMap[T]) insert(key T, value any, node *Node[T]) *Node[T] {
	if node == nil {
		m.size++
		return &Node[T]{key: key, value: value}
	}
	if key < node.key {
		node.left = m.insert(key, value, node.left)
	}
	if key > node.key {
		node.right = m.insert(key, value, node.right)
	}
	node.value = value // key == node.key
	return node
}

func (m *OrderedMap[T]) Erase(key T) {
	m.root = m.erase(key, m.root)
}

func (m *OrderedMap[T]) erase(key T, node *Node[T]) *Node[T] {
	if node == nil {
		return nil
	}
	if key < node.key {
		node.left = m.erase(key, node.left)
		return node
	}
	if key > node.key {
		node.right = m.erase(key, node.right)
		return node
	}
	m.size--
	if node.left == nil {
		return node.right
	}
	if node.right == nil {
		return node.left
	}
	lastChild := node.right
	for lastChild.left != nil {
		lastChild = lastChild.left
	}
	lastChild.left = node.left
	lastChild.right = node.right
	node = lastChild
	return node
}

func (m *OrderedMap[T]) Contains(key T) bool {
	return contains(key, m.root)
}

func contains[T cmp.Ordered](key T, node *Node[T]) bool {
	if node == nil {
		return false
	}
	if key < node.key {
		return contains(key, node.left)
	}
	if key > node.key {
		return contains(key, node.right)
	}
	return true // key == node.key
}

func (m *OrderedMap[T]) Size() int {
	return m.size
}

func (m *OrderedMap[T]) ForEach(action func(T, any)) {
	if m.root != nil {
		forEach(action, m.root)
	}
}

func forEach[T cmp.Ordered](action func(T, any), node *Node[T]) {
	if node.left != nil {
		forEach(action, node.left)
	}
	action(node.key, node.value)
	if node.right != nil {
		forEach(action, node.right)
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key int, _ any) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key int, _ any) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
