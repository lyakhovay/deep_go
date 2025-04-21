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
	node := m.root
	var lastNode *Node[T]
	for node != nil {
		if node.key == key {
			node.value = value
		} else {
			lastNode = node
			if key < node.key {
				node = node.left
			} else {
				node = node.right
			}
		}
	}
	newNode := Node[T]{key: key, value: value}
	if lastNode == nil {
		m.root = &newNode
	} else {
		if key < lastNode.key {
			lastNode.left = &newNode
		} else {
			lastNode.right = &newNode
		}
	}
	m.size++
}

func (m *OrderedMap[T]) Erase(key T) {
	node := m.root
	var lastNode *Node[T]
	for node != nil {
		if node.key == key {
			break
		}
		lastNode = node
		if key < node.key {
			node = node.left
		} else {
			node = node.right
		}
	}
	if node == nil || lastNode == nil {
		return
	}
	m.size--
	if node.right == nil {
		if node == lastNode.left {
			lastNode.left = node.left
			return
		}
		lastNode.right = node.left
		return
	}
	leftChild := node.right
	var lastChild *Node[T]
	for leftChild.left != nil {
		lastChild = leftChild
		leftChild = leftChild.left
	}
	if lastChild != nil {
		lastChild.left = leftChild.right
	} else {
		node.right = leftChild.right
	}
	node.key = leftChild.key
	node.value = leftChild.value
}

func (m *OrderedMap[T]) Contains(key T) bool {
	node := m.root
	for node != nil {
		if node.key == key {
			return true
		}
		if key < node.key {
			node = node.left
		} else {
			node = node.right
		}
	}
	return false
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
