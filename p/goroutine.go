package p

import (
	"sync"

	"github.com/petermattis/goid"
)

// GoID 获取当前 Goroutine 的 ID
func GoID() int {
	return int(goid.Get())
}

var _globals struct {
	ms map[int]rwmap
	sync.RWMutex
}

// G 获取当前协程内的全局变量
// 参考自 http://php.net/manual/zh/reserved.variables.globals.php
func G() *rwmap {
	_globals.Lock()

	if _globals.ms == nil {
		_globals.ms = make(map[int]rwmap, 32)
	}

	m, ok := _globals.ms[GoID()]
	if !ok {
		_globals.ms[GoID()] = *newRWMap()
		m = _globals.ms[GoID()]
	}

	_globals.Unlock()

	return &m
}

// _global alias to $_GLOBAL in current goroutine
type rwmap struct {
	data map[interface{}]interface{}
	rw   sync.RWMutex
}

func newRWMap() *rwmap {
	return &rwmap{data: make(map[interface{}]interface{}, 4)}
}

func (m *rwmap) Set(key interface{}, value interface{}) {
	m.rw.Lock()
	m.data[key] = value
	m.rw.Unlock()
}

func (m *rwmap) Get(key interface{}) interface{} {
	m.rw.Lock()
	defer m.rw.Unlock()
	return m.data[key]
}
