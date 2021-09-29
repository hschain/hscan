package server

import (
	"log"
	"sync"
)

const (
	CacheSize = 10 * 1024
)

type Cache struct {
	sync.RWMutex
	ids         []uint32
	index_all   map[string][]uint32
	index_all_1 map[string][]uint32
	index_all_2 map[string][]uint32

	ids_denom     map[string][]uint32
	index_denom   map[string]map[string][]uint32
	index_denom_1 map[string]map[string][]uint32
	index_denom_2 map[string]map[string][]uint32
}

func NewCache() *Cache {
	return &Cache{
		RWMutex: sync.RWMutex{},

		index_all:     make(map[string][]uint32),
		index_all_1:   make(map[string][]uint32),
		index_all_2:   make(map[string][]uint32),
		ids_denom:     make(map[string][]uint32),
		index_denom:   make(map[string]map[string][]uint32),
		index_denom_1: make(map[string]map[string][]uint32),
		index_denom_2: make(map[string]map[string][]uint32),
	}
}

func (c *Cache) Print() {
	length := len(c.ids)
	if length > 100 {
		length = 100
	}
	log.Printf("ids: %+v", c.ids[0:length])

}

func (c *Cache) Init(id uint32, index1, index2 string, denom string) {
	c.Lock()
	defer c.Unlock()

	c.ids = append(c.ids, id)

	c.index_all[index1] = append(c.index_all[index1], id)
	if index1 != index2 {
		c.index_all[index2] = append(c.index_all[index2], id)
	}

	c.index_all_1[index1] = append(c.index_all_1[index1], id)
	c.index_all_2[index2] = append(c.index_all_2[index2], id)

	c.ids_denom[denom] = append(c.ids_denom[denom], id)

	if _, ok := c.index_denom[denom]; !ok {
		c.index_denom[denom] = make(map[string][]uint32)
		c.index_denom_1[denom] = make(map[string][]uint32)
		c.index_denom_2[denom] = make(map[string][]uint32)
	}

	c.index_denom[denom][index1] = append(c.index_denom[denom][index1], id)
	if index1 != index2 {
		c.index_denom[denom][index2] = append(c.index_denom[denom][index2], id)
	}

	c.index_denom_1[denom][index1] = append(c.index_denom_1[denom][index1], id)
	c.index_denom_2[denom][index2] = append(c.index_denom_2[denom][index2], id)

}

func (c *Cache) Add(id uint32, index1, index2 string, denom string) {
	c.Lock()
	defer c.Unlock()

	c.ids = append([]uint32{id}, c.ids...)

	c.index_all[index1] = append([]uint32{id}, c.index_all[index1]...)
	if index1 != index2 {
		c.index_all[index2] = append([]uint32{id}, c.index_all[index2]...)
	}

	c.index_all_1[index1] = append([]uint32{id}, c.index_all_1[index1]...)
	c.index_all_2[index2] = append([]uint32{id}, c.index_all_2[index2]...)

	c.ids_denom[denom] = append([]uint32{id}, c.ids_denom[denom]...)
	if _, ok := c.index_denom[denom]; !ok {
		c.index_denom[denom] = make(map[string][]uint32)
		c.index_denom_1[denom] = make(map[string][]uint32)
		c.index_denom_2[denom] = make(map[string][]uint32)
	}

	c.index_denom[denom][index1] = append([]uint32{id}, c.index_denom[denom][index1]...)
	if index1 != index2 {
		c.index_denom[denom][index2] = append([]uint32{id}, c.index_denom[denom][index2]...)
	}

	c.index_denom_1[denom][index1] = append([]uint32{id}, c.index_denom_1[denom][index1]...)
	c.index_denom_2[denom][index2] = append([]uint32{id}, c.index_denom_2[denom][index2]...)

}

func (c *Cache) GetTotal(address string, denom string, typ int) int {
	c.RLock()
	defer c.RUnlock()

	if address == "null" {
		if denom == "null" {
			return len(c.ids)
		} else {
			return len(c.ids_denom[denom])
		}
	} else {
		if denom == "null" {
			if typ == 0 {
				return len(c.index_all[address])
			}
			if typ == 1 {
				return len(c.index_all_1[address])
			}
			if typ == 2 {
				return len(c.index_all_2[address])
			}
		} else {
			if _, ok := c.index_denom[denom]; !ok {
				return 0
			}
			if typ == 0 {
				return len(c.index_denom[denom][address])
			}
			if typ == 1 {
				return len(c.index_denom_1[denom][address])
			}
			if typ == 2 {
				return len(c.index_denom_2[denom][address])
			}
		}

	}
	return 0
}

func (c *Cache) GetTxids(address, denom string, typ int, page, iLimit int64) []uint32 {
	c.RLock()
	defer c.RUnlock()
	var ids []uint32
	if address == "null" {
		if denom == "null" {
			ids = c.ids
		} else {
			ids = c.ids_denom[denom]
		}
	} else {
		if denom == "null" {
			if typ == 0 {
				ids = c.index_all[address]
			}
			if typ == 1 {
				ids = c.index_all_1[address]
			}
			if typ == 2 {
				return c.index_all_2[address]
			}
		} else {
			if _, ok := c.index_denom[denom]; ok {
				if typ == 0 {
					ids = c.index_denom[denom][address]
				}
				if typ == 1 {
					ids = c.index_denom_1[denom][address]
				}
				if typ == 2 {
					ids = c.index_denom_2[denom][address]
				}
			}

		}
	}

	if (int64)(len(ids)) < iLimit*(page-1) {
		return nil
	}

	if (int64)(len(ids)) <= iLimit*page {
		return ids[iLimit*(page-1):]
	}

	return ids[iLimit*(page-1) : iLimit*page]

}
