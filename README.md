# ACME Parser
This application is designed to handle a stream of contacts and to
save all the channels linked to a contact into a database (SQLite3).

## Installation
Make sure [Golang is installed](https://golang.org/doc/install) on
your system, then follow these steps:
```
$ go get github.com/mattn/go-sqlite3
$ go get github.com/abordji/acme-parser
```

## Usage
The application can take up to three arguments:
* The path of the configuration (see *default.json*)
* The path of the database
* The URI of the stream
```
$ $GOPATH/bin/myapp \
    -conf=$GOPATH/src/myapp/conf.json \
    -db=mydb.sqlite \
    "http://api.local/contacts"
```

## Error handling
* An error during initialization or during the process will stop
  the application.
* A parsing error results in the display of the contact on error
  output.

## To go further
Indexes can be created after the import in order to query efficiently
the database.
