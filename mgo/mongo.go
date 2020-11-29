package mgo

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"gopkg.in/mgo.v2/bson"
)

// ErrNoDocuments error
var ErrNoDocuments = mongo.ErrNoDocuments

//Mongo configurations for mongo
type Mongo struct {
	Host string
	Name string
	User string

	URI      *url.URL
	Database *mongo.Database
}

//Init mongo instance
func Init(rel *url.URL, database string) Mongo {
	return Mongo{
		URI:  rel,
		Name: database,
	}
}

// Parse url
func Parse(uri string) (rel *url.URL, err error) {
	rawURI, err := base64.RawStdEncoding.DecodeString(uri)
	if err != nil {
		return
	}

	rel, err = url.Parse(string(rawURI))
	if err != nil {
		return
	}
	return
}

//Connect to mongo
func (m *Mongo) Connect() (err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(m.URI.String()))
	if err != nil {
		return fmt.Errorf("could not create mongo object, %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return fmt.Errorf("could not connect to database, %s", err)
	}

	m.Database = client.Database(m.Name)
	m.Host, m.User = m.URI.Hostname(), m.URI.User.Username()

	return err
}

//Ping mongo
func (m *Mongo) Ping() (err error) {
	if m.Database == nil {
		return fmt.Errorf("could not connect to database")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = m.Database.Client().Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("could not ping to database, %s", err)
	}

	return err
}

//Insert record in mongo
func Insert(ctx context.Context, collection *mongo.Collection, doc interface{}) (*mongo.InsertOneResult, error) {
	return collection.InsertOne(ctx, doc, options.InsertOne())
}

//FindOne record in mongo
func FindOne(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID) (*mongo.SingleResult, error) {
	filter := bson.M{"_id": id}

	result := collection.FindOne(ctx, filter, options.FindOne())
	return result, result.Err()
}

//FindMany records in mongo
func FindMany(ctx context.Context, collection *mongo.Collection, filter []bson.M) (*mongo.Cursor, error) {
	return collection.Aggregate(ctx, filter, options.Aggregate())
}

//Update record in mongo
func Update(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID, doc interface{}) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}

	return collection.UpdateOne(context.Background(), filter, bson.M{"$set": doc}, options.Update().SetUpsert(true))
}

//Delete record in mongo
func Delete(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}

	return collection.DeleteOne(context.Background(), filter, options.Delete())
}
