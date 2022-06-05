# shared | client

A simple http client that make client calls quicker and lighter while still allowing for customization. Intended on being used with the `shared/config` package. The package only allows for the following request methods: GET, POST, PUT, DELETE with or without contexts.

## How To Use

Here is an example of what a config JSON file could look like.

```json
{
    "clients": {
        "myApp1": {
            "headers": {
                "content-type": [
                    "application/json"
                ]
            },
            "health": "/ping",
            "timeout": 10,
            "url": "https://www.test.com"
        },
        "myApp2": {
            "headers": {
                "content-type": [
                    "application/json"
                ]
            },
            "health": "/v2/health",
            "timeout": 10,
            "url": "https://www.fake.com"
        }
    }
}
```

Below is an example on how to we use the configs from above to create two clients.

``` go
package main

import (
    "fmt"

    "github.com/jobaldw/shared/client"
    "github.com/jobaldw/shared/config"
)

func main() {
    newStruct := config.Clients{}

    if err := config.Unmarshal(&newStruct); err != nil {
        // handle error
    }

    // because the built-in clients configs have a map of string to clients, we can loop through
    // each client within the config. You may also use the built-in client config allowing for only
    // one client.
    clients := make(map[string]*client.Client)
    for k, v := range newStruct.Clients {
        newClient, err := client.New(v)
        if err != nil {
            // handle error
        }
        clients[k] = newClient
        
        // you now have the ability to make http requests with the given client.
        fmt.Println(newClient)
    }
}

```

OUTPUT

``` text
$ go run main.go 
&{/ping 0xc00007f1d0 map[content-type:[application/json]] https://www.test.com}
&{/v2/health 0xc00007f290 map[content-type:[application/json]] https://www.fake.com}
```
