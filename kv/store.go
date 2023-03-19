package kv

import (
	"time"
)

// KVStore是一个接口类型，使用泛型来定义了KV存储的基本操作
// [golang] 接口定义，一些函数的集合，与java的interface类似
// `[K comparable, V any]`定义了这个接口使用哪些泛型类型。
// `K`表示key的类型，要求K必须是`comparable`(这里不具体展开了)
// `V`表示value的类型，V可以是任何类型
type KVStore[K comparable, V any] interface {
	// Put函数将一个键值对存储到KV存储中
	// Put方法接收两个参数，一个是键key，一个是值val
	// [golang] 参数变量名在前，变量类型在后。`key K`中key是变量名，K是key的类型
	Put(key K, val V)
	// Get函数返回与指定键相关联的值
	// Get方法接收两个参数，一个是键key，另一个是超时时间timeout。
	// golang中函数可以返回多个值，在这里返回了两个值，第一个值的类型是V，第二个值的类型是error
	// Get函数的逻辑如下：
	// 如果能够找到key对应的值val，则立即返回val。
	// 如果不能找到key对应的值val，则会阻塞等待，
	// 直到key对应的值被其他协程(goroutine)设置，或者超时(timeout时间)。
	Get(key K, timeout time.Duration) (V, error)
}

// 实现这个KVStore的`类`在kv/concurrent.go中，见`ConcurrentMap`
