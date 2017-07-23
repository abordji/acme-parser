package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"flag"
)

// ChannelDef represents the definition of a channel
type ChannelDef struct {
	Name    string `json:"name"`
	Prefix  string `json:"prefix"`
	Pattern string `json:"pattern"`
	regExp  *regexp.Regexp
}

// Channel represents a single channel of a contact
type Channel struct {
	id    int
	name  string
	value string
}

// loadConf loads the configuration
func loadConf(filename string) ([]*ChannelDef, error) {

	var cds []*ChannelDef

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(file, &cds); err != nil {
		return nil, err
	}
	for _, cd := range cds {
		cd.regExp, err = regexp.Compile(cd.Pattern)
		if err != nil {
			return nil, err
		}
	}

	return cds, nil
}

// openDb initializes the database
func openDb(dbname string) (*sql.DB, error) {

	os.Remove(dbname)

	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec(`
	CREATE TABLE contacts (
		id    INTEGER NOT NULL,
		name  TEXT    NOT NULL,
		value TEXT    NOT NULL
	);
	DELETE FROM contacts;
	`); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// openStream opens a stream from an given URI
func openStream(uri string) (io.ReadCloser, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// findChannelDef finds the next type of channel present in the buffer
func findChannelDef(buf string, cds []*ChannelDef) *ChannelDef {
	for _, cd := range cds {
		if strings.HasPrefix(buf, cd.Prefix) {
			return cd
		}
	}

	// Unknown prefix
	return nil
}

func main() {

	confPath := flag.String("conf", "", "The path to the configuration")
	dbPath := flag.String("db", "", "The path to the database")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Println(fmt.Sprintf("Usage of %s: The URI of the stream is required", os.Args[0]))
		return
	}
	if *confPath == "" {
		*confPath = fmt.Sprintf("%s.json", os.Args[0])
	}
	if *dbPath == "" {
		*dbPath = fmt.Sprintf("%s.sqlite", os.Args[0])
	}

	cf, err := loadConf(*confPath)
	if err != nil {
		log.Fatalf("Cannot load configuration: %s", err)
	}

	db, err := openDb(*dbPath)
	if err != nil {
		log.Fatalf("Cannot open database: %s", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Cannot start transaction: %s", err)
	}

	stmt, err := tx.Prepare("INSERT INTO contacts(id, name, value) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatalf("Cannot prepare statement: %s", err)
	}
	defer stmt.Close()

	stream, err := openStream(flag.Arg(0))
	if err != nil {
		log.Fatalf("Cannot open stream: %s", err)
	}
	defer stream.Close()

	// Create a buffer for all the channels of a contact
	cs := make([]*Channel, 0, 8)

	scanner := bufio.NewScanner(stream)
	for id := 1; scanner.Scan(); id++ {

		parseErr := false
		contact := strings.TrimSpace(scanner.Text())

		// Work on a copy in case of parsing error
		buf := contact

		for buf != "" {

			// Consume prefix
			cd := findChannelDef(buf, cf)
			if cd == nil {
				parseErr = true
				break
			}
			buf = strings.TrimPrefix(buf, cd.Prefix)

			// Consume value
			val := cd.regExp.FindString(buf)
			if val == "" {
				parseErr = true
				break
			}
			buf = strings.TrimPrefix(buf, val)

			// All the channels of a contact share the same id
			cs = append(cs, &Channel{id, cd.Name, val})
		}

		if parseErr {
			fmt.Fprintln(os.Stderr, contact)
			continue
		}

		for _, c := range cs {
			_, err = stmt.Exec(c.id, c.name, c.value)
			if err != nil {
				log.Fatalf("Cannot execute statement: %s", err)
			}
		}

		// Empty the buffer but preserve its capacity
		cs = cs[:0]
	}

	tx.Commit()
}
