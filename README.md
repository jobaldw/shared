# shared

> A library of code that can be used in Go applications.

![JB Desgins](https://github.com/jobaldw/shared/blob/master/jb-icon.jpg)

## config

Reads in **json** key value pairs that are unmarshaled into one configuration struct that can be passed around through an application.

### Example config/application.json:

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
	App App
	Datasource Datasource
}
```

This package relies on a directory that should be at the root of the application called `config`. The **Unmarshal()** function looks for two json files named `application.json` and `datasource.json`.

- config
	- application.json  // *gets read into the App struct*
	- datasource.json // *gets read into the Datasource struct*

*App - configurables for common application related objects.*

``` go
type App struct {
	Name string
	Port int
	LogLevel string
}
```

*Datasource - configurations for one or more mongo database objects.*

``` go
type Datasource struct {
	Database  Database
	Databases map[string]Database
}

type Database struct {
	Database string
	URI string
	Collections map[string]string
}
```

## log

Utilizes [sirupsen/logrus](https://github.com/sirupsen/logrus "sirupsen/logrus") logging pacakge. Wraps a few functions to help configure and help make log messages more informational.

### Example json message:

```json
{
	"application":"api",
	"file":"file.go:10",
	"function":"func()",
	"level":"info",
	"msg":"a message",
	"time":"2020-06-09T12:30:00-05:00"
}
```

**Note**: all log messages will have *application*, *level*, *msg* and *time* fields by default. See **Usage**  for how to log with additional fields.

### Set Up

``` go
package main

import "github.com/jobaldw/shared/log"

func main() {
    ...
    log.Configure("api", "info")
    ...
}
```

### Usage
``` go
func foo() {
    log.Info("message") // prints default fields
    log.Details().Info("message") // prints default fields, file, and function
    log.Add(log.Fields{"id": 123-456-7890}).Info("message") // prints Details() and any additional fields specified
}
```