package hotstuff

// use a local leveldb for now. must convert to use chaindb - bradley
import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

func OpenDB(path string) (*leveldb.DB, error) {
	return leveldb.OpenFile(path, nil)
}

func NewMemDB() *leveldb.DB {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		panic(err.Error())
	}
	return db
}
