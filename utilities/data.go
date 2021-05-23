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
	logfile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer logfile.Close()
	opts.Logger = MyLogger(logfile)

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

func (d *Datastore) GetKeyValue(key string) ([]byte, error) {
	var valCopy []byte
	err := d.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			log.Fatal(err)
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

func (d *Datastore) SetKeyValue(key []byte, value []byte) error {
	// Not implemented
	return nil
}

//func (d *Datastore) WriteBatch(data []string) error {
func (d *Datastore) WriteBatch(data []KeyValue) error {
	wb := d.DB.NewWriteBatch()
	defer wb.Cancel()

	for _, kv := range data {
		// Check for existing key && value at this key
		if d.KeyIsSet(string(kv.Key)) {
			continue
		}
		err := wb.Set(kv.Key, kv.Value) // Will create txns as needed.
		if err != nil {
			return err
		}
	}
	err := wb.Flush() // Wait for all txns to finish.
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// Send all key value pairs to stdout
func (d *Datastore) OutputAll(f io.Writer, valuesOnly bool) error {
	err := d.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		//		format := "%s\t%s\n"
		//		tw := new(tabwriter.Writer).Init(f, 12, 8, 2, ' ', 0)
		//		fmt.Fprintf(tw, format, "Key", "Value")
		//		fmt.Fprintf(tw, format, "---", "-----")
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
			if valuesOnly {
				fmt.Fprintf(f, "%s\n", string(val))
			} else {
				fmt.Fprintf(f, "%s\t%s\n", decodeKey(key), string(val))

			}
			//			fmt.Fprintf(tw, format, decodeKey(key), string(val))

			//			fmt.Printf(
			//				"%s%s: %s%s%s\n",
			//				string(colorYellow),
			//				decodeKey(key),
			//				string(colorWhite),
			//				string(val),
			//				string(colorReset),
			//			)
		}
		//		tw.Flush()
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
