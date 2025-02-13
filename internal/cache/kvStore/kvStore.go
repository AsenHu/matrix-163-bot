package kvstore

import (
	"encoding/gob"
	"os"
	"sync"
)

type KVStore struct {
	mu    sync.RWMutex
	store map[interface{}]interface{}
}

func NewKVStore() *KVStore {
	return &KVStore{
		store: make(map[interface{}]interface{}),
	}
}

func (kv *KVStore) Set(key, value interface{}) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.store[key] = value
}

func (kv *KVStore) Get(key interface{}) (interface{}, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	val, ok := kv.store[key]
	return val, ok
}

func (kv *KVStore) Delete(key interface{}) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.store, key)
}

type KVPair struct {
	Key, Value interface{}
}

func (kv *KVStore) Iter() <-chan KVPair {
	ch := make(chan KVPair)
	go func() {
		kv.mu.RLock()
		defer kv.mu.RUnlock()
		for k, v := range kv.store {
			ch <- KVPair{k, v}
		}
		close(ch)
	}()
	return ch
}

func (kv *KVStore) Save(filename string) error {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	return enc.Encode(kv.store)
}

func (kv *KVStore) Load(filename string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	return dec.Decode(&kv.store)
}
