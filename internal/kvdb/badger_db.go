package kvdb

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync/atomic"

	badger "github.com/dgraph-io/badger/v4"
)

type Badger struct {
	db   *badger.DB
	path string
}

func (s *Badger) WithDataPath(path string) *Badger {
	s.path = path
	return s
}

func (s *Badger) Open() error {
	DataDir := s.GetDbPath()
	if err := os.MkdirAll(path.Dir(DataDir), os.ModePerm); err != nil {
		return err
	}
	option := badger.DefaultOptions(DataDir).
		WithNumVersionsToKeep(1).
		WithLoggingLevel(badger.WARNING)
	db, err := badger.Open(option)
	if err != nil {
		return err
	} else {
		s.db = db
		return nil
	}
}

func (s *Badger) CheckAndGC() {
	lsmSize1, vlogSize1 := s.db.Size()
	for {
		if err := s.db.RunValueLogGC(0.5); err == badger.ErrNoRewrite || err == badger.ErrRejected {
			break
		}
	}
	lsmSize2, vlogSize2 := s.db.Size()
	if vlogSize2 < vlogSize1 {
		fmt.Printf("badger before GC", lsmSize1, lsmSize2)
	} else {
		fmt.Printf("collect zero garbadge")
	}
}

func (s *Badger) Set(k, v []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(k, v)
	})
	return err
}

func (s *Badger) BatchSet(keys, values [][]byte) error {
	if len(keys) != len(values) {
		return errors.New("keys and values length mismatch")
	}
	var err error
	txn := s.db.NewTransaction(true)
	for i, key := range keys {
		value := values[i]
		if err = txn.Set(key, value); err != nil {
			_ = txn.Commit()
			txn = s.db.NewTransaction(true)
			_ = txn.Set(key, value)
		}
	}
	txn.Commit()
	return err
}

func (s *Badger) Get(k []byte) ([]byte, error) {
	var v []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			v = append([]byte{}, val...)
			return nil
		})
		return err
	})
	if err == badger.ErrKeyNotFound {
		return nil, NoDataErr
	}
	return v, err
}

func (s *Badger) BatchGet(keys [][]byte) ([][]byte, error) {
	values := make([][]byte, len(keys))
	err := s.db.View(func(txn *badger.Txn) error {
		for i, key := range keys {
			item, err := txn.Get(key)
			if err == nil {
				err = item.Value(func(val []byte) error {
					values[i] = append([]byte{}, val...)
					return nil
				})
			}
		}
		return nil
	})
	return values, err
}

func (s *Badger) Delete(k []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(k)
	})
}

func (s *Badger) BatchDelete(keys [][]byte) error {
	var err error
	txn := s.db.NewTransaction(true)
	for _, key := range keys {
		if err = txn.Delete(key); err != nil {
			_ = txn.Commit()
			txn = s.db.NewTransaction(true)
			_ = txn.Delete(key)
		}
	}
	return txn.Commit()
}

func (s *Badger) Has(k []byte) bool {
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(k)
		return err
	})
	return err == nil
}

func (s *Badger) IterKey(fn func(k []byte) error) int64 {
	var total int64
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			if err := fn(k); err != nil {
				return err
			}
			atomic.AddInt64(&total, 1)
		}
		return nil
	})
	if err != nil {
		return 0
	}
	return atomic.LoadInt64(&total)
}

func (s *Badger) IterDB(fn func(k, v []byte) error) int64 {
	var total int64
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				if err := fn(k, v); err != nil {
					return err
				}
				atomic.AddInt64(&total, 1)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0
	}
	return atomic.LoadInt64(&total)
}

func (s *Badger) GetDbPath() string {
	return s.path
}

func (s *Badger) WALName() string {
	return s.path
}

func (s *Badger) Close() error {
	return s.db.Close()
}
