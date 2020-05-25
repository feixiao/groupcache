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

// Package consistenthash provides an implementation of a ring hash.
package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash函数
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int   // 每个key的副本数量
	keys     []int // Sorted，key为哈希环上面的一个点(节点哈希值)
	// hashMap表的key是一个表示cache服务器或者副本的hash值，value为一个具体的cache服务器，
	// 这样就完成了Cache A、Cache A1、Cache A2等副本全部映射到Cache A的功能。
	hashMap map[int]string // 哈希环上面的一个点到服务器名的映射
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		// 默认的hash函数
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// IsEmpty returns true if there are no items available.
func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

// Add adds some keys to the hash.
// 将缓存服务器加到Map中，比如Cache A、Cache B作为keys，
// 如果副本数指定的是2，那么Map中存的数据是Cache A#1、Cache A#2、Cache B#1、Cache B#2的hash结果
// keys类似"host1"，"host2"等，表示cache服务器
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 内循环实现副本数量要求
		for i := 0; i < m.replicas; i++ {
			// key和数字编号一起计算哈希
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// 哈希值存入切片
			m.keys = append(m.keys, hash)
			// 哈希值和服务器名字的对应关系
			m.hashMap[hash] = key
		}
	}

	// 从小到大排序
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key.
// 如果有一个key要保存到某个cache服务器，Get函数返回对应的cache服务器。
func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}

	// 计算哈希值
	hash := int(m.hash([]byte(key)))

	// Binary search for appropriate replica.
	// 查找m.keys[i] >= hash成立的最小值i，i前面的元素都不满足>hash
	idx := sort.Search(len(m.keys), func(i int) bool { return m.keys[i] >= hash })

	// Means we have cycled back to the first replica.
	if idx == len(m.keys) {
		idx = 0
	}

	// 返回对应的服务器信息
	return m.hashMap[m.keys[idx]]
}
