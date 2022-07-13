package cache

import (
	"encoding/json"
	"sync"
	"time"

	logger "github.com/NpoolPlatform/go-service-framework/pkg/logger"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
)

const expireDuration = 4 * time.Hour

var lock sync.Mutex

func AddEntry(key string, value interface{}) {
	lock.Lock()
	defer lock.Unlock()
	body, err := json.Marshal(value)
	if err != nil {
		logger.Sugar().Errorf("fail marshal %v value: %v", key, err)
		return
	}
	if err := redis2.Set(key, body, expireDuration); err != nil {
		logger.Sugar().Errorf("fail update cache %v: %v", key, err)
	}
}

func DelEntry(key string) {
	lock.Lock()
	defer lock.Unlock()
	if err := redis2.Del(key); err != nil {
		logger.Sugar().Errorf("fail del cache %v: %v", key, err)
	}
}

func GetEntry(key string, f func(data []byte) (interface{}, error)) interface{} {
	lock.Lock()
	defer lock.Unlock()
	v, err := redis2.Get(key)
	if err != nil {
		logger.Sugar().Errorf("fail get cache %v: %v", key, err)
		return nil
	}
	if v.(string) == "" {
		return nil
	}
	if f != nil {
		s, err := f([]byte(v.(string)))
		if err != nil {
			logger.Sugar().Errorf("fail unmarshal cache %v: %v", key, err)
			return nil
		}
		return s
	}
	return v
}
