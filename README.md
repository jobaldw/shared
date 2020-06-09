# shared

> A library of code that can be used across Go applications.

![JB Desgins](https://github.com/jobaldw/shared/blob/master/jb-icon.jpg)

## config

Reads in **json** key value pairs that are unmarshaled into one configuration struct that can be passed around through an application.

*Config - configuration that holds all configurables.*

``` golang
type Config struct {
	App App
	Datasource Datasource
}
```

This package relies on a directory that should be at the root of the application called `config`. The **Unmarshal()** function looks for two json files named `application.json` and `datasource.json`.

- config<br>
-- application.json  // *gets read into the App struct*<br>
-- datasource.json // *gets read into the Datasource struct*<br>

*App - configurables for common application related objects.*<br>
*Datasource - configurations for one or more mongo database objects.*<br>

``` golang
type App struct {
	Name        string
	Port          int
	LogLevel   string
}
```

``` golang

type Datasource struct {
	Database  Database
	Databases map[string]Database
}

type Database struct {
	Database   string
	URI            string
	Collections map[string]string
}
```




## log