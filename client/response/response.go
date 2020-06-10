package response

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Response struct
type Response struct {
	Status string
	Code   int

	Body Body
}

// Body struct
type Body struct {
	String string
	Bytes  []byte
	IO     io.Reader
}

// Marshal response body
func Marshal(body io.Reader, structure interface{}) error {
	err := json.NewDecoder(body).Decode(structure)
	return fmt.Errorf("could not marshal response, %s", err)
}

// Save response body
func (r *Response) Save(resp *http.Response) {
	r.Status = resp.Status
	r.Code = resp.StatusCode
	r.Body.save(resp.Body)
}

// helper functions
func (b *Body) save(body io.Reader) {
	b.IO = body

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		b.String = ""
		b.Bytes = nil
		return
	}

	b.String = string(bodyBytes)
	b.Bytes = bodyBytes
}
