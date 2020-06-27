# shared [![Go Report Card](https://goreportcard.com/badge/github.com/jobaldw/shared)](https://goreportcard.com/report/github.com/jobaldw/shared) [![Last Commit](https://img.shields.io/github/last-commit/jobaldw/shared)](https://img.shields.io/github/last-commit/jobaldw/shared) [![Release Version](https://img.shields.io/github/v/release/jobaldw/shared)](https://img.shields.io/github/v/release/jobaldw/shared)

> A library of code that can be used in Go applications.

![JB Desgins](https://github.com/jobaldw/shared/blob/master/jb-icon.jpg)

## config

Reads in **json** key value pairs that are unmarshaled into one configuration struct that can be passed around through an application.

### Example config/application.json

```json
{
    "application": "api",
    "log_level": "info",
    "port": 3000
}
```

*Config - configuration that holds all configurables.*

``` go
type Config struct {
    Application Application
    Datasource  Datasource
    Dependents  Dependents
}
```

This package relies on a directory that should be at the root of the application called `config`. The **Unmarshal()** function looks for two json files named `application.json` and `datasource.json`.

- config
  - application.json // *gets read into the Application struct*
  - datasource.json  // *gets read into the Datasource struct*
  - dependents.json  // *gets read into the Dependents struct*

*App - configurables for common application related objects.*

``` go
type Application struct {
    Name     string `json:"application,omitempty"`
    Port     int    `json:"port,omitempty"`
    LogLevel string `json:"log_level,omitempty"`
    Auth0    Auth0  `json:"auth0,omitempty"`
```

*Auth0 - configurables for Auth0 middleware authentication.*

``` go
type Auth0 struct {
    Identifier string `json:"identifier,omitempty"`
    Domain     string `json:"domain,omitempty"`
}
```
*Datasource - configurations for one or more mongo database objects.*

``` go
type Datasource struct {
    Mongo  Mongo            `json:"database,omitempty"`
    Mongos map[string]Mongo `json:"databases,omitempty"`


type Mongo struct {
    Database    string            `json:"database,omitempty"`
    URI         string            `json:"uri,omitempty"`
    Collections map[string]string `json:"collections,omitempty"`
}
```

*Dependents - configurations for one or more api client objects.*

``` go
type Dependents struct {
    Client  Client            `json:"client,omitempty"`
    Clients map[string]Client `json:"clients,omitempty"`
}

type Client struct {
    Health     string            `json:"health,omitempty"`
    URL        string            `json:"url,omitempty"`
    Timeout    int               `json:"timeout,omitempty"`
    Headers    map[string]string `json:"headers,omitempty"`
}
```

## log

Utilizes the [sirupsen/logrus](https://github.com/sirupsen/logrus "sirupsen/logrus - v1.6.0") logging pacakge. Wraps a few functions to help configure and help make log messages more informational.

### Example json message

```json
{
    "application":"api",
    "file":"file.go:10",
    "function":"func()",
    "id": 123-456-7890,
    "level":"info",
    "msg":"a message",
    "time":"2020-06-09T12:30:00-05:00"
}
```

**Note**: all log messages will have *application*, *level*, *msg* and *time* fields by default. See **Usage**  for how to log with additional fields.

### log Set Up

``` go
package main

import "github.com/jobaldw/shared/log"

func main() {
    ...
    log.Configure("api", "info")
    ...
}
```

### log Usage

``` go
func foo() {
    log.Info("a message") // prints default fields
    log.Details().Info("a message") // prints default fields, file, and function
    log.Add(log.Fields{"id": 123-456-7890}).Info("a message") // prints Details() and any additional fields specified
}
```

## client

Makes http client calls simpler while still allowing for customization for any CRUD application.

*Client - client specific variables.*

``` go
type Client struct {
    Health     string
    Headers    map[string]string

    URL    *url.URL
    Client *http.Client
}
```

*Response - dependent client responses.*

``` go
type Response struct {
    Status string
    Code   int

    Body Body
}

type Body struct {
    String string
    Bytes  []byte
    IO     io.Reader
}
```

### client Set Up

For mulitple clients, range over every client object keeping the same key and retreiving that key's values neccessary to build a new client. For one client, do the smae without ranging.

``` go
package main

import "github.com/jobaldw/shared/client"

func main() {
    ...
    clients := make(map[string]client.Client)
    for key, value := range conf.Dependents {
        newClient, err := client.New(value.URL, value.Health, value.Timeout, value.Headers)
        if err != nil {
            return Controller{}, err
        }

        clients[key] = newClient
    }
    ...
}
```

### client Usage

#### http request

These http request all return a response.Response or an error.

``` go
func main() {
    ...

    clients["key"].Get("path", params) /* or */ client.Get("path", params)
    clients["key"].Put("path", params, body) /* or */ client.Put("path", params, body)
    clients["key"].Post("path", params, body) /* or */ client.Post("path", params, body)
    clients["key"].Delete("path", params) /* or */ client.Delete("path", params)
}
```

## router

Utilizes the [gorilla/mux](https://github.com/gorilla/mux "gorilla/mux - v1.7.4") routing and url matcher pacakge. Primarly initializes a mux.Router with health and ready check endpoints and writes to the client.

**Note**: This package is dependent on the **client** package.

*Resp - api client responses.*

``` go
type Resp struct {
    ID      interface{} `json:"id,omitempty"`
    Payload interface{} `json:"payload,omitempty"`

    Status string `json:"status,omitempty"`
    MSG    string `json:"msg,omitempty"`
    ERR    string `json:"error,omitempty"`
}
```

### router Set Up

In order to use the `router.New()` function, we will use the `clients` created in the **client** section.

``` go
package main

import "github.com/jobaldw/shared/router"

func main() {
    ...
    newRouter := router.New("api", clients) // see config to view "clients" declaration
}
```

### router Usage

You can add as many endpoints as you want. Each instantiation needs a function to call and each function can have its own personalized response.

``` go
func main() {
    ...
    newRouter.HandleFunc("/endpoint", foo()).Methods(http.MethodGet)

    ...

    // starts the server and keeps it open
    http.ListenAndServe(
        8000,
        r,
    )
}

func foo() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        payload := struct{
            ID string
        }{
            ID: "abc123",
        }

        resp := router.Resp{Status: "Up", MSG: "hello world", Payload: payload}

        router.Response(w, 200, resp)
        return
    }
}
```

#### Output

The status of this response is `200 OK`

``` json
{
    "payload": {
        "ID": "abc123"
    },
    "status": "Up",
    "msg": "hello world"
}
```

## middlware

Utilizes the [auth0](https://auth0.com/ "auth0") middlware capability to validate user access to API endpoints using Bearer token authentication.

### middleware Set Up

``` go
package main

import "github.com/jobaldw/shared/config"
import "github.com/jobaldw/shared/middleware"

func main() {
    ... 
    
    // pass in configurables
    middleware.New(conf.Application.Auth0)

    ...

    // wrap function with middleware.Auth0() wrapper
	newRouter.HandleFunc("/endpoint", middleware.Auth0(foo())).Methods(http.MethodGet)

    ...

    http.ListenAndServe(
        8000,
        middleware.Handler(r), // wrap the router with the middle handler options
    )
}
```

## mgo

Utilizes the [offical mongo driver](https://godoc.org/go.mongodb.org/mongo-driver/mongo "mongo-driver") to make mongo request simpler.

*Mongo - client configurations.*

``` go
type Mongo struct {
    Host string
    Name string
    User string

    URI         *url.URL
    Database    *mongo.Database
    Collections map[string]string
}
```

### mgo Set Up

``` go
import  "github.com/jobaldw/shared/mgo"

func main() {
    uri, err := mgo.Parse("bW9uZ28rc3ZyOi8vZmFrZVVSTFBhdGhUb01vbmdv")
    if err != nil {
        return
    }

    ds = mgo.Init(uri, "movieDB", map[string]string{"Act":"action", "Adv":"adventure", "Rom":"romance"})
    if err := ds.Connect(); err != nil {
        return ds, err
    }

    if err := ds.Ping(); err != nil {
        return ds, err
    }
}
```

After successful connection you will have access to do the following mongo request:
* Insert()
* FindOne()
* FindMany()
* Update()
* Delete()
