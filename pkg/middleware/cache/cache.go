package cache

import (
	"sync"
	"time"

	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
)

const expireDuration = 4 * time.Hour

type entry struct {
	createAt time.Time
	value    interface{}
}

var (
	entries = map[string]*entry{}
	lock    sync.Mutex
)

func AddEntry(key string, value interface{}) {
	lock.Lock()
	defer lock.Unlock()
	redis2.Set(key, value, expireDuration)
}

func DelEntry(key string) {
	lock.Lock()
	defer lock.Unlock()
	redis2.Del(key)
}

func GetEntry(key string) interface{} {
	lock.Lock()
	defer lock.Unlock()
	v, err := redis2.Get(key)
	if err != nil {
		return nil
	}
	return v
}
