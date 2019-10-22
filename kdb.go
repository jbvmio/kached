package kached

import (
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/ristretto"
)

// ItemCost is used as the default cost when adding an item to cache.
var ItemCost int64 = 1

// KDB is a cached KV database.
type KDB struct {
	Encoder  EncodeFunc
	Decoder  DecodeFunc
	cache    *ristretto.Cache
	database *badger.DB
}

// New returns a new KDB using the given config.
func New(config *Config) (*KDB, error) {
	c, err := ristretto.NewCache(config.Cache)
	if err != nil {
		return &KDB{}, err
	}
	d, err := badger.Open(config.Database)
	if err != nil {
		return &KDB{}, err
	}
	return &KDB{
		Encoder:  encode,
		Decoder:  decode,
		cache:    c,
		database: d,
	}, nil
}

// Set inserts a key value pair into cache and the underlying DB.
func (k *KDB) Set(key, value interface{}, cost ...int64) error {
	var C int64
	var err ErrCode
	var errMsg string
	switch len(cost) {
	case 0:
		C = ItemCost
	default:
		C = cost[0]
	}
	ok := k.cache.Set(key, value, C)
	if !ok {
		err = ErrUnableToCache
	}
	errd := k.database.Update(func(txn *badger.Txn) error {
		key := k.Encoder(key)
		val := k.Encoder(value)
		e := badger.NewEntry(key, val)
		err := txn.SetEntry(e)
		return err
	})
	if errd != nil {
		err = ErrUnableToCacheOrSave
		errMsg = errd.Error()
	}
	if err == ErrNoErr {
		return nil
	}
	return ErrMSG{
		code:  err,
		stack: errMsg,
	}
}

// Get retrieves a value by key from cache, falling back to the DB if needed. It will reload the key into cache if found in the DB.
func (k *KDB) Get(key interface{}) (value interface{}, err error) {
	value, there := k.cache.Get(key)
	if !there {
		value, err = k.DBGet(key)
		switch {
		case err == nil && value != nil:
			k.cache.Set(key, value, 1)
		case err == ErrNotFoundDB:
			err = ErrNotFoundCacheOrDB
		}
	}
	return
}

// Delete removes a KV pair from cache and DB.
func (k *KDB) Delete(key interface{}) error {
	k.cache.Del(key)
	K := k.Encoder(key)
	return k.database.Update(func(txn *badger.Txn) error {
		return txn.Delete(K)
	})
}

// CacheSet inserts a key value pair into cache only.
func (k *KDB) CacheSet(key, value interface{}, cost int64) error {
	if ok := k.cache.Set(key, value, cost); !ok {
		return ErrUnableToCache
	}
	return nil
}

// CacheGet retrieves a value by key from cache without fallback to the DB.
func (k *KDB) CacheGet(key interface{}) (value interface{}, err error) {
	value, there := k.cache.Get(key)
	if !there {
		err = ErrNotFoundCache
	}
	return
}

// DBSet inserts a key value pair into the DB bypassing cache.
func (k *KDB) DBSet(key, value interface{}) error {
	var err ErrMSG
	errd := k.database.Update(func(txn *badger.Txn) error {
		K := k.Encoder(key)
		val := k.Encoder(value)
		e := badger.NewEntry(K, val)
		err := txn.SetEntry(e)
		return err
	})
	if errd != nil {
		err.code = ErrUnableToSave
		err.stack = errd.Error()
	}
	if err.code == ErrNoErr {
		return nil
	}
	return err
}

// DBGet retrieves a value by key from the DB bypassing cache.
func (k *KDB) DBGet(key interface{}) (value interface{}, err error) {
	K := k.Encoder(key)
	err = k.database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(K)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				err = ErrNotFoundDB
			}
			return err
		}
		V, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		value = k.Decoder(V)
		return err
	})
	return
}

// Close closes the KDB.
func (k *KDB) Close() error {
	k.cache.Close()
	return k.database.Close()
}
