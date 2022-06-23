package bann

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/giffone/forum-security/internal/config"
)

// to change added time to ban_list
const (
	minus2hours = -time.Hour * 2
	plus2hours  = time.Hour * 2
	plus2month  = time.Hour * 24 * 31 * 2
)

type Banned interface { // origin interface ../forum-security/internal/adapters/api
	Add(key string, banTill time.Time)
	UnlockBan(deadline time.Time, key string) (exist bool, expire bool)
	// testing func
	tBannedAddRandom(addIP string, n int, till time.Time) error
	tBannedCheckUnlock(addedIP string, toCheck int, deadline time.Time, res *resultBanned)
	saveMap() (int, error)
	saveList() (int, error)
	lengthListFixed() int
	lengthMap() int
	saveFile(t *testing.T, n int, nameTest string)
}

func tNewBanned() Banned {
	return &banned{
		bm: bannedMap{
			banned: make(map[string]time.Time),
		},
		bl: bannedList{},
	}
}

func TestBannedAdd(t *testing.T) {
	// will add 10k ips
	add := 10000 // add ip between 0 and "n" like "127.0.0.n"
	b := tNewBanned()
	schemaIP := "127.0.0." // + 0... n
	err := b.tBannedAddRandom(schemaIP, add, time.Now())
	if err != nil {
		t.Errorf("add: %s", err.Error())
	}
	// if need result save to file
	b.saveFile(t, add, "add:")
}

func TestBannedAddRead(t *testing.T) {
	// expected
	result := struct{
		total int
	}{
		total : 20000,
	}
	// will add 10k (2 goroutines: check 10k/2 + add 10k) = actual ban is 20k in saved files
	add := 10000       // add ip between 0 and "n" like "127.0.0.n"
	toCheck := add / 2 // check ip for unlock
	banTill := time.Now().Add(config.BanDuration)
	b := tNewBanned()
	schemaIP := "127.0.0." // + 0... n

	err := b.tBannedAddRandom(schemaIP, add, banTill)
	if err != nil {
		t.Errorf("add+(add+read): add: %s", err.Error())
	}

	res := new(resultBanned)
	c := make(chan error)

	// read recently added
	go func(ch chan error) {
		b.tBannedCheckUnlock(schemaIP, toCheck, time.Now(), res)
		if res.readed != toCheck {
			ch <- fmt.Errorf("add+(add+read): read: expected %d but got %d; deleted old %d, not found %d", toCheck, res.readed, res.deleted, res.lost)
			// if some deleted - check function that expires date
		} else {
			ch <- nil
		}
	}(c)

	// write new ip while read
	newSchemaIP := "127.1.0."
	go func(ch chan error) {
		err := b.tBannedAddRandom(newSchemaIP, add, time.Now())
		if err != nil {
			ch <- fmt.Errorf("add+(add+read): add2: %s", err.Error())
		} else {
			ch <- nil
		}
	}(c)

	err1 := <-c
	err2 := <-c
	if err1 != nil || err2 != nil {
		t.Errorf("err1: %v\nerr2: %v", err1, err2)
	}
	lBanned := b.lengthMap()
	if lBanned != add*2 {
		t.Errorf("add+(add+read): added %d but got %d", add*2, lBanned)
	}
	// if need result save to file
	b.saveFile(t, add*2, "add+(add+read):")
}

func TestBannedAddDelete(t *testing.T) {
	// need to add ips less than "NodeMaxLength" to not activate auto-clear for old records in lists = 10k/3 = 3,3k
	// will add 3,3k ip (2 goroutines: check and delete 3,3k/2 + add 3,3k) = actual ban is 20k in saved files
	add := config.NodeMaxLength / 3 // add ip between 0 and "n" like "127.0.0.n"
	toCheck := add / 2
	banTill := time.Now().Add(-config.BanDuration) // -1h to delete some bans
	banTill2 := time.Now().Add(config.BanDuration) // +1h to not delete
	b := tNewBanned()
	schemaIP := "127.0.0." // + 0... n

	err := b.tBannedAddRandom(schemaIP, add, banTill)
	if err != nil {
		t.Errorf("test add+(add+read+delete): add: %s", err.Error())
	}

	res := new(resultBanned)
	c := make(chan error)

	b.UnlockBan(time.Now(), "127.0.0.34")

	// delete recently added
	go func(ch chan error) {
		b.tBannedCheckUnlock(schemaIP, toCheck, time.Now(), res)
		if err != nil {
			ch <- fmt.Errorf("test add+(add+read+delete): read+delete: %s", err.Error())
		} else {
			ch <- nil
		}
	}(c)

	// write new ip while delete
	newSchemaIP := "127.1.0."
	go func(ch chan error) {
		err := b.tBannedAddRandom(newSchemaIP, add, banTill2)
		if err != nil {
			ch <- fmt.Errorf("test add+(add+read+delete): add2: %s", err.Error())
		} else {
			ch <- nil
		}
	}(c)

	err1 := <-c
	err2 := <-c
	if err1 != nil || err2 != nil {
		t.Errorf("err1: %v\nerr2: %v", err1, err2)
	}
	lBanned := b.lengthMap()
	added := add * 2
	if lBanned != added-toCheck {
		t.Errorf("added %d and deleted %d, should stay %d but got %d", added, toCheck, added-toCheck, lBanned)
	}
	// if need result save to file
	b.saveFile(t, added-toCheck, "test add+(add+read+delete):")
}

