package excelize

import (
	"container/list"
	"sync"
)

type formulaResult struct {
	value string
	err   error
}

type Cache struct {
	cache    map[string]*list.Element
	list     *list.List
	mutex    sync.RWMutex
	limit    int
	disabled bool
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]*list.Element),
		list:  list.New(),
		limit: 0, // Unlimited entries by default
	}
}

func (c *Cache) SetLimit(limit int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.limit = limit

	// If the limit is reduced, evict excess entries
	if limit > 0 {
		for c.list.Len() > limit {
			oldest := c.list.Back()
			if oldest != nil {
				cellRef := oldest.Value.(string)
				delete(c.cache, cellRef)
				c.list.Remove(oldest)
			}
		}
	}
}

func (c *Cache) Add(key string, value formulaResult) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.disabled {
		return
	}

	// If the cache has a positive limit and is full, evict the least recently used entry
	if c.limit > 0 && c.list.Len() >= c.limit {
		oldest := c.list.Back()
		if oldest != nil {
			// Get the key of the oldest entry
			for k, v := range c.cache {
				if v == oldest {
					delete(c.cache, k) // Remove the oldest entry from the cache
					break
				}
			}
			c.list.Remove(oldest) // Remove the oldest entry from the list
		}
	}

	// Add the new entry to the cache and the front of the list
	elem := c.list.PushFront(value)
	c.cache[key] = elem
}

func (c *Cache) Get(key string) (formulaResult, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.disabled {
		return formulaResult{}, false
	}

	if elem, ok := c.cache[key]; ok {
		// Move the accessed element to the front of the list (LRU)
		c.list.MoveToFront(elem)
		return elem.Value.(formulaResult), true
	}
	return formulaResult{}, false
}

func (c *Cache) DisableCache() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.disabled = true
	c.cache = make(map[string]*list.Element)
	c.list = list.New()
}

func (c *Cache) Invalidate() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// If the cache is disabled, do nothing
	if c.disabled {
		return
	}

	// Clear the cache map and reset the LRU list
	c.cache = make(map[string]*list.Element)
	c.list = list.New()
}
