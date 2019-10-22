package kached

import (
	"testing"
)

func TestCodeStrings(t *testing.T) {
	for code := range ErrCodeStrings {
		c := ErrCode(code)
		err := c
		if !IsErrCode(err, c) {
			t.Errorf("%v does not match %v", err, c)
			t.FailNow()
		}
		errMsg := ErrMSG{
			code:  c,
			stack: `More Details Here`,
		}
		if !IsErrCode(errMsg, c) {
			t.Errorf("%v does not match %v", err, c)
			t.FailNow()
		}
		str1 := errMsg.Error()
		errMsg.msg("Edited Message Here")
		str2 := errMsg.Error()
		if str1 == str2 {
			t.Errorf("failed to edit err stack message")
			t.FailNow()
		}
	}
}

func TestErrNotFound(t *testing.T) {
	kdb := testingKDB(t)
	defer kdb.Close()
	val, err := kdb.CacheGet("testKey")
	if err != ErrNotFoundCache {
		t.Fatalf("expected [ErrNotFoundCache] error for non existent key: [testKey] in cache, received error: [%v], found value %v", err, val)
	}
	if !IsErrCode(err, ErrNotFoundCache) {
		t.Fatalf("%v does not match %v", err, ErrNotFoundCache)
	}
	val, err = kdb.DBGet("testKey")
	if err != ErrNotFoundDB {
		t.Fatalf("expected [ErrNotFoundDB] error for non existent key: [testKey] in DB, received error: [%v], found value %v", err, val)
	}
	if !IsErrCode(err, ErrNotFoundDB) {
		t.Fatalf("%v does not match %v", err, ErrNotFoundDB)
	}
	val, err = kdb.Get("testKey")
	if err != ErrNotFoundCacheOrDB {
		t.Fatalf("expected [ErrNotFoundCacheOrDB] error for non existent key: [testKey] in cache and DB, received error: [%v], found value %v", err, val)
	}
	if !IsErrCode(err, ErrNotFoundCacheOrDB) {
		t.Fatalf("%v does not match %v", err, ErrNotFoundCacheOrDB)
	}
}
