package response

import (
	"io"
	"io/ioutil"
	"net/http"
)

// Response struct
type Response struct {
	Status string
	Code   int

	Body    Body
	Request *http.Request
}

// Body struct
type Body struct {
	String string
	Bytes  []byte
	IO     io.Reader
}

// Save response body
func (r *Response) Save(resp *http.Response) {
	r.Status = resp.Status
	r.Code = resp.StatusCode
	r.Request = resp.Request
	r.Body.save(resp.Body)
}

// helper function
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
