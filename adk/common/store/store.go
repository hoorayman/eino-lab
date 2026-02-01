package store

import (
	"context"
	"sync"

	"github.com/cloudwego/eino/compose"
)

// NewInMemoryStore 创建一个线程安全的内存 CheckPointStore
// 用于存储 Agent 执行中断时的检查点状态
func NewInMemoryStore() compose.CheckPointStore {
	return &inMemoryStore{
		mem: make(map[string][]byte),
	}
}

// inMemoryStore 内存实现的 CheckPointStore
type inMemoryStore struct {
	mu  sync.RWMutex        // 读写锁，保护 map 的并发访问
	mem map[string][]byte    // 存储 key-value 数据
}

// Set 设置 key-value 到存储中
// 这是线程安全的操作，可以并发调用
func (i *inMemoryStore) Set(ctx context.Context, key string, value []byte) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.mem[key] = value
	return nil
}

// Get 从存储中获取 key 对应的值
// 这是线程安全的操作，支持并发读
func (i *inMemoryStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	v, ok := i.mem[key]
	return v, ok, nil
}
