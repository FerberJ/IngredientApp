package db

import (
	"fmt"
	"gotth/template/backend/configuration"
	"os"

	c "github.com/ostafen/clover/v2"
)

type CloverProvider struct {
	Database *c.DB
}

var cloverProvider *CloverProvider

func NewCloverProvider(cfg configuration.Configutration) (*CloverProvider, error) {
	dbPath := "clover-db"
	err := os.MkdirAll(dbPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := c.Open("clover-db")
	if err != nil {
		return nil, err
	}

	hasCollection, err := db.HasCollection("recipes")
	if err != nil {
		return nil, err
	}

	if !hasCollection {
		err := db.CreateCollection("recipes")
		if err != nil {
			return nil, err
		}
	}

	setCloverProvider(&CloverProvider{Database: db})
	return cloverProvider, nil
}

func setCloverProvider(provider *CloverProvider) {
	cloverProvider = provider
}

func GetCloverProvider() *CloverProvider {
	return cloverProvider
}

func (cp *CloverProvider) Close() error {
	return cp.Close()
}
