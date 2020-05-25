/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package lru implements an LRU cache.
package lru

import "container/list"

// Cache is an LRU cache. It is not safe for concurrent access.
// Cache结构用于实现LRU cache算法；并发访问不安全
type Cache struct {
	// MaxEntries is the maximum number of cache entries before
	// an item is evicted. Zero means no limit.
	// 最大入口数，也就是缓存中最多存几条数据，超过了就触发数据淘汰；0表示没有限制
	MaxEntries int

	// OnEvicted optionally specifies a callback function to be
	// executed when an entry is purged from the cache.
	// 销毁前回调
	OnEvicted func(key Key, value interface{})

	ll *list.List
	// key为任意类型，值为指向链表一个结点的指针
	cache map[interface{}]*list.Element
}

// A Key may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
// 任意可比较类型
type Key interface{}

// 访问入口结构，包装键值
type entry struct {
	key   Key
	value interface{}
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
// 初始化一个Cache类型实例
func New(maxEntries int) *Cache {
	return &Cache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Element),
	}
}

// Add adds a value to the cache.
// 往缓存中增加一个值
func (c *Cache) Add(key Key, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}

	// 如果key已经存在，则将记录前移到头部，然后设置value
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		return
	}
	// key不存在时，创建一条记录，插入链表头部，ele是这个Element的指针
	// 这里的Element是一个*entry类型，ele是*list.Element类型
	ele := c.ll.PushFront(&entry{key, value})
	// cache这个map设置key为Key类型的key，value为*list.Element类型的ele
	c.cache[key] = ele
	// 链表长度超过最大入口值，触发清理操作
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.RemoveOldest()
	}
}

// Get looks up a key's value from the cache.
// 链表长度超过最大入口值，触发清理操作
func (c *Cache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	// 如果存在
	if ele, hit := c.cache[key]; hit {
		// 将这个Element移动到链表头部
		c.ll.MoveToFront(ele)
		// 返回entry的值
		return ele.Value.(*entry).value, true
	}
	return
}

// Remove removes the provided key from the cache.
// 如果key存在，调用removeElement删除链表and缓存中的元素
func (c *Cache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
// 删除最旧的元素
func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *list.Element) {
	// 链表中删除一个element
	c.ll.Remove(e)
	// e.Value本质是*entry类型，entry结构体就包含了key和value2个属性
	// Value本身是interface{}类型，通过类型断言转成*entry类型
	kv := e.Value.(*entry)
	// 删除cache这个map中key为kv.key这个元素；也就是链表中删了之后缓存中也得删
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

// Len returns the number of items in the cache.
// 返回缓存中的item数，通过链表的Len()方法获取
func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Clear purges all stored items from the cache.
// 删除缓存中所有条目，如果有回调函数OnEvicted()，则先调用所有回调函数，然后置空
func (c *Cache) Clear() {
	if c.OnEvicted != nil {
		for _, e := range c.cache {
			kv := e.Value.(*entry)
			c.OnEvicted(kv.key, kv.value)
		}
	}
	c.ll = nil
	c.cache = nil
}
