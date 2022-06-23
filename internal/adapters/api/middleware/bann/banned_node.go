package bann

// head --> ... --> prev --> current --> next --> ... --> tail

import (
	"sync"
	"time"

	"github.com/giffone/forum-security/internal/config"
)

type banNode struct {
	prev *banNode
	next *banNode
	key  string
}

type bannedList struct {
	sync.RWMutex
	head   *banNode
	tail   *banNode
	length int
}

func (bl *bannedList) insert(key string) {
	bl.Lock()
	defer bl.Unlock()
	newNode := &banNode{
		next: bl.head, // for new make next - current(old)
		key:  key,
	}
	// for current (old) make prev - new
	if bl.head != nil {
		bl.head.prev = newNode
	}
	// and write to head new
	bl.head = newNode
	// make tail
	if bl.tail == nil && bl.length == 0 { // first initialize
		bl.tail = newNode
	}
	bl.length++
}

func (bl *bannedList) append(newNode *banNode) {
	bl.Lock()
	defer bl.Unlock()
	newNode.next = bl.head // for new make next - current(old)
	// for current (old) make prev - new
	if bl.head != nil {
		bl.head.prev = newNode
	}
	// and write to head new
	bl.head = newNode
	bl.length++
}

func (bl *bannedList) findTail() {
	bl.RLock()
	if bl.tail != nil {
		bl.RUnlock()
		return
	}
	bl.RUnlock()
	bl.Lock()
	defer bl.Unlock()
	b := bl.head
	bl.length = 1
	for b.next != nil {
		b = b.next
		bl.length++
	}
	bl.tail = b
}

// cutTail will walk from tail
func (bl *bannedList) cutTail(old map[int]*banNode) {
	bl.RLock()
	length := bl.length
	bl.RUnlock()
	// re-check if another goroutine already cleared
	if length < config.NodeMaxLength {
		return
	}
	end := time.Now().Add(config.ClearDuration)
	bl.Lock()
	defer bl.Unlock()
	for i := 0; bl.tail.prev != nil && i < config.FindOldNodeLoop; i++ {
		// move tail (delete) and appen in to slice
		old[i] = bl.tail
		bl.tail = bl.tail.prev
		bl.length--
		// check deadline in 16th iteration
		if i&0x0f == 0 && time.Now().After(end) { // i&0x0f = i%16 = 0
			break
		}
	}
}

func (bl *bannedList) restore(old map[int]*banNode) {
	for _, v := range old {
		bl.append(v)
	}
}
