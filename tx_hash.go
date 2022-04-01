package flashdb

import (
	"time"
)

func (tx *Tx) HSet(key string, field string, value string) (res int, err error) {
	existVal := tx.HGet(key, field)
	if existVal == value {
		return
	}

	e := newRecordWithValue([]byte(key), []byte(field), []byte(value), HashRecord, HashHSet)
	tx.addRecord(e)
	return
}

func (tx *Tx) HGet(key string, field string) string {
	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return ""
	}

	return toString(tx.db.hashStore.HGet(key, field))
}

func (tx *Tx) HGetAll(key string) []string {
	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return nil
	}

	vals := tx.db.hashStore.HGetAll(key)
	values := make([]string, 0, 1)

	for _, v := range vals {
		values = append(values, toString(v))
	}

	return values
}

func (tx *Tx) HDel(key string, field ...string) (res int, err error) {
	for _, f := range field {
		e := newRecord([]byte(key), []byte(f), HashRecord, HashHDel)
		tx.addRecord(e)
		res++
	}
	return
}

func (tx *Tx) HKeyExists(key string) (ok bool) {
	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return
	}
	return tx.db.hashStore.HKeyExists(key)
}

func (tx *Tx) HExists(key, field string) (ok bool) {
	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return
	}

	return tx.db.hashStore.HExists(key, field)
}

func (tx *Tx) HLen(key string) int {
	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return 0
	}

	return tx.db.hashStore.HLen(key)
}

func (tx *Tx) HKeys(key string) (val []string) {
	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return nil
	}

	return tx.db.hashStore.HKeys(key)
}

func (tx *Tx) HVals(key string) (values []string) {
	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return nil
	}

	vals := tx.db.hashStore.HVals(key)
	for _, v := range vals {
		values = append(values, toString(v))
	}

	return
}

func (tx *Tx) HExpire(key string, duration int64) (err error) {
	if duration <= 0 {
		return ErrInvalidTTL
	}

	if !tx.HKeyExists(key) {
		return ErrInvalidKey
	}

	ttl := time.Now().Unix() + duration
	e := newRecordWithExpire([]byte(key), nil, ttl, HashRecord, HashHExpire)
	tx.addRecord(e)

	return
}

func (tx *Tx) HTTL(key string) (ttl int64) {
	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return
	}

	deadline := tx.db.getTTL(Hash, key)
	if deadline == nil {
		return
	}
	return deadline.(int64) - time.Now().Unix()
}

func (tx *Tx) HClear(key string) (err error) {
	if tx.db.hasExpired(key, Hash) {
		tx.db.evict(key, Hash)
		return
	}

	e := newRecord([]byte(key), nil, HashRecord, HashHClear)
	tx.addRecord(e)
	return
}

func toString(val interface{}) string {
	if val == nil {
		return ""
	}
	return val.(string)
}
