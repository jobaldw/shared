package config

import (
	"os"

	"github.com/jobaldw/shared/log"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
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
	Name       string       `json:"application,omitempty" bson:"application,omitempty"`
	Port       int          `json:"port,omitempty" bson:"port,omitempty"`
	LogLevel   logrus.Level `json:"log_level,omitempty" bson:"log_level,omitempty" `
	StackTrace bool         `json:"stack_trace,omitempty" bson:"stack_trace,omitempty"`
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
		err = gonfig.GetConf(appSource, &conf.App)
		if err != nil {
			return conf, errors.WithMessage(err, "could not read appliation configurations")
		}
	}

	if hasSource(dataSource) {
		err = gonfig.GetConf(dataSource, &conf.Datasource)
		if err != nil {
			return conf, errors.WithMessage(err, "could not read datasource configurations")
		}
	}

	log.Configure(log.Logger{
		App:       conf.App.Name,
		Level:     conf.App.LogLevel,
		ShowStack: conf.App.StackTrace,
	})

	log.Entry.WithField("method", "Marshal").Debugf("%s configured.", conf.App.Name)
	return
}

func hasSource(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	return true
}
