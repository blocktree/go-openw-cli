package openwcli

import (
	"fmt"
	"github.com/asdine/storm"
)

type StormDB struct {
	*storm.DB
	FileName string
	Opened bool
}

//OpenStormDB
func OpenStormDB(filename string, stormOptions ...func(*storm.Options) error) (*StormDB, error) {

	db, err := storm.Open(filename, stormOptions...)
	//fmt.Println("open app db")
	if err != nil {
		return nil, fmt.Errorf("can not open dbfile: '%s', unexpected error: %v", filename, err)
	}

	// Check the metadata.
	stormDB := &StormDB{
		FileName: filename,
		DB:       db,
		Opened: true,
	}

	return stormDB, nil
}

// Close closes the database.
func (db *StormDB) Close() error {
	err := db.DB.Close()
	if err != nil {
		return err
	}
	db.Opened = false
	return nil
}

