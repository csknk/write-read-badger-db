/*
Save stdin to Badger DB
Copyright © 2021 DavidEgan

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package utilities

import (
	"io"
	"log"
)

type myLog struct {
	*log.Logger
}

// The package supplied logger:
// `var defaultLogger = &defaultLog{Logger: log.New(os.Stderr, "x", log.LstdFlags)}`
// ...where myLog has methods that meet the contract specified by the interface.
// Amend to a function returning a pointer to a myLog struct that contains a Logger
// initialised with the supplied io.Writer
func MyLogger(output io.Writer) *myLog {
	return &myLog{Logger: log.New(output, "csknk DBTool ", log.LstdFlags)}
}

func (l *myLog) Errorf(f string, v ...interface{}) {
	l.Printf("ERROR: "+f, v...)
}

func (l *myLog) Warningf(f string, v ...interface{}) {
	l.Printf("WARNING: "+f, v...)
}

func (l *myLog) Infof(f string, v ...interface{}) {
	l.Printf("INFO: "+f, v...)
}

func (l *myLog) Debugf(f string, v ...interface{}) {
	l.Printf("DEBUG: "+f, v...)
}
