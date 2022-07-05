package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func NewKVStoreService() *KVStoreService {
	return &KVStoreService{
		m:      make(map[string]string),
		filter: make(map[string]func(key string)),
	}
}

type KVStoreService struct {
	// map lưu trữ dữ liệu key value
	m map[string]string

	// map chứa danh sách các hàm filter
	// được xác định trong mỗi lời gọi
	filter map[string]func(key string)

	// bảo vệ các thành phần khác khi được truy cập
	// và sửa đổi từ nhiều Goroutine cùng lúc
	mu sync.Mutex
}

func (p *KVStoreService) Get(key string, value *string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if v, ok := p.m[key]; ok {
		*value = v
		return nil
	}

	return fmt.Errorf("not found")
}

func (p *KVStoreService) Set(kv [2]string, reply *struct{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	key, value := kv[0], kv[1]

	if oldValue := p.m[key]; oldValue != value {
		// hàm filter được gọi khi value tương ứng
		// với key bị sửa đổi
		for _, fn := range p.filter {
			fn(key)
		}
	}

	p.m[key] = value

	fmt.Println("Setting")
	return nil
}

// Watch trả về key mỗi khi nhận thấy có thay đổi
func (p *KVStoreService) Watch(timeoutSecond int, keyChanged *string) error {
	// id là một string ghi lại thời gian hiện tại
	id := fmt.Sprintf("watch-%s-%03v", time.Now(), rand.Int())

	// buffered channel chứa key
	ch := make(chan string, 10)

	// filter để theo dõi key thay đổi
	p.mu.Lock()
	p.filter[id] = func(key string) { ch <- key }
	p.mu.Unlock()

	select {
	// trả về timeout sau một khoảng thời gian
	case <-time.After(time.Duration(timeoutSecond) * time.Second):
		fmt.Println("timeout")
		return fmt.Errorf("timeout")
	case key := <-ch:
		*keyChanged = key
		return nil
	}
}
