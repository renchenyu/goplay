// 这个程序主要用来尽可能地展示golang语言的语法，让初学者有个直观的了解

// 一个可运行的go程序需要有一个main package，在这个package中需要有一个main函数作为程序的入口
package main

// 引入依赖的package
import (
	// 这部分都是golang标准库中的package
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	// 这个是项目中的package。在go.mod中定义了这个module的名字(goplay)，kv对应本项目kv/目录
	// 代码使用`kv.`来使用这个package中定义的结构体、函数等
	"goplay/kv"
)

// 主函数，名称必须是main
// 这个程序主要是通过并发地对一个Map进行读和写来演示golang的相关特性
func main() {
	// 这里是一条变量定义并初始化的语句
	// 同大部分语言一样，局部变量的生命周期是一个{}来表示的block
	// 使用`:=`可以省去变量的类型，直接根据`:=`右边的类型来推断出左边的类型。在这里右边的范围值是*ConcurrentMap[int, [2]int]，所以myMap的类型就是*ConcurrentMap[int, [2]int]
	// PS：这个类型有点复杂，先不用纠结细节
	// 接着可以打开kv/store.go继续浏览
	myMap := kv.NewConcurrentMap[int, [2]int]()

	// 这里调用函数playKv
	playKv(context.Background(), myMap)
}

func playKv(ctx context.Context, kv kv.KVStore[int, [2]int]) {
	const size = 10

	keys := make([]int, 0, size)
	for i := 0; i < size; i++ {
		keys = append(keys, rand.Intn(10))
	}

	var wg sync.WaitGroup

	for _, key := range keys {
		wg.Add(1)
		go func(key int, timeout time.Duration) {
			defer wg.Done()

			fmt.Printf("key: %d, timeout: %s\n", key, timeout)
			if v, err := kv.Get(key, timeout); err != nil {
				fmt.Printf("key: %d, err: %s\n", key, err)
			} else {
				fmt.Printf("key: %d, value: %d\n", key, v)
			}
		}(key, time.Duration(rand.Intn(10))*time.Second)
	}

	time.Sleep(time.Second)

	for len(keys) > 0 {
		first := keys[0]
		keys = keys[1:]

		go func(key int, wait time.Duration) {
			time.Sleep(wait)
			kv.Put(key, [2]int{key, key + 1})
		}(first, time.Duration(rand.Intn(5))*time.Second)
	}

	wg.Wait()
}
