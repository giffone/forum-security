package bann

import (
	"sync"
	"time"
)

type bannedMap struct {
	sync.RWMutex
	banned map[string]time.Time // list of banned user_id/ip and expire date
}

// Add adds user_id/ip in ban list and adds expire date
func (bm *bannedMap) add(key string, banTill time.Time) {
	bm.Lock()
	defer bm.Unlock()
	bm.banned[key] = banTill
}

// Unlock checks if user_id/ip in actual banned list
// and can unlock if unlock time expired < "deadline"
func (bm *bannedMap) unlock(deadline time.Time, key string) (exist bool, expire bool) {
	bm.RLock()
	expireDate, exist := bm.banned[key]
	bm.RUnlock()
	if exist {
		// can delete from ban_list
		if expire = expireDate.Before(deadline); expire {
			bm.delete(key)
			return exist, expire // (true, true) exist in list, can unlock
		}
		return exist, expire // (true, false) exist in list, can not unlock
	}
	return exist, expire // (false, false) not exist in list and can not unlock
}

// delete deletes from banned list by "key" = user_id/ip
func (bm *bannedMap) delete(key string) {
	bm.Lock()
	defer bm.Unlock()
	delete(bm.banned, key)
}
