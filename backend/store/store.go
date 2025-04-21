package store

import (
	"context"
	"fmt"
	"gotth/template/backend/db"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v9"
)

type Store struct {
	Store *redisstore.RedisStore
}

var store Store

func InitStore() error {
	provider := db.GetRedisProvider()
	if provider == nil || provider.Client == nil {
		return fmt.Errorf("Redis provider or client is nil")
	}

	var err error
	// New default RedisStore
	s, err := redisstore.NewRedisStore(context.Background(), provider.Client)
	if err != nil {
		return err
	}

	store = Store{s}

	return nil
}

func GetStore() Store {
	return store
}

func (s Store) GetToken(r *http.Request) (*sessions.Session, error) {
	return s.Store.Get(r, "session")
}
func (s Store) SaveToken(accessToken string, w http.ResponseWriter, r *http.Request) error {
	if r == nil {
		return fmt.Errorf("request is nil")
	}
	if w == nil {
		return fmt.Errorf("response writer is nil")
	}
	if s.Store == nil {
		return fmt.Errorf("redis store is nil")
	}

	session, err := s.GetToken(r)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session is nil after GetToken")
	}

	if session.Values == nil {
		session.Values = make(map[interface{}]interface{})
	}

	session.Values["token"] = accessToken

	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) DeleteToken(w http.ResponseWriter, r *http.Request) error {
	session, err := s.GetToken(r)

	if err != nil {
		return err
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

func (s Store) AddValue(key string, value interface{}, w http.ResponseWriter, r *http.Request) error {
	session, err := s.GetToken(r)
	if err != nil {
		return err
	}

	session.Values[key] = value

	err = sessions.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) GetValue(key string, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	session, err := s.GetToken(r)
	if err != nil {
		return nil, err
	}

	value := session.Values[key]

	return value, nil
}
