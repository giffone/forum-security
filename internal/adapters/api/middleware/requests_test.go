package middleware

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/giffone/forum-security/internal/config"
)

func TestRequestsIsDDos(t *testing.T) {
	r := new(requests)
	now := time.Now()
	for i := 0; i < config.CheckAfter; i++ {
		r.req[i] = now
		now = now.Add((config.FrequencyRequest + (time.Millisecond * 100)) * 1) // if freq == 600ms, adding +100ms = 700ms (not often than limit)
		fmt.Println(r.req[i])
		r.count++
	}

	got, interval := r.isDDos()
	if got {
		t.Errorf("limit time is %v for %d requests and got %v for %d requests",
			config.CheckTimeInterval.Seconds(), config.CheckAfter, interval, config.CheckAfter)
	}

	os.Remove("req.json")
	r.saveArr(interval, got, false)

	r = new(requests)

	now = time.Now()
	for i := 0; i < config.CheckAfter; i++ {
		r.req[i] = now
		now = now.Add((config.FrequencyRequest - (time.Millisecond * 100)) * 1) // if freq == 600ms, substract -100ms = 500ms (often than limit)
		fmt.Println(r.req[i])
		r.count++
	}

	got, interval = r.isDDos()
	if !got {
		t.Errorf("limit time is %v for %d requests and got %v for %d requests",
			config.CheckTimeInterval.Seconds(), config.CheckAfter, interval, config.CheckAfter)
	}

	// for save results into file
	r.saveArr(interval, got, true)
}

func (r *requests) saveArr(interval float64, result, close bool) error {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("[{\"count\": \"%v\"},\n", r.count))
	buf.WriteString(fmt.Sprintf("{\"limit_sec\": \"%v\"},\n", config.CheckTimeInterval.Seconds()))
	buf.WriteString(fmt.Sprintf("{\"interval_sec\": \"%v\"},\n", interval))
	buf.WriteString(fmt.Sprintf("{\"banned\": \"%t\"},\n", result))
	for i, date := range r.req {
		buf.WriteByte('{')
		buf.WriteString(fmt.Sprintf("\"%d\": \"%v\"", i, date))
		buf.WriteString("},\n")
	}
	if close {
		buf.WriteByte(']')
	} else {
		buf.WriteByte('\n')
	}

	file, err := os.OpenFile("req.json", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
