package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

const (
	appSource  = "config/application.json"
	dataSource = "config/datasource.json"
	depSource  = "config/dependents.json"
)

// Config struct
type Config struct {
	Application Application
	Datasource  Datasource
	Dependents  Dependents
}

// Application struct
type Application struct {
	Name     string `json:"application,omitempty"`
	Port     int    `json:"port,omitempty"`
	LogLevel string `json:"log_level,omitempty"`
}

// Datasource struct
type Datasource struct {
	Mongo  Mongo            `json:"database,omitempty"`
	Mongos map[string]Mongo `json:"databases,omitempty"`
}

// Mongo struct
type Mongo struct {
	Database    string            `json:"database,omitempty"`
	URI         string            `json:"uri,omitempty"`
	Collections map[string]string `json:"collections,omitempty"`
}

// Dependents struct
type Dependents struct {
	Client  Client            `json:"client,omitempty"`
	Clients map[string]Client `json:"clients,omitempty"`
}

// Client struct
type Client struct {
	URL        string            `json:"url,omitempty"`
	Timeout    int               `json:"timeout,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
	Health     string            `json:"health,omitempty"`
}

// Unmarshal configurables
func Unmarshal() (conf Config, err error) {
	if hasSource(appSource) {
		err = read(appSource, &conf.Application)
		if err != nil {
			return conf, fmt.Errorf("%s, %s '%s'", err, "could not read", appSource)
		}
	}

	if hasSource(dataSource) {
		err = read(dataSource, &conf.Datasource)
		if err != nil {
			return conf, fmt.Errorf("%s, %s '%s'", err, "could not read", dataSource)
		}
	}

	if hasSource(depSource) {
		err = read(depSource, &conf.Dependents)
		if err != nil {
			return conf, fmt.Errorf("%s, %s '%s'", err, "could not read", depSource)
		}
	}

	return
}

// Helper functions
func hasSource(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	return true
}

func read(filename string, configuration interface{}) (err error) {
	configValue := reflect.ValueOf(configuration)
	if typ := configValue.Type(); typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return errors.New("configuration should be a pointer to a struct type")
	}

	return getValues(filename, configuration)
}

func getValues(filename string, configuration interface{}) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &configuration)
	if err != nil {
		return
	}

	return
}
