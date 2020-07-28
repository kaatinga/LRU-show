package LRU

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewCache(t *testing.T) {
	var cacheSize1 byte = 5
	var cacheSize2 byte = 0
	tests := []struct {
		name      string
		cacheSize byte
		want      *Cache
		wantErr   bool
	}{
		{"ok", cacheSize1, &Cache{
			mx:       sync.RWMutex{},
			items:    make(map[string]*item, cacheSize1),
			size:     0,
			capacity: cacheSize1,
			order:    order{},
		}, false},
		{"!ok", cacheSize2, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCache(tt.cacheSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}
