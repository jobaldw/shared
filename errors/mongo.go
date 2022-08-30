package errors

import "fmt"

// Mongo operation key
const (
	Insert = "insert"
	Update = "update"
	Read   = "read"
	Delete = "delete"
)

var (
	// an error ocurred trying to delete a document
	ErrDeleteFailed = fmt.Errorf("could not remove document")

	// an error ocurred trying to create a new document
	ErrInsertFailed = fmt.Errorf("could not create document")

	// an error ocurred trying to update a document
	ErrUpdateFailed = fmt.Errorf("could not update document")

	// an error ocurred trying to retrieve a document
	ErrNoBudgetsFound = fmt.Errorf("no documents were found")
)

// struct for mongo errors
type MongoErr struct {
	Operation  string `json:"operation,omitempty"`
	Collection string `json:"collection,omitempty"`
	Err        error  `json:"error,omitempty"`
}

// NewMongoErr
//
//	Creates and fills a new MongoErr struct.
//	* @param op: the operation
//	* @param col: the mongo collection
//	* @param err: an error
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
//
// Error implements the error interface.
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
