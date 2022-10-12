package errors

import "fmt"

const (
	Insert = "insert" // Insert Mongo operation key
	Update = "update" // Update Mongo operation key
	Read   = "read"   // Read Mongo operation key
	Delete = "delete" // Delete Mongo operation key
)

var (
	ErrDeleteFailed   = fmt.Errorf("could not remove document") // an error occurred trying to delete a document
	ErrInsertFailed   = fmt.Errorf("could not create document") // an error occurred trying to create a new document
	ErrUpdateFailed   = fmt.Errorf("could not update document") // an error occurred trying to update a document
	ErrNoBudgetsFound = fmt.Errorf("no documents were found")   // an error occurred trying to retrieve a document
)

// struct for mongo errors
type MongoErr struct {
	// Mongo operation that the error occurred
	Operation string `json:"operation,omitempty"`

	// the Mongo collection name
	Collection string `json:"collection,omitempty"`

	// the error that occurred
	Err error `json:"error,omitempty"`
}

// NewMongoErr
// creates and fills a new MongoErr struct.
func NewMongoErr(op, col string, err error) MongoErr {
	operation := ""
	switch op {
	case Delete, Insert, Update, Read:
		operation = op
	}

	collection := "unknown"
	if col != "" {
		collection = col
	}

	return MongoErr{
		Operation:  operation,
		Collection: collection,
		Err:        err,
	}
}

// Error
// implements the error interface.
func (me MongoErr) Error() string {
	err := me.Err
	if err == nil {
		switch me.Operation {
		case Delete:
			err = ErrDeleteFailed
		case Insert:
			err = ErrInsertFailed
		case Update:
			err = ErrUpdateFailed
		case Read:
			err = ErrNoBudgetsFound
		}
	}

	if me.Operation == "" {
		return fmt.Sprintf("[%s] %v", me.Collection, err)
	}

	return fmt.Sprintf("[%s] %s error: %v", me.Collection, me.Operation, err)
}
