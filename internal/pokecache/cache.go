package pokecache

import ("time"
				"sync"
				)

type CacheEntry struct {
	timeCreated time.Time
	val []byte
}

type Cache struct {
	vals map[string]CacheEntry
	interval time.Duration
	mu   sync.Mutex
}


func NewCache(interval time.Duration) Cache {
	result := Cache{ vals : make(map[string]CacheEntry),
									 interval : interval,
  }
	go result.reapLoop()
	return result
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

  newEntry := CacheEntry{timeCreated : time.Now(),
                         val : val,
                        }
  c.vals[key] = newEntry
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	result, exists := c.vals[key]
	if !exists {
		return []byte{}, exists
	}
	return result.val, exists
}

func (c Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for tick := range ticker.C {
		for key, val := range c.vals {
			if tick.Sub(val.timeCreated) > c.interval {
				c.mu.Lock()
				delete(c.vals, key)
				c.mu.Unlock()
			}
		}
	}
}

