Read & Write to Badger DB
=========================

Simple utility that:

* Writes lines from a source file into a Badger DB database
* Reads entries from a Badger DB database to stdout/file

Lines in the source file are stored as key-value pairs, with the line number converted to byte slice used as the  key and raw bytes of the line data held as the value.

TODO: Makefile
TODO: Prevents storage of duplicated
TODO: Access individual records, lookup by key
TODO: Benefits of storing the total number of records?
