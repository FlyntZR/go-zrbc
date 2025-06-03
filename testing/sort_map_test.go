package main

import (
	"sort"
	"testing"
)

func TestSortMap(t *testing.T) {
	m := map[int]string{
		3: "three",
		1: "one",
		4: "four",
		2: "two",
	}

	// 从 map 中获取所有的键
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// 对键进行排序
	sort.Ints(keys)

	// 按顺序遍历有序集合
	for _, k := range keys {
		v := m[k]
		t.Logf("k:(%d), v:(%s)", k, v)
	}
}
