/*
Read and write data to Badger DB
Copyright © 2021 David Egan

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

package main

import (
	"backup/utilities"
	"bufio"
	"flag"
	"log"
	"os"
)

var dbPath, inFilePath string
var offset uint64
var logfile *os.File

func init() {
	flag.StringVar(&inFilePath, "infile", "", "Please specify an input file")
	flag.StringVar(&dbPath, "db", "", "Please specify a database")
	flag.Uint64Var(&offset, "offset", 0, "Specify offset (default is 0)")
	flag.Parse()
}

func main() {
	data, err := utilities.NewDatastore(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(inFilePath)
	if err != nil {
		log.Fatal("infile", err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	kvs := []utilities.KeyValue{}
	var i uint64 = offset
	for sc.Scan() {
		k, err := utilities.Uint64ToBytes(i)
		if err != nil {
			log.Fatal(err)
		}
		kvs = append(kvs, utilities.KeyValue{Key: k, Value: []byte(sc.Text())})
		i++
	}
	err = data.WriteBatch(kvs)
	if err != nil {
		log.Fatal(err)
	}
}
