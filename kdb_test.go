package kached

import (
	"strconv"
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	var err error
	testKey := `foo`
	testVal := `bar`
	invalidKey := `bam`

	kdb := testingKDB(t)
	defer func() {
		err = kdb.Close()
		if err != nil {
			t.Fatal("error closing kdb:", err)
		}
	}()
	err = kdb.Set(testKey, testVal, 1)
	if err != nil {
		t.Error("error setting value:", err)
		return
	}
	val, err := kdb.Get(testKey)
	if err != nil {
		t.Error("error getting value:", err)
		return
	}
	v, str := val.(string)
	if !str {
		t.Errorf("get value is not string, is %T\n", val)
		return
	}
	if v != testVal {
		t.Errorf("get expected value %v, got %v:", testVal, v)
		return
	}
	val, err = kdb.CacheGet(testKey)
	if err != nil {
		t.Error("cache get error getting value:", err)
		return
	}
	v, str = val.(string)
	if !str {
		t.Errorf("cached get value is not string, is %T\n", val)
		return
	}
	if v != testVal {
		t.Errorf("cached get expected value %v, got %v:", testVal, v)
		return
	}
	val, err = kdb.DBGet(testKey)
	if err != nil {
		t.Error("db get error getting value:", err)
		return
	}
	v, str = val.(string)
	if !str {
		t.Errorf("db get value is not string, is %T\n", val)
		return
	}
	if v != testVal {
		t.Errorf("db get expected value %v, got %v:", testVal, v)
		return
	}
	val, err = kdb.Get(invalidKey)
	if err == nil {
		t.Errorf("key %v should not exist!", invalidKey)
		return
	}
	val, err = kdb.DBGet(invalidKey)
	if err == nil {
		t.Errorf("key %v should not exist!", invalidKey)
		return
	}
	err = kdb.DBSet(`bypassCache`, 777)
	if err != nil {
		t.Errorf("error bypassing cache and setting db value: %v", err)
		return
	}
	val, err = kdb.CacheGet(`bypassCache`)
	if !IsErrCode(err, ErrNotFoundCache) {
		t.Errorf("cache not bypassed! Found Value: %v", val)
		return
	}
	val, err = kdb.DBGet(`bypassCache`)
	if err != nil {
		t.Errorf("error getting value from DB: %v", err)
		return
	}
	err = kdb.Delete(testKey)
	if err != nil {
		t.Errorf("error deleting key: %v", err)
		return
	}
	// wait for value to pass through buffers
	time.Sleep(10 * time.Millisecond)
	val, err = kdb.Get(testKey)
	if !IsErrCode(err, ErrNotFoundCacheOrDB) {
		t.Errorf("error code %v does not match %v", err, ErrNotFoundCacheOrDB)
		return
	}
}

func TestEvict(t *testing.T) {
	var del bool
	var count int
	evicted := make(map[uint64]interface{})
	config := NewConfig("./dbDir")
	if config == nil {
		t.Error("Config is NIL")
		return
	}
	config.Cache.NumCounters = 10
	config.Cache.MaxCost = 1
	config.Cache.OnEvict = func(key uint64, value interface{}, cost int64) {
		switch {
		case del:
			t.Logf("deleting key %v from eviction list\n", key)
			delete(evicted, key)
		default:
			if count < 10 {
				t.Logf("adding key %v to eviction list\n", key)
				evicted[key] = value
			}
		}
	}
	kdb, err := New(config)
	if err != nil {
		t.Error("error creating kdb:", err)
		return
	}
	if kdb == (&KDB{}) {
		t.Error("kdb is empty!")
		return
	}
	defer kdb.Close()
	for i := 0; i < 11; i++ {
		k := strconv.Itoa(i)
		err := kdb.Set(k, i, 1)
		if err != nil {
			t.Error("error setting value:", err)
			return
		}
		count++
	}
	// wait for value to pass through buffers
	time.Sleep(10 * time.Millisecond)
	if len(evicted) != 9 {
		t.Errorf("expected 9 keys evicted, found %v\n", len(evicted))
		return
	}
	del = true
	kdb.Set(`flushKey`, 777, 1)
	err = kdb.Delete(`flushKey`)
	if err != nil {
		t.Errorf("error deleting flushKey: %v", err)
		return
	}
	// wait for value to pass through buffers
	time.Sleep(10 * time.Millisecond)
	for i := 0; i < 11; i++ {
		k := strconv.Itoa(i)
		_, err := kdb.Get(k)
		if err != nil {
			t.Error("error getting value:", err)
			return
		}
	}
	// wait for value to pass through buffers
	time.Sleep(10 * time.Millisecond)
	if len(evicted) != 0 {
		t.Errorf("expected 0 keys left, found %v", len(evicted))
		return
	}
}

func testingKDB(t *testing.T) *KDB {
	config := NewConfig("./dbDir")
	if config == nil {
		t.Error("Config is NIL")
		t.FailNow()
	}
	config.Cache.NumCounters = 10
	config.Cache.MaxCost = 1
	kdb, err := New(config)
	if err != nil {
		t.Error("error creating kdb:", err)
		t.FailNow()
	}
	if kdb == (&KDB{}) {
		t.Error("kdb is empty!")
		t.FailNow()
	}
	return kdb
}
