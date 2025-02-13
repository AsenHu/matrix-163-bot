package manager

import (
	"log"
	"os"
	"sync"
	"time"

	kvstore "matrix-163-bot/internal/cache/kvStore"
)

type Cache interface {
	IsExpired() bool
}

type CacheManager struct {
	Duration time.Duration
	Kv       *kvstore.KVStore
	Logger   Logger
	stopChan chan struct{}
	wg       sync.WaitGroup
}

type Logger interface {
	Printf(format string, v ...interface{})
}

func NewCacheManager(duration time.Duration, kv *kvstore.KVStore, logger Logger) *CacheManager {
	if logger == nil {
		logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	}
	return &CacheManager{
		Duration: duration,
		Kv:       kv,
		Logger:   logger,
		stopChan: make(chan struct{}),
	}
}

func (c *CacheManager) Run() {
	c.wg.Add(1)
	defer c.wg.Done()
	ticker := time.NewTicker(c.Duration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			go c.deleteExpired()
		case <-c.stopChan:
			c.Logger.Printf("CacheManager stopped.\n")
			return
		}
	}
}

func (c *CacheManager) Stop() {
	close(c.stopChan)
	c.Logger.Printf("Waiting for all goroutines to exit...\n")
	c.wg.Wait()
	c.Logger.Printf("GoodBye\n")
}

func (c *CacheManager) deleteExpired() {
	c.wg.Add(1)
	defer c.wg.Done()
	for kvItem := range c.Kv.Iter() {
		cacheItem, ok := kvItem.Value.(Cache)
		if !ok {
			c.Logger.Printf("Key %v does not implement Cache interface, skipping...\n", kvItem.Key)
			continue
		}
		if cacheItem.IsExpired() {
			c.Logger.Printf("Key %v has expired, deleting...\n", kvItem.Key)
			go c.Kv.Delete(kvItem.Key)
		}
	}
}
