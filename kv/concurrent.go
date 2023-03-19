package kv

import (
	"context"
	"sync"
	"time"
)

// 为了实现Get函数的功能，我们先定义一个future的概念。
// 当调用Get(key)时，无论key对应的val是否存在，总是能获得一个future。
// 对于这个future有两种状态
// 1. future已经完成了，那么可以立即获得val
// 2. future还没有完成，需要等待
// 这里定义了3个成员变量
// 1. val是future的结果
// 2. doneCh用于在future完成时，恢复等待中的get
// 3. closed用来标识doneCh是否已经被关闭
// [golang] struct类似java中的class。golang中参数名和参数类型一起出现时，总是名称在前，类型在后。
type future[V any] struct {
	val    V
	doneCh chan struct{}
	closed bool
}

// [golang] struct的method定义与实现
// 这个方法等价于 func done(f *future[V]) <-chan struct{}
func (f *future[V]) done() <-chan struct{} {
	return f.doneCh
}

func (f *future[V]) set(val V) {
	f.val = val
	if !f.closed {
		f.closed = true
		close(f.doneCh)
	}
}

func (f *future[V]) get() V {
	<-f.doneCh
	return f.val
}

func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		store: make(map[K]*future[V]),
	}
}

type ConcurrentMap[K comparable, V any] struct {
	sync.Mutex
	store map[K]*future[V]
}

func (c *ConcurrentMap[K, V]) getOrNewFuture(key K) *future[V] {
	c.Lock()
	defer c.Unlock()

	if f, ok := c.store[key]; ok {
		return f
	} else {
		f := future[V]{
			doneCh: make(chan struct{}),
		}
		c.store[key] = &f
		return &f
	}
}

func (c *ConcurrentMap[K, V]) Put(key K, val V) {
	f := c.getOrNewFuture(key)
	f.set(val)
}

func (c *ConcurrentMap[K, V]) Get(key K, timeout time.Duration) (V, error) {
	f := c.getOrNewFuture(key)
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	var v V
	select {
	case <-ctx.Done():
		return v, ctx.Err()
	case <-f.done():
		return f.get(), nil
	}
}
