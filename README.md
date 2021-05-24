Read & Write to Badger DB
=========================

Simple utility that:

* Writes lines from a source file into a Badger DB database
* Reads entries from a Badger DB database to stdout/file

Lines in the source file are stored as key-value pairs, with the line number converted to byte slice used as the  key and raw bytes of the line data held as the value.

Records are stored in an append-only manner, with the exception of `height` which tracks DB entries and allows each batch to have unique keys based on an offset.

* TODO: Add Makefile, build binaries for reading & writing
