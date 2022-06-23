package bann

import (
	"time"

	"github.com/giffone/forum-security/internal/adapters/api"
)

type banned struct {
	bm bannedMap
	bl bannedList
}

func NewBanned() api.Banned {
	return &banned{
		bm: bannedMap{
			banned: make(map[string]time.Time),
		},
		bl: bannedList{},
	}
}

func (b *banned) Add(key string, banTill time.Time) {
	b.bm.add(key, banTill)
	go b.bl.insert(key)
}

func (b *banned) UnlockBan(deadline time.Time, key string) (exist bool, expire bool) {
	if exist, expire = b.bm.unlock(deadline, key); exist && expire {
		// fix tail if need
		b.bl.findTail()
		// make slice to append old nodes
		old := make(map[int]*banNode)
		// cut nodes from tail (old records)
		b.bl.cutTail(old)
		// check slice for old records
		now := time.Now()
		for i, v := range old {
			if exist, expire := b.bm.unlock(now, v.key); expire || !exist {
				delete(old, i)
			}
		}
		// after delete complete, remaining (actual) restore in list
		go b.bl.restore(old)
	}
	return
}
