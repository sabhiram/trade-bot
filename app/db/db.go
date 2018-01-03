package db

////////////////////////////////////////////////////////////////////////////////

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/sabhiram/trade-bot/types"
)

////////////////////////////////////////////////////////////////////////////////

// db is a JSON serialize-able structure.
type db struct {
	Balances []*types.Balance `json:"Balances"`
	Sessions []*types.Session `json:"Sessions"`
}

////////////////////////////////////////////////////////////////////////////////

// DB is the externally visible db object with locks to guarantee sync
// operations.
type DB struct {
	*sync.RWMutex // Mutex to guard the database

	dbPath string // Path to the db
	db     *db    // Instance to internal db
}

// New returns a db instance from the JSON file specified by `dbPath`.  If the
// file does not exist, it is created and an empty "database" is created.  All
// read and write access to the map and the "database" is guarded by the DB
// mutex.
func New(dbPath string) (*DB, error) {
	d := &DB{
		RWMutex: &sync.RWMutex{},

		dbPath: dbPath,
		db: &db{
			Balances: []*types.Balance{},
			Sessions: []*types.Session{},
		},
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

	return json.Unmarshal(bs, &d.db)
}

// Flush writes the database to the stored `dbPath`.
func (d *DB) Flush() error {
	d.Lock()
	defer d.Unlock()

	bs, err := json.MarshalIndent(d.db, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(d.dbPath, bs, 0755)
}

// Dump prints a debug output of the currently contained database.
func (d *DB) Dump() error {
	d.RLock()
	defer d.RUnlock()

	fmt.Printf("Dumping Balances\n")
	for _, bal := range d.db.Balances {
		fmt.Printf("Found balance: %#v\n", bal)
	}

	fmt.Printf("Dumping Sessions\n")
	for _, ses := range d.db.Sessions {
		fmt.Printf("Found session: %#v\n", ses)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
