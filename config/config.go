package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

// The root directory of any project
const root = "."

var (
	// The config directory name where files with JSON value configs will live.
	configDirectory = "configs"

	// This error will normally be thrown if the configDirectory does not exist.
	ErrConfigsNotFound = errors.New("could not locate config values within the project")

	// The struct must be a pointer for json.unmarshaling purposes. If the value
	// is nil or not a pointer, json.Unmarshal returns an InvalidUnmarshalError.
	ErrNonPointerStruct = errors.New("configuration should be a pointer to a struct type")
)

// The application struct holds values specific to the apps configuration.
// Intended on being used alone or in combination with other structs for
// users to utilize.
type Application struct {
	// Application name primarily used for logging/debugging purposes.
	Name string `json:"name,omitempty"`

	// Server port that the microservice communicates through.
	Port int `json:"port,omitempty"`

	// Used to set logging severity. Field is a string value to users can use
	// this value with any logging packages such as zerolog, logrus, viper
	// or an internal logging package.
	LogLevel string `json:"log_level,omitempty"`
}

// Holds multiple Client objects that can be used within the app via a map to
// allow users to keep client configurations separate.
type Clients struct {
	// A map of string to Client configs
	Clients map[string]Client `json:"clients,omitempty"`
}

// This struct holds configurations for a client from health checks and base
// urls to client timeouts and retry maxes.
type Client struct {
	// Map of string to string header options to be used with http.Headers.
	// For client requests, certain headers such as Content-Length
	// and Connection are automatically written when needed and
	// values in Header may be ignored.
	Headers map[string]string `json:"headers,omitempty"`

	// Health check path used for pinging the client
	Health string `json:"health,omitempty"`

	// A time limit for requests made by the client. The duration includes
	// connection time, redirects and reading the response. A Timeout of
	// zero or omitted means no timeout.
	Timeout int `json:"timeout,omitempty"`

	// The client's base url.
	URL string `json:"url,omitempty"`
}

// Holds multiple Mongo objects that can be used within the app via a map to
// allow users to keep mongo configurations separate.
type Datasource struct {
	// A map of string to Mongo configs
	Mongo map[string]Mongo `json:"mongo,omitempty"`
}

// The mongo struct stores configurations that are primarily used with the
// official mongo driver package.
type Mongo struct {
	// The name of the database to connect to.
	Database string `json:"database,omitempty"`

	// mongo uri with authentication encoded in base64. Should be in mongodb+svr://
	// form before encoding.
	URI string `json:"uri,omitempty"`

	// collections that exist within the defined database.
	Collections map[string]string `json:"collections,omitempty"`
}

/** Unmarshal
 * Reads in located JSON config files to parse into any given config struct
 * @param config: A pointer to a struct where values will be read into */
func Unmarshal(config interface{}) error {
	// check if given configuration is a pointer to a struct
	configType := reflect.ValueOf(config).Type()
	if configType.Kind() != reflect.Ptr || configType.Elem().Kind() != reflect.Struct {
		return ErrNonPointerStruct
	}

	// read the root director of the project
	root, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	// search for the config directory and its files
	for _, folder := range root {
		if folder.IsDir() && folder.Name() == configDirectory {
			configs, err := ioutil.ReadDir(folder.Name())
			if err != nil {
				return err
			}
			for _, file := range configs {
				path := fmt.Sprintf("./%s/%s", folder.Name(), file.Name())
				if err := unmarshal(path, config); err != nil {
					return err
				}
			}
			return nil
		}
	}

	return ErrConfigsNotFound
}

/** unmarshal
 * Helper function that parses config json values into any given struct
 * @param path: File location where config values are stored
 * @param config: A pointer to a struct where values will be read into */
func unmarshal(path string, config interface{}) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		// open file location
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf(`%v, could not open "%s"`, err, path)
		}
		defer file.Close()

		// read file contents
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return fmt.Errorf(`%v, could not read "%s"`, err, path)
		}

		// Store JSON data into struct. To unmarshal JSON into a pointer, this
		// function first handles the case of the JSON being the JSON literal null.
		// In that case, it sets the pointer to nil. Otherwise, it unmarshals the
		// JSON into the value pointed at by the pointer. If the pointer is nil,
		// jon.Unmarshal allocates a new value for it to point to.
		return json.Unmarshal(data, config)
	}
	return ErrConfigsNotFound
}
