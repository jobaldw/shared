# shared | mongo

This package utilizes the [official mongo driver](https://godoc.org/go.mongodb.org/mongo-driver/mongo@v1.7.2 "v1.9.1") making it simpler to configure and connect to a mongo instance for reading and writing.

How To Use

``` go
package main

import (
	"context"

	"github.com/jobaldw/shared/config"
	"github.com/jobaldw/shared/mongo"
)

func main() {

	// Use the shared config.Mongo() struct for mongo client configurations.
	conf := config.Mongo{
		Database: "dbName",
		URI:      "bW9uZ28rc3ZyOi8vPHlvdXJNb25nb1VSST4",
		Collections: map[string]string{
			"key1": "value1",
		},
	}

	// Create new mongo client.
	mongoClient, err := mongo.New(conf)
	if err != nil {
		// handle error
	}

	// The GetCollection() gives you the ability interact with the collection's CRUD operations.
	key1 := mongoClient.GetCollection("key1")
	key1.FindOne(context.Background(), nil)
}
```
