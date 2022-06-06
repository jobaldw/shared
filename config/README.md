# shared | config

Reads in **JSON** key value pairs and unmarshals them into any configuration struct. This package relies on a directory that should be at the root of the application called `./configs` and only has 1 public function, **Unmarshal()**.

## How To Use

Here is a view of what the **./configs** directory should look like. Notice that there is no limit on how many config files there are. Filenames can also be whatever you would like.

```text
<working directory>
         |__ /cmd
         |   |__ /svr
         |       |__ main.go
         |
         |__ /configs                <--- Unmarshal() will attempt to read in every file within this folder 
         |   |__ file1.json
         |   |__ file2.json
         |   |__ ...
         |
         |__ ...
```

This is an example of what a config JSON file could look like.

```json
{
    "field1": "a value",
    "name": "api",
    "log_level": "info",
    "port": 3000
}
```

In this example, we created a new struct that has a combination of the package's built-in application struct with the addition of a new field.

``` go
package main

import "github.com/jobaldw/shared/v2/config"

func main() {
    newStruct := struct{
        Field1 string `json:"field1"` // new user defined field 

        config.Application // package built-in struct
    }{}

    err := config.Unmarshal(&newStruct)
    if err != nil {
        // handle error
    }

    // newStruct has been populated
}
```

## Built-In Structs

*Application - configs that are common to an application.*

``` go
type Application struct {
    Name     string `json:"name,omitempty"`
    Port     int    `json:"port,omitempty"`
    LogLevel string `json:"log_level,omitempty"`
}
```

*Clients - configs for one or more api client objects.*

``` go
type Clients struct {
    Clients map[string]Client `json:"clients,omitempty"`
}

type Client struct {
    Headers    map[string]string `json:"headers,omitempty"`
    Health     string            `json:"health,omitempty"`
    Timeout    int               `json:"timeout,omitempty"`
    URL        string            `json:"url,omitempty"`
}
```

*Datasource - configs for one or more mongo database objects.*

``` go
type Datasource struct {
    Mongo map[string]Mongo `json:"mongo,omitempty"`
}

type Mongo struct {
    Database    string            `json:"database,omitempty"`
    URI         string            `json:"uri,omitempty"`
    Collections map[string]string `json:"collections,omitempty"`
}
```
