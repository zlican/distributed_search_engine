package kvdb

import (
	"errors"
	"sync/atomic"

	bolt "go.etcd.io/bbolt"
)

var NoDataErr = errors.New("no data")

type Bolt struct {
	db     *bolt.DB
	path   string
	bucket []byte
}

func (s *Bolt) WithDataPath(path string) *Bolt {
	s.path = path
	return s
}

func (s *Bolt) WithBucket(bucket string) *Bolt {
	s.bucket = []byte(bucket)
	return s
}

func (s *Bolt) Open() error {
	DataDir := s.GetDbPath()
	db, err := bolt.Open(DataDir, 0o600, bolt.DefaultOptions)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(s.bucket)
		return err
	})
	if err != nil {
		db.Close()
		return err
	} else {
		s.db = db
		return nil
	}
}

func (s *Bolt) GetDbPath() string {
	return s.path
}

func (s *Bolt) WALName() string {
	return s.db.Path()
}

func (s *Bolt) Set(k, v []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		return b.Put(k, v)
	})
}

func (s *Bolt) BatchSet(keys, values [][]byte) error {
	if len(keys) != len(values) {
		return errors.New("keys and values length mismatch")
	}
	var err error
	s.db.Batch(func(tx *bolt.Tx) error {
		for i, key := range keys {
			value := values[i]
			tx.Bucket(s.bucket).Put(key, value)
		}
		return nil
	})
	return err
}

func (s *Bolt) Get(k []byte) ([]byte, error) {
	var v []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		v = b.Get(k)
		return nil
	})
	if len(v) == 0 {
		return nil, NoDataErr
	}
	return v, err
}

func (s *Bolt) BatchGet(keys [][]byte) ([][]byte, error) {
	var err error
	values := make([][]byte, len(keys))
	s.db.Batch(func(tx *bolt.Tx) error {
		for i, key := range keys {
			ival := tx.Bucket(s.bucket).Get(key)
			values[i] = ival
		}
		return nil
	})
	return values, err
}

func (s *Bolt) Delete(k []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		return b.Delete(k)
	})
}

func (s *Bolt) BatchDelete(keys [][]byte) error {
	var err error
	s.db.Batch(func(tx *bolt.Tx) error {
		for _, key := range keys {
			tx.Bucket(s.bucket).Delete(key)
		}
		return nil
	})
	return err
}

func (s *Bolt) Has(k []byte) bool {
	var v []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		v = b.Get(k)
		return nil
	})
	if err != nil || string(v) == "" {
		return false
	}

	return true
}

func (s *Bolt) IterKey(fn func(k []byte) error) int64 {
	var total int64
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if err := fn(k); err != nil {
				return err
			}
			atomic.AddInt64(&total, 1)
		}
		return nil
	})
	return atomic.LoadInt64(&total)
}

func (s *Bolt) IterDB(fn func(k, v []byte) error) int64 {
	var total int64
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := fn(k, v); err != nil {
				return err
			}
			atomic.AddInt64(&total, 1)
		}
		return nil
	})
	return atomic.LoadInt64(&total)
}

func (s *Bolt) Close() error {
	return s.db.Close()
}
