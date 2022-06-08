package mongo

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/jobaldw/shared/v2/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// package logging key
const packageKey = "mongo"

var (
	// There needs to be collections to read and write from the database.
	ErrNoCollections = errors.New("no configured collections")

	// No database has been connected.
	ErrNoDatabase = errors.New("no database found")

	// This error will occur is no mongo documents are in the result set based on the CRUD operation.
	ErrNoDocuments = errors.New("no documents in result")
)

// Holds the configurations for a mongo connection using the official mongo driver package.
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
// 	Implements a new mongo object that connects to a database and creates a handle for mongo collections.
// 	* @param conf: mongo configuration
func New(conf config.Mongo) (*Mongo, error) {
	// decode the Base64 encoded mongo uri
	rawURI, err := base64.StdEncoding.DecodeString(conf.URI)
	if err != nil {
		return nil, fmt.Errorf("%s: %s, could not decode url", packageKey, err)
	}

	// builds the mongo uri
	url, err := url.Parse(string(rawURI))
	if err != nil {
		return nil, fmt.Errorf("%s: %s, could not parse url", packageKey, err)
	}

	// build the mongo object with a connected mongo instance and database collection objects
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
// 	Retrieves a mongo collection object.
// 	* @param key: collection key to retrieve mongo collection
func (m *Mongo) GetCollection(key string) *mongo.Collection {
	return m.Collections[key]
}

// Ping
// 	Test the mongo database connection.
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
// 	Test the mongo database connection with any passed in context.Context().
// 	* @param ctx: context used to handle any cancellations
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

/********** Helper functions **********/

// connect
// 	Creates, configures and connects to a mongo client.
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
