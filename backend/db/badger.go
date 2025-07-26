package db

import (
	"gotth/template/backend/configuration"

	"github.com/dgraph-io/badger/v4"
)

type BadgerProvider struct {
	DB *badger.DB
}

var badgerProvider *BadgerProvider

func NewBadgerProvider(cfg configuration.Configutration) (*BadgerProvider, error) {
	dbPath := cfg.BadgerPath
	if dbPath == "" {
		dbPath = "./badger"
	}

	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	setBadgerProvider(&BadgerProvider{DB: db})
	return badgerProvider, nil
}

func setBadgerProvider(provider *BadgerProvider) {
	badgerProvider = provider
}

func GetBadgerProvider() *BadgerProvider {
	return badgerProvider
}

func (bp *BadgerProvider) Close() error {
	return bp.DB.Close()
}
