package db

////////////////////////////////////////////////////////////////////////////////

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/sabhiram/walley/types"
	"github.com/sabhiram/walley/uuid"
)

////////////////////////////////////////////////////////////////////////////////

var (
	ErrInvalidCurrency = errors.New("invalid currency specified")
)

////////////////////////////////////////////////////////////////////////////////

// Currencies keeps track of the list of track-able currencies.  This is a
// non-mutable map of items and is NOT intended to be modified during runtime.
var Currencies = map[string]*types.Currency{
	"bitcoin":  types.NewCurrency("AA0B55A6-3F0B-422F-A6D6-7EFF68FC5E45", "Bitcoin", "BTC"),
	"ethereum": types.NewCurrency("8C6E9D45-57A0-4B0E-BF5C-9DD1079625F4", "Ethereum", "ETH"),
	"pivx":     types.NewCurrency("C0B6B568-7EDF-4510-A30F-5E3B8D794C0E", "PIVX", "PIVX"),
}

////////////////////////////////////////////////////////////////////////////////

type DB struct {
	*sync.RWMutex // Mutex to guard the database

	dbPath string                               // Path to the db
	m      map[uuid.UUID][]*types.CurrencyEntry // Map of entries
}

// New returns a db instance from the JSON file specified by `dbPath`.  If the
// file does not exist, it is created and an empty "database" is created.  All
// read and write access to the map and the "database" is guarded by the DB
// mutex.
func New(dbPath string) (*DB, error) {
	d := &DB{
		RWMutex: &sync.RWMutex{},
		dbPath:  dbPath,
		m:       map[uuid.UUID][]*types.CurrencyEntry{},
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

func (d *DB) Add(currencyName, publicKey string) error {
	cn := strings.ToLower(currencyName)
	if c, ok := Currencies[cn]; ok {
		id := uuid.New()
		ne, err := types.NewEntry(c, publicKey)
		if err != nil {
			return err
		}

		d.Lock()
		d.m[id] = []*types.CurrencyEntry{ne}
		d.Unlock()
		return nil
	}

	return ErrInvalidCurrency
}

////////////////////////////////////////////////////////////////////////////////
