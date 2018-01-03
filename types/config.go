package types

////////////////////////////////////////////////////////////////////////////////

import "time"

////////////////////////////////////////////////////////////////////////////////

// Config encapsulates app wide configuration settings.
type Config struct {
	RefreshInterval time.Duration // conditions check refresh interval
	ApiKey          string        // bittrex api key
	Secret          string        // bittrex secret
	DbPath          string        // path to local session db
	Args            []string      // other command line args
}

////////////////////////////////////////////////////////////////////////////////