// func TestBannedNotExist(t *testing.T) {
// 	// b := newBanned()
// 	schemaIP := "127.0.0." // + 0... n
// 	n := 10000             // add ip between 0 and "n" like "127.0.0.n"
// 	err := b.tBannedAddRandom(schemaIP, n)
// 	if err != nil {
// 		t.Errorf("test add: %s", err.Error())
// 	}

// 	NoneExistant := "127.1.0.0"
// 	exist, expire := b.unlock(time.Now().Add(time.Hour*1), NoneExistant) // add 1 hour to "expireDate" to delete entry
// 	if exist || expire {
// 		t.Errorf("none-existant ip shows exist-%t and expire-%t", exist, expire)
// 	}
// 	lBanned := len(b.banned)
// 	if lBanned != n {
// 		t.Errorf("addded %d but got %d", n, lBanned)
// 	}
// }

func (b *banned) tBannedAddRandom(addIP string, n int, banTill time.Time) error {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(random int) {
			ip := fmt.Sprintf("%s%d", addIP, random)
			b.Add(ip, banTill)
			wg.Done()
		}(i)
	}
	wg.Wait()
	lBanned := len(b.bm.banned)
	if lBanned < n {
		return fmt.Errorf("addded %d but got %d", n, lBanned)
	}
	return nil
}

func (b *banned) lengthListFixed() int {
	return b.bl.length
}

func (b *banned) lengthMap() int {
	return len(b.bm.banned)
}

func (b *banned) tBannedCheckUnlock(addedIP string, toCheck int, deadline time.Time, res *resultBanned) {
	var wg sync.WaitGroup

	for i := 0; i < toCheck; i++ {
		wg.Add(1)
		go func(random int) {
			ip := fmt.Sprintf("%s%d", addedIP, random)
			exist, expire := b.UnlockBan(deadline, ip)
			if exist && !expire {
				res.add(readed)
			} else if exist && expire {
				res.add(deleted)
			} else {
				res.add(lost)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func (b *banned) tBannedReadAndDelete(addedIP string, n int) error {
	var wg sync.WaitGroup
	res := resultBanned{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(random int) {
			ip := fmt.Sprintf("%s%d", addedIP, random)
			exist, expire := b.UnlockBan(time.Now().Add(time.Hour*1), ip)
			if exist && expire {
				res.add(deleted)
			} else if !exist {
				res.add(lost)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	if res.deleted != n {
		return fmt.Errorf("read_and_delete: deleted %d but expected %d, not found %d", res.deleted, n, res.lost)
	}
	return nil
}

func (b *banned) saveMap() (int, error) {
	count := 0
	buf := bytes.Buffer{}
	buf.WriteByte('[')
	for ip, date := range b.bm.banned {
		buf.WriteString(fmt.Sprintf("{\"%s\": \"%v\"},\n", ip, date))
		count++
	}
	buf.WriteByte(']')

	os.Remove("ips_map.json")
	file, err := os.Create("ips_map.json")
	if err != nil {
		return 0, err
	}
	defer file.Close()
	_, err = buf.WriteTo(file)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b *banned) saveList() (int, error) {
	count := 0
	buf := bytes.Buffer{}
	l := b.bl.head
	buf.WriteByte('[')
	for l != nil {
		buf.WriteString(fmt.Sprintf("{\"key\": \"%s\"},\n", l.key))
		l = l.next
		count++
	}
	buf.WriteByte(']')

	os.Remove("ips_list.json")
	file, err := os.Create("ips_list.json")
	if err != nil {
		return 0, err
	}
	defer file.Close()
	_, err = buf.WriteTo(file)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b *banned) saveFile(t *testing.T, n int, nameTest string) {
	nCount, err := b.saveMap()
	if err != nil {
		t.Errorf("%s saveMap: %s", nameTest, err.Error())
	}
	if nCount != n {
		t.Errorf("%s added %d list but saved to map_file %d", nameTest, n, nCount)
	}

	nCount, err = b.saveList()
	if err != nil {
		t.Errorf("%s saveList: %s", nameTest, err.Error())
	}
	if nCount != n {
		t.Errorf("%s added %d list but saved to list_file %d", nameTest, n, nCount)
	}
	if b.lengthListFixed() != n {
		t.Errorf("%s added %d list but list length %d", nameTest, n, b.lengthListFixed())
	}
}

// to count the number of records (goroutines)
type resultBanned struct {
	sync.Mutex
	readed, deleted, lost int
}

const (
	readed  = "r"
	deleted = "d"
	lost    = "l"
)

func (rb *resultBanned) add(who string) {
	rb.Lock()
	defer rb.Unlock()
	switch who {
	case readed:
		rb.readed++
	case deleted:
		rb.deleted++
	case lost:
		rb.lost++
	}
}
