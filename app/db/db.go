package db

////////////////////////////////////////////////////////////////////////////////

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////

type DB struct {
	*sync.RWMutex // Mutex to guard the database

	dbPath string                 // Path to the db
	m      map[string]interface{} // Map of entries
}

// New returns a db instance from the JSON file specified by `dbPath`.  If the
// file does not exist, it is created and an empty "database" is created.  All
// read and write access to the map and the "database" is guarded by the DB
// mutex.
func New(dbPath string) (*DB, error) {
	d := &DB{
		RWMutex: &sync.RWMutex{},
		dbPath:  dbPath,
		m:       map[string]interface{}{},
	}

	return d, d.Load()
}

////////////////////////////////////////////////////////////////////////////////

// Load loads the database from the stored `dbPath`.
func (d *DB) Load() error {
	if _, err := os.Stat(d.dbPath); os.IsNotExist(err) {
		if err := d.Flush(); err != nil {
			return err
		}
	}

	d.Lock()
	defer d.Unlock()

	bs, err := ioutil.ReadFile(d.dbPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, &d.m)
}

// Flush writes the database to the stored `dbPath`.
func (d *DB) Flush() error {
	d.Lock()
	defer d.Unlock()

	bs, err := json.MarshalIndent(d.m, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(d.dbPath, bs, 0755)
}

// Dump prints a debug output of the currently contained database.
func (d *DB) Dump() error {
	d.RLock()
	defer d.RUnlock()

	fmt.Printf("Dumping DB:\n")
	for k, v := range d.m {
		fmt.Printf("Found key: %#v with value: %#v\n", k, v)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
