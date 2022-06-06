# shared | router

Utilizing the [gorilla/mux](https://github.com/gorilla/mux "gorilla/mux - v1.8.0") routing and url matcher package, this package quickly creates mux.Routers with health and ready check endpoints and has built-in http writer functions.

## How To Use

Building a new router only takes a few of lines of code.

``` go
package main

import (
    "encoding/json"
    "net/http"

    "github.com/jobaldw/shared/v2/router"
)

func main() {
    // here we are implementing a sever and passing back the created router
    srv, r := router.New(3001, nil)

    // The reason for sending the router back is for customization. Below are 2 server handlers examples:
    //  - 1) We attach the router with only its default endpoints as the server's handler.
    //      e.g.     
    //          srv.Handler = r
    // 
    //  - 2) We wrap the handler in a middleware object and open up the cors policy allowing the router's 
    //      endpoints to have RBAC authentication.
    //      e.g.
    //          srv.Handler = cors.AllowAll().Handler(middleware.Handler(router))
    
    srv.Handler = r // for this example we went with option (1)

    // start and listen on the port configured above
    if err := srv.ListenAndServe(); err != nil {
        // handle err
    }
}
```

### Rename Liveliness & Readiness Endpoints

By default `/health` and `/ready` endpoints are added to the router, but they can be renamed using the `router.New()`.

``` go
package main

import (
    "encoding/json"
    "net/http"

    "github.com/jobaldw/shared/v2/router"
)

func main() {
    // The liveliness and readiness endpoints are updated based on whats populated (not populated) 
    // by the variadic paths parameter. 
    //  - One path input will only update the liveliness path
    //  - Two path inputs will update liveliness and readiness paths
    //  - Zero or more than two will use the default values
    srv, r := router.New(3001, nil, "live","ping")
    srv.Handler = r
    if err := srv.ListenAndServe(); err != nil {
        //...
    }
}
```

### Add More Handlers

We can even add more handlers in addition to the liveliness and readiness handlers.

``` go
package main

import (
    "encoding/json"
    "net/http"

    "github.com/jobaldw/shared/v2/router"
)

func main() {
    srv, r := router.New(3001, nil, "live","ping")
    srv.Handler = r

    // Here we are adding another GET endpoint with the path "/health"
    r.HandleFunc("/hello", helloWorld()).Methods(http.MethodGet)

    if err := srv.ListenAndServe(); err != nil {
        //...
    }
}

func helloWorld() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        router.Respond(w, json.Marshal, 200, router.Message{MSG: "hello world"})
    }
}
```

## Output

When the `/hello` GET endpoint gets invoked, the status would output `200 OK` and the body we look like below.

``` json
{
    "msg": "hello world"
}
```
