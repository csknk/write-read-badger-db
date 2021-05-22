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
	"fmt"
	"log"
	"os"

	badger "github.com/dgraph-io/badger"
)

var dbPath, inFilePath string
var logfile *os.File
var write bool

func init() {
	flag.StringVar(&inFilePath, "infile", "", "Please specify an input file")
	flag.StringVar(&dbPath, "db", "", "Please specify a database")
	flag.BoolVar(&write, "write", false, "Specify true if you want to write")
	flag.Parse()
}

func main() {
	fmt.Println("dbPath: ", dbPath)
	fmt.Println("inFilePath: ", inFilePath)
	fmt.Println("write: ", write)

	opts := badger.DefaultOptions(dbPath)
	logfile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("logfile", err)
	}
	defer logfile.Close()
	opts.Logger = utilities.MyLogger(logfile)

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal("Database", err)
	}
	data := utilities.Datastore{db, 0}

	if write {
		f, err := os.Open(inFilePath)
		if err != nil {
			log.Fatal("infile", err)
		}
		defer f.Close()
		sc := bufio.NewScanner(f)
		lines := []string{}
		for sc.Scan() {
			lines = append(lines, sc.Text())
		}
		err = data.WriteBatch(lines)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		data.OutputAll(os.Stdout, true)
		out, err := os.OpenFile("out.txt", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		data.OutputAll(out, true)
	}
}
