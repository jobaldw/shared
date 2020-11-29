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
	appSource  = "application.json"
	dataSource = "datasource.json"
	depSource  = "dependents.json"
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
	Auth0    Auth0  `json:"auth0,omitempty"`
}

// Auth0 struct
type Auth0 struct {
	Identifier string `json:"identifier,omitempty"`
	Domain     string `json:"domain,omitempty"`
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
	Health  string            `json:"health,omitempty"`
	URL     string            `json:"url,omitempty"`
	Timeout int               `json:"timeout,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

// Unmarshal configurables
func Unmarshal(dir string) (conf Config, err error) {
	if hasSource(dir + appSource) {
		if err = Read(dir+appSource, &conf.Application); err != nil {
			return conf, fmt.Errorf("%s %s'", appSource, err)
		}
	}

	if hasSource(dir + dataSource) {
		if err = Read(dir+dataSource, &conf.Datasource); err != nil {
			return conf, fmt.Errorf("%s %s'", dataSource, err)
		}
	}

	if hasSource(dir + depSource) {
		if err = Read(dir+depSource, &conf.Dependents); err != nil {
			return conf, fmt.Errorf("%s %s'", depSource, err)
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

func Read(filename string, configuration interface{}) error {
	configValue := reflect.ValueOf(configuration)
	if typ := configValue.Type(); typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return errors.New("configuration should be a pointer to a struct type")
	}

	return getValues(filename, configuration)
}

func getValues(filename string, configuration interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &configuration)
	if err != nil {
		return err
	}

	return err
}
