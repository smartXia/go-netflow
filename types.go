package netflow

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const (
	TCP_ESTABLISHED = iota + 1
	TCP_SYN_SENT
	TCP_SYN_RECV
	TCP_FIN_WAIT1
	TCP_FIN_WAIT2
	TCP_TIME_WAIT
	TCP_CLOSE
	TCP_CLOSE_WAIT
	TCP_LAST_ACK
	TCP_LISTEN
	TCP_CLOSING
	//TCP_NEW_SYN_RECV
	//TCP_MAX_STATES
)

var states = map[int]string{
	TCP_ESTABLISHED: "ESTABLISHED",
	TCP_SYN_SENT:    "SYN_SENT",
	TCP_SYN_RECV:    "SYN_RECV",
	TCP_FIN_WAIT1:   "FIN_WAIT1",
	TCP_FIN_WAIT2:   "FIN_WAIT2",
	TCP_TIME_WAIT:   "TIME_WAIT",
	TCP_CLOSE:       "CLOSE",
	TCP_CLOSE_WAIT:  "CLOSE_WAIT",
	TCP_LAST_ACK:    "LAST_ACK",
	TCP_LISTEN:      "LISTEN",
	TCP_CLOSING:     "CLOSING",
	//TCP_NEW_SYN_RECV: "NEW_SYN_RECV",
	//TCP_MAX_STATES:   "MAX_STATES",
}

// https://github.com/torvalds/linux/blob/master/include/net/tcp_states.h
var StateMapping = map[string]string{
	"01": "ESTABLISHED",
	"02": "SYN_SENT",
	"03": "SYN_RECV",
	"04": "FIN_WAIT1",
	"05": "FIN_WAIT2",
	"06": "TIME_WAIT",
	"07": "CLOSE",
	"08": "CLOSE_WAIT",
	"09": "LAST_ACK",
	"0A": "LISTEN",
	"0B": "CLOSING",
}

//type Mapping struct {
//	cb   func()
//	dict map[string]string
//	sync.RWMutex
//}

type Mapping struct {
	cb   func()
	dict map[string]*entry
	sync.RWMutex
}
type entry struct {
	value     string
	timestamp time.Time
}

func (m *Mapping) Length() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.dict)
}
func NewMapping() *Mapping {
	size := 1000
	return &Mapping{
		dict: make(map[string]*entry, size),
	}
}

func (m *Mapping) Handle() {
}

// Add 向 Mapping 中添加键值对，并记录时间戳
func (m *Mapping) Add(key string, value string) {
	m.Lock()
	defer m.Unlock()
	m.dict[key] = &entry{
		value:     value,
		timestamp: time.Now(),
	}
}

// Cleanup 清理过期的条目
func (m *Mapping) Cleanup(expiration time.Duration) {
	m.Lock()
	defer m.Unlock()
	now := time.Now()
	println("清理数据")
	for key, e := range m.dict {
		if now.Sub(e.timestamp) > expiration {
			delete(m.dict, key)
		}
	}
}
func (m *Mapping) Exists(key, value string) bool {
	m.RLock()
	defer m.RUnlock()
	v, exists := m.dict[key]
	return exists && v.value == value
}

func (m *Mapping) Get(key string) (string, bool) {
	m.RLock()
	defer m.RUnlock()
	ent, exists := m.dict[key]
	if !exists {
		return "", false
	}
	return ent.value, true
}
func (m *Mapping) Delete(k string) {
	m.Lock()
	defer m.Unlock()

	delete(m.dict, k)
}

func (m *Mapping) String() string {
	m.RLock()
	defer m.RUnlock()

	bs, _ := json.Marshal(m.dict)
	return string(bs)
}

type Null struct{}

type LoggerInterface interface {
	Debug(...interface{})
	Error(...interface{})
}

var defaultLogger string

type logger struct{}

func (l *logger) Debug(msg ...interface{}) {
	fmt.Println(msg...)
}

func (l *logger) Error(msg ...interface{}) {
	fmt.Println(msg...)
}
