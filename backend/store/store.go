package store

import (
	"github.com/quasoft/memstore"
)

var store *memstore.MemStore

func InitStore() {
	store = memstore.NewMemStore(
		[]byte("authkey123"),
		[]byte("enckey12341234567890123456789012"),
	)
}

func GetStore() *memstore.MemStore {
	return store
}
