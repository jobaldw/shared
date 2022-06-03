# shared | router

Utilizing the [gorilla/mux](https://github.com/gorilla/mux "gorilla/mux - v1.8.0") routing and url matcher package, this package quickly creates mux.Routers with health and ready check endpoints and has built-in http writer functions.

## How To Use

Building a new router only takes a couple of lines of code.

``` go
package main

import (
    "encoding/json"
    "net/http"

    "github.com/jobaldw/shared/router"
)

func main() {
    srv, r := router.New(3001, nil) // here we are implementing a sever and passing back the created router

    // The reason for sending the router back is for customization. Below we have 2 sever configurations.
    //  1) We attach the router as the server's handler. 
    //  2) We wrap the handler in a middle object and open up the cors policy.
    
    // configuration (1)
    srv.Handler = r

    // configuration (2)
    // srv.Handler = cors.AllowAll().Handler(middleware.Handler(router))

    // We can even add more handlers
    r.HandleFunc("/hello", helloWorld()).Methods(http.MethodGet)


    // start and listen on the port configured above
    if err := srv.ListenAndServe(); err != nil {
        // handle err
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
