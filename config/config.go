package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

const (
	dataSource = "config/datasource.json"
	appSource  = "config/application.json"
)

// Config struct
type Config struct {
	App        App
	Datasource Datasource
}

// App struct
type App struct {
	Name       string `json:"application,omitempty" bson:"application,omitempty"`
	Port       int    `json:"port,omitempty" bson:"port,omitempty"`
	LogLevel   string `json:"log_level,omitempty" bson:"log_level,omitempty" `
	StackTrace bool   `json:"stack_trace,omitempty" bson:"stack_trace,omitempty"`
}

// Datasource struct
type Datasource struct {
	Database  Database            `json:"database,omitempty"`
	Databases map[string]Database `json:"databases,omitempty"`
}

// Database struct
type Database struct {
	Database    string            `json:"database,omitempty"`
	URI         string            `json:"uri,omitempty"`
	Collections map[string]string `json:"collections,omitempty"`
}

// Marshal configurables and configure logging
func Marshal() (conf Config, err error) {
	if hasSource(appSource) {
		err = read(appSource, &conf.App)
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

	return
}

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

	err = yaml.Unmarshal(data, &configuration)
	if err != nil {
		return
	}

	return
}
