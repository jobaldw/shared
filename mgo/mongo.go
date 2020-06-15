package mgo

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"gopkg.in/mgo.v2/bson"
)

//Mongo configurations for mongo
type Mongo struct {
	Host string
	Name string
	User string

	URI         *url.URL
	Database    *mongo.Database
	Collections map[string]string
}

//Init mongo instance
func Init(uri *url.URL, database string, collections map[string]string) Mongo {
	return Mongo{
		URI:         uri,
		Name:        database,
		Collections: collections,
	}
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
func (m *Mongo) Insert(doc interface{}, title, key string) (*mongo.InsertOneResult, error) {
	collection := m.Database.Collection(m.Collections[key])
	return collection.InsertOne(context.Background(), doc, options.InsertOne())
}

//FindOne record in mongo
func (m *Mongo) FindOne(id primitive.ObjectID, key string) *mongo.SingleResult {
	filter := bson.M{"_id": id}
	collection := m.Database.Collection(m.Collections[key])
	return collection.FindOne(context.Background(), filter, options.FindOne())
}

//FindMany record in mongo
func (m *Mongo) FindMany(filter []bson.M, key string) (*mongo.Cursor, error) {
	collection := m.Database.Collection(m.Collections[key])
	return collection.Aggregate(context.Background(), filter)
}

//Update record in mongo
func (m *Mongo) Update(id primitive.ObjectID, doc interface{}, key string) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	collection := m.Database.Collection(m.Collections[key])
	return collection.UpdateOne(context.Background(), filter, bson.M{"$set": doc})
}

//Delete record in mongo
func (m *Mongo) Delete(id primitive.ObjectID, key string) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	collection := m.Database.Collection(m.Collections[key])
	return collection.DeleteOne(context.Background(), filter, options.Delete())
}
