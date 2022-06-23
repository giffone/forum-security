package middleware

import "testing"

func TestUnregAdd(t *testing.T) {
	// u := newUnreg()
	// schemaIP := "127.0.0." // + 0... n
	// n := 10000             // add ip between 0 and "n" like "127.0.0.n"
	// err := b.tBannedAddRandom(schemaIP, n)
	// if err != nil {
	// 	t.Errorf("test add: %s", err.Error())
	// }
	// // if need result save to file
	// b.saveMap()
}

func newUnreg() *unRegistered {
	return &unRegistered{
		ip: make(map[string]*requests),
	}
}
