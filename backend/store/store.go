package store

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"gotth/template/backend/db"
	"net/http"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gorilla/sessions"
)

type Store struct {
	provider *db.BadgerProvider
	options  *sessions.Options
}

var store Store

const (
	sessionKeyPrefix = "session:"
	defaultMaxAge    = 86400 * 30 // 30 days
)

func InitStore() error {
	provider := db.GetBadgerProvider()
	if provider == nil || provider.DB == nil {
		return errors.New("Badger provider or DB is nil")
	}

	store = Store{
		provider: provider,
		options: &sessions.Options{
			Path:     "/",
			MaxAge:   defaultMaxAge,
			HttpOnly: true,
			Secure:   false,
		},
	}

	return nil
}

func GetStore() Store {
	return store
}

func (s Store) GetToken(r *http.Request) (*sessions.Session, error) {
	session := sessions.NewSession(&badgerStore{s.provider}, "session")
	session.Options = s.options
	session.IsNew = true

	// Try to get session ID from cookie
	cookie, err := r.Cookie("session")
	if err == nil && cookie.Value != "" {
		sessionID := cookie.Value

		// Load session data from Badger
		err = s.provider.DB.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(sessionKeyPrefix + sessionID))
			if err != nil {
				return err
			}

			return item.Value(func(val []byte) error {
				var sessionData map[interface{}]interface{}
				decoder := gob.NewDecoder(bytes.NewReader(val))
				if err := decoder.Decode(&sessionData); err != nil {
					return err
				}

				session.Values = sessionData
				session.ID = sessionID
				session.IsNew = false
				return nil
			})
		})

		if err != nil && err != badger.ErrKeyNotFound {
			return nil, err
		}
	}

	// Generate new session ID if needed
	if session.IsNew {
		sessionID, err := generateSessionID()
		if err != nil {
			return nil, err
		}
		session.ID = sessionID
		session.Values = make(map[interface{}]interface{})
	}

	return session, nil
}

func (s Store) SaveToken(accessToken string, w http.ResponseWriter, r *http.Request) error {
	if r == nil {
		return fmt.Errorf("request is nil")
	}
	if w == nil {
		return fmt.Errorf("response writer is nil")
	}
	if s.provider == nil {
		return fmt.Errorf("badger provider is nil")
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

	err = s.saveSession(session, w, r)
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

	// Delete from Badger
	err = s.provider.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(sessionKeyPrefix + session.ID))
	})
	if err != nil {
		return err
	}

	// Set cookie to expire immediately
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     s.options.Path,
		MaxAge:   -1,
		HttpOnly: s.options.HttpOnly,
		Secure:   s.options.Secure,
	}
	http.SetCookie(w, cookie)

	return nil
}

func (s Store) AddValue(key string, value interface{}, w http.ResponseWriter, r *http.Request) error {
	session, err := s.GetToken(r)
	if err != nil {
		return err
	}

	session.Values[key] = value

	err = s.saveSession(session, w, r)
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

// Helper methods

func (s Store) saveSession(session *sessions.Session, w http.ResponseWriter, r *http.Request) error {
	// Serialize session data
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(session.Values); err != nil {
		return err
	}

	// Save to Badger with TTL
	err := s.provider.DB.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(sessionKeyPrefix+session.ID), buf.Bytes())
		if s.options.MaxAge > 0 {
			entry = entry.WithTTL(time.Duration(s.options.MaxAge) * time.Second)
		}
		return txn.SetEntry(entry)
	})
	if err != nil {
		return err
	}

	// Set cookie
	cookie := &http.Cookie{
		Name:     "session",
		Value:    session.ID,
		Path:     s.options.Path,
		MaxAge:   s.options.MaxAge,
		HttpOnly: s.options.HttpOnly,
		Secure:   s.options.Secure,
	}
	http.SetCookie(w, cookie)

	return nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// badgerStore implements the gorilla/sessions Store interface
type badgerStore struct {
	provider *db.BadgerProvider
}

func (bs *badgerStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	// This method is required by the Store interface but we handle session retrieval
	// in our custom GetToken method above
	return sessions.NewSession(bs, name), nil
}

func (bs *badgerStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(bs, name)
	session.IsNew = true
	return session, nil
}

func (bs *badgerStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	// This would be called by sessions.Save(), but we handle saving in our custom methods
	return nil
}

/*
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
*/
