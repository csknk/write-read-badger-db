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
	"flag"
	"log"
	"os"
)

var dbPath, inFilePath string
var logfile *os.File

func init() {
	flag.StringVar(&dbPath, "database", "", "Specify a database")
	flag.StringVar(&dbPath, "db", "", "Specify a database (shorthand)")
	flag.Parse()
}

func main() {
	data, err := utilities.NewDatastore(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	data.OutputAll(os.Stdout, true)
	out, err := os.OpenFile("out.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	data.OutputAll(out, true)
	data.OutputAll(out, false)
}
