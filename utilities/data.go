/*
Qredochain Backup
Copyright 2021 David Egan

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utilities

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	badger "github.com/dgraph-io/badger"
)

type DBName string

type Datastore struct {
	DB    *badger.DB
	Count uint64
}

type KeyValue struct {
	Key   []byte
	Value []byte
}

func NewDatastore(dbPath string) (*Datastore, error) {
	opts := badger.DefaultOptions(dbPath)
	badgerLogfile, err := os.OpenFile("badger-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer badgerLogfile.Close()
	opts.Logger = MyLogger(badgerLogfile)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &Datastore{db, 0}, nil
}

func (d *Datastore) IsSet() bool {
	if d.DB == nil {
		return false
	}
	return true
}

// Make a call to GetKeyValue and interpret the response to determine if the DB has the given key.
func (d *Datastore) Has(key []byte) (ok bool) {
	_, err := d.GetKeyValue(key)
	switch err {
	case badger.ErrKeyNotFound:
		ok = false
	case nil:
		ok = true
	}
	return
}

func (d *Datastore) GetKeyValue(key []byte) ([]byte, error) {
	var valCopy []byte
	err := d.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		return nil
	})
	return valCopy, err
}

func (d *Datastore) KeyIsSet(key string) bool {
	result := false
	err := d.DB.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err == nil {
			result = true
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func (d *Datastore) SetKeyValue(key, value []byte) error {
	//	d.Logger.Infof("Setting key %s; value %s", string(key), string(value))
	err := d.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
	if err != nil {
		//		d.Logger.Errorf("Problem setting key %s; value %s", string(key), string(value))
		return err
	}

	return nil
}

func (d *Datastore) WriteBatch(data []KeyValue) error {
	wb := d.DB.NewWriteBatch()
	defer wb.Cancel()
	var recordsWritten uint64 = 0
	height, err := d.Height()
	if err != nil {
		return err
	}

	for _, kv := range data {
		if d.KeyIsSet(string(kv.Key)) {
			continue
		}
		err := wb.Set(kv.Key, kv.Value)
		if err != nil {
			return err
		}
		recordsWritten++
	}
	err = wb.Flush()
	if err != nil {
		return err
	}
	newHeight := height + recordsWritten
	heightBytes, err := Uint64ToBytes(newHeight)
	if err != nil {
		return err
	}
	d.SetKeyValue([]byte("height"), heightBytes)
	return nil
}

func (d *Datastore) Height() (uint64, error) {
	var height uint64
	hasHeight := d.Has([]byte("height"))
	if hasHeight {
		oldHeight, err := d.GetKeyValue([]byte("height"))
		if err != nil {
			return 0, err
		}
		height = BytesToUint64(oldHeight)
	}

	return height, nil
}

// Send all key value pairs to stdout
func (d *Datastore) OutputAll(f io.Writer, valuesOnly bool) error {
	err := d.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var val []byte
			key := item.Key()
			err := item.Value(func(v []byte) error {
				val = v
				return nil
			})
			if err != nil {
				return err
			}
			var printVal string
			if string(key) == "height" {
				printVal = strconv.FormatUint(BytesToUint64(val), 10)
			} else {
				printVal = string(val)
			}
			if valuesOnly {
				fmt.Fprintf(f, "%s\n", printVal)
				continue
			}
			fmt.Fprintf(f, "%s\t%s\n", decodeKey(key), printVal)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func Uint64ToBytes(num uint64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func BytesToUint64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

//func (d *Datastore) Errorf(format string, v ...interface{}) {
//	if d.Logger == nil {
//		return
//	}
//	d.Logger.Errorf(format, v...)
//}
//
//// Infof logs an INFO message to the logger specified in d..
//func (d *Datastore) Infof(format string, v ...interface{}) {
//	if d.Logger == nil {
//		return
//	}
//	d.Logger.Infof(format, v...)
//}
//
//// Warningf logs a WARNING message to the logger specified in d..
//func (d *Datastore) Warningf(format string, v ...interface{}) {
//	if d.Logger == nil {
//		return
//	}
//	d.Logger.Warningf(format, v...)
//}
//
//// Debugf logs a DEBUG message to the logger specified in d..
//func (d *Datastore) Debugf(format string, v ...interface{}) {
//	if d.Logger == nil {
//		return
//	}
//	d.Logger.Debugf(format, v...)
//}
