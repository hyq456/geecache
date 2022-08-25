package test

import (
	"fmt"
	"geecache/geecache"
	"log"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadCount := make(map[string]int, len(db))
	gee := geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCount[key]; !ok {
					loadCount[key] = 0
				}
				loadCount[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exit", key)
		}))
	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}
		if _, err := gee.Get(k); err != nil || loadCount[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
