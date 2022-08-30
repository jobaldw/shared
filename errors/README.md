# shared | errors

Error package with built in structs designed for dependent client controller errors and mongo operation errors.

## Built-In Error Structs

*Client error structs.*

``` go
type ClientErr struct {
	Client     string `json:"client,omitempty"`
	URI        string `json:"uri,omitempty"`
	Status     string `json:"status,omitempty"`
	StatusCode int    `json:"code,omitempty"`
	Msg        string `json:"message,omitempty"`
	Err        error  `json:"error,omitempty"`
}
```

*Mongo error structs.*

``` go
type MongoErr struct {
	Operation  string `json:"operation,omitempty"`
	Collection string `json:"collection,omitempty"`
	Err        error  `json:"error,omitempty"`
}
```
