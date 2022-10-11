/*
Package mongo implements a Mongo instance that connects to a
database and allows interactions to the datasource.

All functions that require a context to be passed should be given one from
the service handler request to correctly handle cancellations.
*/
package mongo

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jobaldw/shared/v2/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// package logging key
const packageKey = "mongo"

var (
	ErrNoCollections = errors.New("no configured collections") // no collection to read from or write to
	ErrNoDatabase    = errors.New("no database found")         // no database has been connected
	ErrNoDocuments   = errors.New("no documents in result")    // no mongo documents are in the result set
)

// Holds the configurations for a mongo connection using the
// official mongo driver package.
type Mongo struct {
	// host or host:port of the mongo uri
	Host string

	// mongo database name
	Name string

	// username/password information for uri authentication
	User string

	// handles to a MongoDB collection
	Collections map[string]*mongo.Collection

	// handle to a MongoDB database
	Database *mongo.Database

	// representation of a parsed URL (technically, a URI reference)
	URI *url.URL
}

// New
// implements a new mongo object that connects to a database and creates a
// handle for mongo collections.
func New(conf config.Mongo) (*Mongo, error) {
	// builds the mongo uri
	url, err := getURI(conf)
	if err != nil {
		return nil, err
	}

	// build the mongo object with a connected mongo instance and database
	// collection objects
	db := &Mongo{URI: url, Name: conf.Database, Collections: make(map[string]*mongo.Collection)}
	if err := db.connect(); err != nil {
		return db, fmt.Errorf("%s: %s", packageKey, err)
	}

	if conf.Collections != nil {
		for k, v := range conf.Collections {
			db.Collections[k] = db.Database.Collection(v)
		}
		return db, nil
	}

	return db, fmt.Errorf("%s: %s", packageKey, ErrNoCollections)
}

// GetCollection
// retrieves a mongo collection object.
func (m *Mongo) GetCollection(key string) *mongo.Collection {
	return m.Collections[key]
}

// Ping
// test the mongo database connection.
func (m *Mongo) Ping() (err error) {
	if m.Database == nil {
		return fmt.Errorf("%s: %s", packageKey, ErrNoDatabase)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = m.Database.Client().Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("%s: %s, could not ping database, %s", packageKey, err, m.Database.Name())
	}
	return err
}

// PingWithContext
// test the mongo database connection with any passed in context.Context().
func (m *Mongo) PingWithContext(ctx context.Context) (err error) {
	if m.Database == nil {
		return fmt.Errorf("%s: %s", packageKey, ErrNoDatabase)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err = m.Database.Client().Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("%s: %s, could not ping database, %s", packageKey, err, m.Database.Name())
	}
	return err
}

/********** helper functions **********/

// connect
// creates, configures and connects to a mongo client.
func (m *Mongo) connect() (err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(m.URI.String()))
	if err != nil {
		return fmt.Errorf("%s: %s, could not create mongo object", packageKey, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return fmt.Errorf("%s: %s, could not connect to %s", packageKey, err, m.Database.Name())
	}

	m.Database = client.Database(m.Name)
	m.Host, m.User = m.URI.Hostname(), m.URI.User.Username()
	return err
}

// getURI
// builds the decodes and builds the mongo db url.
func getURI(conf config.Mongo) (*url.URL, error) {
	// decode the Base64 encoded mongo uri
	rawURI, err := base64.StdEncoding.DecodeString(conf.URI)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", packageKey, err)
	}

	// use configured username and password over the uri
	uri := string(rawURI)
	if conf.Password != "" && conf.Username != "" {
		user := url.QueryEscape(conf.Username)
		pass := url.QueryEscape(conf.Password)
		split := strings.Split(uri, "@")
		uri = "mongodb+srv://" + user + ":" + pass + "@" + split[1]
	}

	// build url
	url, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", packageKey, err)
	}
	return url, nil
}
